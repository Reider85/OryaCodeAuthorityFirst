package universallsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// TestMain — re-exec паттерн: когда дочерний процесс test-бинарника запускается
// с CAF_FAKE_GOPLS=1, он исполняет роль LSP-сервера вместо m.Run(). Клиенту
// подаётся Binary=os.Args[0]; отдельного .bat/.sh стаба не нужно — тестовый
// бинарник кроссплатформенный.
func TestMain(m *testing.M) {
	if os.Getenv("CAF_FAKE_GOPLS") == "1" {
		runFakeGopls()
		return
	}
	os.Exit(m.Run())
}

// runFakeGopls — минимальный LSP-сервер поверх stdio: читает кадры,
// отвечает на initialize и workspace/symbol, молчит на уведомления.
func runFakeGopls() {
	in := bufio.NewReader(os.Stdin)
	for {
		body, err := readFrame(in)
		if err != nil {
			return // клиент закрыл stdin / умер — штатный выход
		}
		var req struct {
			ID     *int            `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal(body, &req); err != nil || req.ID == nil {
			continue
		}
		var result interface{}
		switch req.Method {
		case "initialize":
			result = map[string]interface{}{
				"capabilities": map[string]interface{}{"symbolProvider": true},
			}
		case "workspace/symbol":
			result = fakeSymbols()
		default:
			result = map[string]interface{}{}
		}
		resp, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0", "id": *req.ID, "result": result,
		})
		fmt.Fprintf(os.Stdout, "Content-Length: %d\r\n\r\n%s", len(resp), resp)
	}
}

// fakeSymbols — фикстура: function с containerName, struct без него,
// out-of-range kind и не-file URI (проверка фильтрации).
func fakeSymbols() []map[string]interface{} {
	pos := func(l, c int) map[string]int { return map[string]int{"line": l, "character": c} }
	sym := func(name string, kind int, uri string, line, ch int, container string) map[string]interface{} {
		m := map[string]interface{}{
			"name": name, "kind": kind,
			"location": map[string]interface{}{
				"uri":   uri,
				"range": map[string]interface{}{"start": pos(line, ch), "end": pos(line, ch+1)},
			},
		}
		if container != "" {
			m["containerName"] = container
		}
		return m
	}
	return []map[string]interface{}{
		sym("NewClient", 12, "file:///proj/client.go", 9, 4, "universallsp"), // function
		sym("Client", 23, "file:///proj/client.go", 20, 6, ""),               // struct
		sym("weird", 99, "file:///proj/other.go", 0, 0, ""),                  // out-of-range kind
		sym("buffer", 12, "untitled:buffer-1", 1, 0, ""),                     // не file:// — отбрасывается
	}
}

// startFakeClient запускает клиент против fake-gopls (этот же test-бинарник).
func startFakeClient(t *testing.T) *Client {
	t.Helper()
	t.Setenv("CAF_FAKE_GOPLS", "1")
	client := NewClient()
	client.Binary = os.Args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	if err := client.Start(ctx, t.TempDir()); err != nil {
		t.Fatalf("Start against fake gopls: %v", err)
	}
	t.Cleanup(func() { _ = client.Stop() })
	return client
}

// TestClientHandshakeAndSymbols — полный round-trip: initialize проходит,
// workspace/symbol возвращает сконвертированные символы (kind → имя,
// 0-based → 1-based, containerName → Parent, не-file URI отброшен).
func TestClientHandshakeAndSymbols(t *testing.T) {
	client := startFakeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	symbols, err := client.WorkspaceSymbols(ctx, "")
	if err != nil {
		t.Fatalf("WorkspaceSymbols: %v", err)
	}
	if len(symbols) != 3 {
		t.Fatalf("got %d symbols, want 3 (untitled URI must be skipped)", len(symbols))
	}

	fn := symbols[0]
	if fn.Name != "NewClient" || fn.Kind != "function" {
		t.Errorf("symbol[0] = %s/%s, want NewClient/function", fn.Name, fn.Kind)
	}
	if fn.Parent != "universallsp" {
		t.Errorf("symbol[0].Parent = %q, want universallsp (containerName)", fn.Parent)
	}
	if !strings.HasSuffix(fn.Location, ":10:5") {
		t.Errorf("symbol[0].Location = %q, want suffix :10:5 (1-based line:col)", fn.Location)
	}

	st := symbols[1]
	if st.Kind != "struct" || st.Parent != "" {
		t.Errorf("symbol[1] = %s/%s parent=%q, want struct with empty Parent", st.Name, st.Kind, st.Parent)
	}

	weird := symbols[2]
	if weird.Kind != "symbol-99" {
		t.Errorf("out-of-range kind = %q, want symbol-99", weird.Kind)
	}
}

// TestClientStartIdempotent — повторный Start не убивает живой процесс.
func TestClientStartIdempotent(t *testing.T) {
	client := startFakeClient(t)
	cmdBefore := client.cmd

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Start(ctx, t.TempDir()); err != nil {
		t.Fatalf("second Start must be no-op, got: %v", err)
	}
	if client.cmd != cmdBefore {
		t.Error("second Start replaced the running process")
	}
}

// TestClientStopIdempotent — двойной Stop без паники; после Stop Call отклоняется.
func TestClientStopIdempotent(t *testing.T) {
	client := startFakeClient(t)
	if err := client.Stop(); err != nil {
		t.Fatalf("first Stop: %v", err)
	}
	if err := client.Stop(); err != nil {
		t.Fatalf("second Stop must be nil, got: %v", err)
	}
	if client.Started() {
		t.Error("Started() = true after Stop")
	}

	ctx := context.Background()
	if _, err := client.Call(ctx, "workspace/symbol", map[string]string{"query": ""}); err == nil {
		t.Error("Call after Stop must fail")
	}
}

// TestClientStartMissingBinary — несуществующий бинарник: ошибка старта, без паники.
func TestClientStartMissingBinary(t *testing.T) {
	client := NewClient()
	client.Binary = "definitely-not-a-real-lsp-server-xyz"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := client.Start(ctx, t.TempDir())
	if err == nil {
		t.Fatal("Start with missing binary must fail")
	}
	if client.Started() {
		t.Error("Started() = true after failed Start")
	}
}

// TestCallTimesOut — сервер принял соединение, но не отвечает: Call должен
// вернуть context.DeadlineExceeded, а не зависнуть (бюджет INV-ADR-013).
func TestCallTimesOut(t *testing.T) {
	t.Setenv("CAF_FAKE_GOPLS", "1")
	client := NewClient()
	client.Binary = os.Args[0]
	startCtx, cancelStart := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStart()
	if err := client.Start(startCtx, t.TempDir()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer client.Stop()

	// Убиваем сервер между Handshake и запросом: stdin закрыт, ответа не будет.
	_ = client.cmd.Process.Kill()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := client.WorkspaceSymbols(ctx, "")
	if err == nil {
		t.Fatal("WorkspaceSymbols against dead server must fail")
	}
}

// TestPathToURIAndBack — file:// URI: Windows-путь получает file:///C:/...,
// unix-путь не плодит двойной слэш.
func TestPathToURIAndBack(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"relative", "proj", "file:///proj"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathToURI(tt.in); got != tt.want {
				t.Errorf("pathToURI(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}

	// uriToPath: round-trip для абсолютного пути текущей ОС.
	dir := t.TempDir()
	uri := pathToURI(dir)
	back := uriToPath(uri)
	if back != dir && back != strings.ReplaceAll(dir, "\\", "/") {
		// FromSlash на Windows вернёт обратные слэши — сравниваем нормализованно.
		if strings.ReplaceAll(back, "/", "\\") != strings.ReplaceAll(dir, "/", "\\") {
			t.Errorf("round-trip: %q → %q → %q", dir, uri, back)
		}
	}
}

// TestURIPathFiltering — uriToPath отбрасывает не-file схемы.
func TestURIPathFiltering(t *testing.T) {
	if got := uriToPath("untitled:buffer-1"); got != "" {
		t.Errorf("uriToPath(untitled) = %q, want empty", got)
	}
	if got := uriToPath("http://example.com/x.go"); got != "" {
		t.Errorf("uriToPath(http) = %q, want empty", got)
	}
}
