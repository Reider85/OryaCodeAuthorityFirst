// Package universallsp — Tier-4 universal LSP fallback адаптер (FR-ADP-004,
// ADR-013). Символьный граф строит gopls через стандартный LSP-протокол
// поверх stdio; никаких выдуманных TCP-адресов из stderr.
//
// Контракт ARCHITECT.md §3.1: Install()/Uninstall() идемпотентны, Detect()
// отвечает «есть ли gopls на машине». Ядро пакета — client.go (JSON-RPC 2.0
// с Content-Length фреймингом); adapter.go добавляет 5 методов интерфейса
// Adapter поверх него.
package universallsp

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultBinary — имя LSP-сервера в $PATH. MVP не бандлит gopls (ADR-026:
// require install), детекция и запуск идут по системному PATH.
const DefaultBinary = "gopls"

// ResponseError — ошибка JSON-RPC от LSP-сервера.
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("lsp error %d: %s", e.Code, e.Message)
}

// response — входящее JSON-RPC сообщение. ID *int: у уведомлений (server →
// client, например window/logMessage) id отсутствует — такие кадры мы
// молча пропускаем, они не являются ответом на наши запросы.
type response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *int            `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// request — исходящий JSON-RPC запрос.
type request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Client — минимальный LSP-клиент поверх stdio.
//
// Я отказался от прежней схемы «gopls -rpc.address + парсинг stderr»: это
// отладочный RPC gopls, а не LSP-протокол, и без Content-Length-фрейминга
// сервер вообще не понимает наши сообщения. Стандартный транспорт LSP —
// stdin/stdout процесса с заголовками Content-Length; здесь он и
// реализован. Вся синхронизация — через mutex + каналы ожидания ответа
// по id; читательская горутина единственная владеет stdout.
type Client struct {
	// Binary — путь к LSP-серверу; пусто = DefaultBinary. Отдельное поле,
	// чтобы тесты могли подставить fake-сервер без подмены $PATH.
	Binary string

	mu      sync.Mutex
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	nextID  int
	pending map[int]chan response
	done    chan struct{} // закрывается читательской горутиной при смерти stdout
	readErr error
	started bool
}

// NewClient создаёт клиент (процесс не запускается до Start).
func NewClient() *Client {
	return &Client{}
}

func (c *Client) binName() string {
	if c.Binary != "" {
		return c.Binary
	}
	return DefaultBinary
}

// Started сообщает, жив ли LSP-процесс (после Start и до Stop/смерти).
func (c *Client) Started() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.started
}

// Start запускает LSP-сервер и выполняет handshake (initialize +
// initialized). ctx ограничивает время handshake — инвариант INV-ADR-013:
// бюджет старта ~500ms на малом репозитории; жёсткий потолок задаёт
// вызывающий (FallbackManager — 2s, см. startupTimeout).
func (c *Client) Start(ctx context.Context, rootPath string) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return nil // идемпотентность: повторный Start не убивает живой процесс
	}
	c.mu.Unlock()

	cmd := exec.Command(c.binName())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("lsp: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("lsp: stdout pipe: %w", err)
	}
	// stderr обязателен к сливаю: если буфер pipe переполнится, gopls
	// заблокируется — это ровно тот класс зависаний, который был у прежней
	// реализации с брошенными pipe'ами.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("lsp: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("lsp: start %s: %w", c.binName(), err)
	}

	c.mu.Lock()
	c.cmd = cmd
	c.stdin = stdin
	c.pending = make(map[int]chan response)
	c.done = make(chan struct{})
	c.readErr = nil
	c.mu.Unlock()

	go func() { _, _ = io.Copy(io.Discard, stderr) }() // дренаж stderr
	go c.readLoop(stdout)

	if err := c.handshake(ctx, rootPath); err != nil {
		_ = c.Stop()
		return err
	}

	c.mu.Lock()
	c.started = true
	c.mu.Unlock()
	return nil
}

// handshake: initialize (ждём ответ) → initialized (уведомление, без ответа).
func (c *Client) handshake(ctx context.Context, rootPath string) error {
	params := map[string]interface{}{
		"processId": os.Getpid(),
		"rootUri":   pathToURI(rootPath),
		"capabilities": map[string]interface{}{
			"workspace": map[string]interface{}{
				"symbol": map[string]interface{}{"dynamicRegistration": false},
			},
		},
		"trace": "off",
	}
	if _, err := c.Call(ctx, "initialize", params); err != nil {
		return fmt.Errorf("lsp: initialize: %w", err)
	}
	if err := c.Notify("initialized", map[string]interface{}{}); err != nil {
		return fmt.Errorf("lsp: initialized: %w", err)
	}
	return nil
}

// Call отправляет запрос и ждёт ответ по id, ограниченный ctx.
func (c *Client) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	if !c.started && c.cmd == nil {
		c.mu.Unlock()
		return nil, errors.New("lsp: client not started")
	}
	c.nextID++
	id := c.nextID
	ch := make(chan response, 1)
	if c.pending == nil {
		c.pending = make(map[int]chan response)
	}
	c.pending[id] = ch
	done := c.done
	c.mu.Unlock()

	req := request{JSONRPC: "2.0", ID: id, Method: method, Params: params}
	if err := c.writeMessage(req); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, resp.Error
		}
		return resp.Result, nil
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("lsp: %s: %w", method, ctx.Err())
	case <-done:
		c.mu.Lock()
		err := c.readErr
		c.mu.Unlock()
		if err == nil {
			err = errors.New("lsp: connection closed")
		}
		return nil, fmt.Errorf("lsp: %s: %w", method, err)
	}
}

// Notify отправляет уведомление (без id, ответа не ждём).
func (c *Client) Notify(method string, params interface{}) error {
	return c.writeMessage(request{JSONRPC: "2.0", Method: method, Params: params})
}

// writeMessage сериализует JSON-RPC сообщение с Content-Length заголовком
// (LSP base protocol): "Content-Length: N\r\n\r\n{body}".
func (c *Client) writeMessage(msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("lsp: marshal: %w", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stdin == nil {
		return errors.New("lsp: stdin closed")
	}
	var buf strings.Builder
	fmt.Fprintf(&buf, "Content-Length: %d\r\n\r\n", len(body))
	buf.Write(body)
	if _, err := io.WriteString(c.stdin, buf.String()); err != nil {
		return fmt.Errorf("lsp: write: %w", err)
	}
	return nil
}

// readLoop — единственная читательница stdout: разбирает кадры и маршрутизирует
// ответы ожидающим горутинам по id. Уведомления (id == nil) отбрасываются —
// MVP не подписывается на diagnostics/logMessage.
func (c *Client) readLoop(stdout io.Reader) {
	defer close(c.done)
	reader := bufio.NewReader(stdout)
	for {
		body, err := readFrame(reader)
		if err != nil {
			c.mu.Lock()
			// io.EOF при штатном Stop — не ошибка, а конец жизни процесса.
			if !errors.Is(err, io.EOF) && c.readErr == nil {
				c.readErr = err
			}
			c.mu.Unlock()
			return
		}
		var resp response
		if err := json.Unmarshal(body, &resp); err != nil {
			continue // битый кадр — пропускаем, соединение живо
		}
		if resp.ID == nil {
			continue // уведомление сервера
		}
		c.mu.Lock()
		ch := c.pending[*resp.ID]
		delete(c.pending, *resp.ID)
		c.mu.Unlock()
		if ch != nil {
			// Неблокирующая отправка: если Stop уже разбудил ожидающих
			// ошибкой, канал заполнен — второй ответ никому не нужен.
			select {
			case ch <- resp:
			default:
			}
		}
	}
}

// readFrame читает один LSP-кадр: заголовки до пустой строки, затем body
// ровно Content-Length байт. Поддерживает заголовки в любом порядке и
// игнорирует прочие (Content-Type и т.п.).
func readFrame(r *bufio.Reader) ([]byte, error) {
	length := -1
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break // конец заголовков
		}
		name, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(name), "Content-Length") {
			n, err := strconv.Atoi(strings.TrimSpace(value))
			if err != nil || n < 0 {
				return nil, fmt.Errorf("lsp: bad Content-Length %q", value)
			}
			length = n
		}
	}
	if length < 0 {
		return nil, errors.New("lsp: missing Content-Length header")
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}
	return body, nil
}

// Stop убивает LSP-процесс (kill + короткий wait, чтобы не оставлять зомби)
// и закрывает stdin. Идемпотентен: повторный Stop — no-op.
func (c *Client) Stop() error {
	c.mu.Lock()
	cmd := c.cmd
	stdin := c.stdin
	c.cmd = nil
	c.stdin = nil
	c.started = false
	pending := c.pending
	c.pending = nil
	c.mu.Unlock()

	// Разбудить ожидающих Call явной ошибкой — иначе они получили бы
	// нулевой ответ (закрытый канал) и интерпретировали бы его как успех
	// с пустым result.
	for _, ch := range pending {
		ch <- response{Error: &ResponseError{Code: -32000, Message: "lsp: connection closed"}}
	}
	if stdin != nil {
		_ = stdin.Close()
	}
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	_ = cmd.Process.Kill()
	waitCh := make(chan error, 1)
	go func() { waitCh <- cmd.Wait() }()
	select {
	case <-waitCh:
	case <-time.After(2 * time.Second):
		// Процесс не умер за 2s — считаем остановленным настолько, насколько
		// это важно для MVP; осироты не оставляем: Wait уже запущен.
	}
	return nil
}

// --- маленькие помощники, вынесенные для тестируемости ---

// pathToURI конвертирует файловый путь в file:// URI (LSP rootUri).
func pathToURI(path string) string {
	if path == "" {
		return ""
	}
	// Windows: C:\foo → file:///C:/foo; Unix: /foo → file:///foo.
	p := "/" + strings.ReplaceAll(path, "\\", "/")
	if runtime.GOOS != "windows" && strings.HasPrefix(p, "//") {
		p = p[1:] // UNC-подобные двойные слэши не плодим
	}
	return "file://" + p
}
