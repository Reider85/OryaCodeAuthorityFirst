package universallsp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// TestReadFrameBasic — стандартный кадр Content-Length читается целиком.
func TestReadFrameBasic(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":1,"result":{}}`
	raw := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)

	got, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatalf("readFrame: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", got, body)
	}
}

// TestReadFrameExtraHeaders — Content-Type и прочие заголовки игнорируются;
// порядок заголовков произволен.
func TestReadFrameExtraHeaders(t *testing.T) {
	body := `{"jsonrpc":"2.0","id":2,"result":null}`
	raw := "Content-Type: application/vscode-jsonrpc; charset=utf-8\r\n" +
		fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)

	got, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatalf("readFrame: %v", err)
	}
	if string(got) != body {
		t.Errorf("got %q, want %q", got, body)
	}
}

// TestReadFrameSequential — два кадра подряд читаются последовательно.
func TestReadFrameSequential(t *testing.T) {
	var buf bytes.Buffer
	for i, body := range []string{`{"id":1}`, `{"id":2}`} {
		fmt.Fprintf(&buf, "Content-Length: %d\r\n\r\n%s", len(body), body)
		_ = i
	}
	r := bufio.NewReader(&buf)

	first, err := readFrame(r)
	if err != nil || string(first) != `{"id":1}` {
		t.Fatalf("first frame: %q, %v", first, err)
	}
	second, err := readFrame(r)
	if err != nil || string(second) != `{"id":2}` {
		t.Fatalf("second frame: %q, %v", second, err)
	}
}

// TestReadFrameMissingLength — кадр без Content-Length: ошибка, не зависание.
func TestReadFrameMissingLength(t *testing.T) {
	raw := "Content-Type: application/json\r\n\r\n{}"
	_, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err == nil {
		t.Fatal("missing Content-Length must be an error")
	}
}

// TestReadFrameBadLength — нечисловой Content-Length: ошибка.
func TestReadFrameBadLength(t *testing.T) {
	raw := "Content-Length: banana\r\n\r\n{}"
	_, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err == nil {
		t.Fatal("non-numeric Content-Length must be an error")
	}
}

// TestReadFrameTruncatedBody — тело обрывается: io.EOF/ErrUnexpectedEOF, не молчание.
func TestReadFrameTruncatedBody(t *testing.T) {
	raw := "Content-Length: 100\r\n\r\n{}"
	_, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err == nil {
		t.Fatal("truncated body must be an error")
	}
}

// TestReadFrameEOFMidHeaders — обрыв в середине заголовков.
func TestReadFrameEOFMidHeaders(t *testing.T) {
	raw := "Content-Length: 5\r\n" // без пустой строки-терминатора
	_, err := readFrame(bufio.NewReader(strings.NewReader(raw)))
	if err == nil {
		t.Fatal("EOF mid-headers must be an error")
	}
	if err != io.EOF {
		t.Errorf("got %v, want io.EOF", err)
	}
}
