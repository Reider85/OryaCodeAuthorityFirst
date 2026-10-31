package universallsp

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// lspSymbolInformation — структура SymbolInformation из спецификации LSP
// (https://microsoft.github.io/language-server-protocol/). Kind — целое
// 1..26, НЕ строка: прежняя реализация объявляла его string, и каст
// молча отбрасывал каждый символ.
type lspSymbolInformation struct {
	Name          string      `json:"name"`
	Kind          int         `json:"kind"`
	ContainerName string      `json:"containerName,omitempty"`
	Location      lspLocation `json:"location"`
}

type lspLocation struct {
	URI   string   `json:"uri"`
	Range lspRange `json:"range"`
}

type lspRange struct {
	Start lspPosition `json:"start"`
	End   lspPosition `json:"end"`
}

type lspPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// symbolKindNames — числовой SymbolKind → человекочитаемое имя (LSP 3.17).
// Индекс = kind; нулевой элемент не используется (kind начинается с 1).
var symbolKindNames = []string{
	"",
	"file",
	"module",
	"namespace",
	"package",
	"class",
	"method",
	"property",
	"field",
	"constructor",
	"enum",
	"interface",
	"function",
	"variable",
	"constant",
	"string",
	"number",
	"boolean",
	"array",
	"object",
	"key",
	"null",
	"enum-member",
	"struct",
	"event",
	"operator",
	"type-parameter",
}

// WorkspaceSymbols выполняет workspace/symbol и конвертирует результат в
// []caf.Symbol. Пустой query в gopls возвращает символы всего workspace —
// ровно то, что нужно для кэша .caf/mental-model/graph.yml (ADR-013).
func (c *Client) WorkspaceSymbols(ctx context.Context, query string) ([]caf.Symbol, error) {
	raw, err := c.Call(ctx, "workspace/symbol", map[string]string{"query": query})
	if err != nil {
		return nil, fmt.Errorf("workspace/symbol: %w", err)
	}
	// result может быть null (нет совпадений) — это не ошибка.
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var infos []lspSymbolInformation
	if err := json.Unmarshal(raw, &infos); err != nil {
		return nil, fmt.Errorf("workspace/symbol: decode: %w", err)
	}
	return convertSymbols(infos), nil
}

// convertSymbols переводит LSP-символы в формат caf.Symbol.
// Location — "path:line:column" (строки и колонки LSP 0-based, привожу к
// 1-based: человеческий формат, как в `go vet`-сообщениях).
func convertSymbols(infos []lspSymbolInformation) []caf.Symbol {
	symbols := make([]caf.Symbol, 0, len(infos))
	for _, info := range infos {
		path := uriToPath(info.Location.URI)
		if path == "" {
			continue // не-file:// URI (виртуальные буферы) для графа не нужны
		}
		loc := fmt.Sprintf("%s:%d:%d", path,
			info.Location.Range.Start.Line+1,
			info.Location.Range.Start.Character+1)
		sym := caf.Symbol{
			Name:     info.Name,
			Kind:     symbolKindName(info.Kind),
			Parent:   info.ContainerName,
			Location: loc,
		}
		symbols = append(symbols, sym)
	}
	return symbols
}

// symbolKindName безопасно переводит kind в строку; неизвестные виды
// получают "symbol-N" вместо пустой строки — пустой kind в графе хуже
// загадочного, но присутствующего.
func symbolKindName(kind int) string {
	if kind > 0 && kind < len(symbolKindNames) {
		return symbolKindNames[kind]
	}
	return fmt.Sprintf("symbol-%d", kind)
}

// uriToPath конвертирует file:// URI в локальный путь.
// Windows: file:///C:/foo → C:\foo (через filepath.FromSlash).
func uriToPath(uri string) string {
	if !strings.HasPrefix(uri, "file://") {
		return ""
	}
	p := strings.TrimPrefix(uri, "file://")
	// file:///C:/x на Windows: ведущий слэш перед диском — артефакт URI.
	if len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = p[1:]
	}
	return filepath.FromSlash(p)
}
