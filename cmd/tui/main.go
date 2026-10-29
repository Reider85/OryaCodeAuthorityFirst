package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/Reider85/OryaCodeAuthorityFirst/pkg/ui"
)

// TUI represents the TUI application
type TUI struct {
	model *ui.AppModel
}

// NewTUI creates a new TUI application
func NewTUI() *TUI {
	return &TUI{
		model: ui.NewAppModel(),
	}
}

// Run runs the TUI application
func (t *TUI) Run() error {
	// Initialize the model
	if err := t.initialize(); err != nil {
		return fmt.Errorf("failed to initialize TUI: %w", err)
	}

	// Create the tea program
	program := bubbletea.NewProgram(t.model, bubbletea.WithAltScreen(), bubbletea.WithMouseCellMotion())

	// Run the program
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

// initialize initializes the TUI
func (t *TUI) initialize() error {
	// Initialize the model
	cmd := t.model.Init()
	if cmd != nil {
		// In a real implementation, we would run the command
		// For now, we'll just initialize the model
	}
	
	return nil
}

// main is the main entry point for the TUI
func main() {
	// Create a new TUI application
	tui := NewTUI()
	
	// Run the TUI application
	if err := tui.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}