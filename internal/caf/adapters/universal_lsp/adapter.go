package adapters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/Reider85/OryaCodeAuthorityFirst/internal/caf"
)

// UniversalLSPAdapter implements the Adapter interface for universal LSP fallback
type UniversalLSPAdapter struct {
	lspProcess *os.Process
	conn       net.Conn
}

// LSPRequest represents a JSON-RPC request to the LSP server
type LSPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// LSPResponse represents a JSON-RPC response from the LSP server
type LSPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *LSPError   `json:"error,omitempty"`
}

// LSPError represents an LSP error
type LSPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// InitializeParams represents LSP initialize parameters
type InitializeParams struct {
	ProcessID        int                    `json:"processId"`
	RootPath        string                 `json:"rootPath"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	Trace           string                 `json:"trace"`
}

// SymbolInformation represents a symbol from the LSP
type SymbolInformation struct {
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	ContainerName string `json:"containerName,omitempty"`
	Location      Location `json:"location"`
}

// Location represents a location in the code
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// Range represents a range in the code
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in the code
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// ID returns the unique identifier for this adapter
func (a *UniversalLSPAdapter) ID() string {
	return "universal-lsp"
}

// Capabilities returns the capability matrix for this adapter
func (a *UniversalLSPAdapter) Capabilities() caf.CapabilityMatrix {
	return caf.CapabilityMatrix{
		PreSession:  caf.TierWrapper,
		PostEdit:    caf.TierWatcher,
		PreCommit:   caf.TierUniversal,
		MCP:         false,
		Tier:        4,
	}
}

// Install installs the universal LSP adapter (gopls for Go)
func (a *UniversalLSPAdapter) Install() error {
	// Check if gopls is available
	if _, err := exec.LookPath("gopls"); err != nil {
		return fmt.Errorf("gopls not found, please install Go and gopls: %w", err)
	}

	return nil
}

// Uninstall removes the universal LSP adapter
func (a *UniversalLSPAdapter) Uninstall() error {
	// Stop the LSP process if running
	if a.lspProcess != nil {
		a.lspProcess.Kill()
		a.lspProcess = nil
	}

	// Close the connection
	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	return nil
}

// Detect returns true if gopls is available on the system
func (a *UniversalLSPAdapter) Detect() bool {
	if _, err := exec.LookPath("gopls"); err != nil {
		return false
	}

	return true
}

// StartLSP starts the LSP server
func (a *UniversalLSPAdapter) StartLSP(repoRoot string) error {
	// Stop any existing LSP process
	if a.lspProcess != nil {
		a.lspProcess.Kill()
	}

	// Start gopls
	cmd := exec.Command("gopls", "-rpc.address=127.0.0.1:0")
	cmd.Dir = repoRoot

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gopls: %w", err)
	}

	a.lspProcess = cmd.Process

	// Read the address from stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "serving RPC on") {
				// Extract address from line
				parts := strings.Split(line, " ")
				if len(parts) >= 5 {
					address := parts[4]
					conn, err := net.Dial("tcp", address)
					if err == nil {
						a.conn = conn
					}
				}
			}
		}
	}()

	// Initialize the LSP
	if err := a.InitializeLSP(stdin, stdout, repoRoot); err != nil {
		return fmt.Errorf("failed to initialize LSP: %w", err)
	}

	return nil
}

// InitializeLSP initializes the LSP server
func (a *UniversalLSPAdapter) InitializeLSP(stdin io.Writer, stdout io.Reader, repoRoot string) error {
	// Send initialize request
	initReq := LSPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: InitializeParams{
			ProcessID: os.Getpid(),
			RootPath:  repoRoot,
			Capabilities: map[string]interface{}{
				"textDocument": map[string]interface{}{
					"synchronization": map[string]interface{}{
						"didSave": true,
					},
				},
			},
			Trace: "off",
		},
	}

	initData, err := json.Marshal(initReq)
	if err != nil {
		return fmt.Errorf("failed to marshal initialize request: %w", err)
	}

	if _, err := stdin.Write(initData); err != nil {
		return fmt.Errorf("failed to send initialize request: %w", err)
	}

	// Send initialized notification
	initNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
		"params":  map[string]interface{}{},
	}

	initNotifData, err := json.Marshal(initNotif)
	if err != nil {
		return fmt.Errorf("failed to marshal initialized notification: %w", err)
	}

	if _, err := stdin.Write(initNotifData); err != nil {
		return fmt.Errorf("failed to send initialized notification: %w", err)
	}

	return nil
}

// GetSymbols gets symbols from the LSP server
func (a *UniversalLSPAdapter) GetSymbols() ([]caf.Symbol, error) {
	if a.conn == nil {
		return nil, fmt.Errorf("LSP server not connected")
	}

	// Send workspace/symbols request
	symbolsReq := LSPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "workspace/symbol",
		Params: map[string]interface{}{
			"query": "",
		},
	}

	symbolsData, err := json.Marshal(symbolsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symbols request: %w", err)
	}

	if _, err := a.conn.Write(symbolsData); err != nil {
		return nil, fmt.Errorf("failed to send symbols request: %w", err)
	}

	// Read response
	buffer := make([]byte, 4096)
	n, err := a.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read symbols response: %w", err)
	}

	var symbolsResp LSPResponse
	if err := json.Unmarshal(buffer[:n], &symbolsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal symbols response: %w", err)
	}

	if symbolsResp.Error != nil {
		return nil, fmt.Errorf("LSP error: %s", symbolsResp.Error.Message)
	}

	// Convert LSP symbols to our format
	var symbols []caf.Symbol
	if symbolList, ok := symbolsResp.Result.([]interface{}); ok {
		for _, symbol := range symbolList {
			if symbolMap, ok := symbol.(map[string]interface{}); ok {
				name, _ := symbolMap["name"].(string)
				kind, _ := symbolMap["kind"].(string)
				location, _ := symbolMap["location"].(map[string]interface{})
				
				if location != nil {
					uri, _ := location["uri"].(string)
					rangeMap, _ := location["range"].(map[string]interface{})
					
					if rangeMap != nil {
						start, _ := rangeMap["start"].(map[string]interface{})
						end, _ := rangeMap["end"].(map[string]interface{})
						
						if start != nil && end != nil {
							symbol := caf.Symbol{
								Name: name,
								Kind: kind,
								Location: uri + ":" + fmt.Sprintf("%d:%d", 
									start["line"], start["character"]),
							}
							symbols = append(symbols, symbol)
						}
					}
				}
			}
		}
	}

	return symbols, nil
}

// Stop stops the LSP server
func (a *UniversalLSPAdapter) Stop() error {
	if a.lspProcess != nil {
		a.lspProcess.Kill()
		a.lspProcess = nil
	}

	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	return nil
}

// Initialize initializes the universal LSP adapter
func Initialize() {
	adapter := &UniversalLSPAdapter{}
	caf.RegisterAdapter(adapter)
}