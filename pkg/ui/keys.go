package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
)

// Keys defines the universal keyboard map
// INV-UI-009: Keyboard-only operation is 100% complete; no feature requires mouse
// UI.md §0.6: Universal keyboard map
type Keys struct {
	// Navigation
	Tab        string
	ShiftTab   string
	Enter      string
	Esc        string
	Space      string
	Up         string
	Down       string
	Left       string
	Right      string
	PageUp     string
	PageDown   string
	Home       string
	End        string
	
	// Actions
	Refresh    string
	Quit       string
	Help       string
	Back       string
	
	// Screen-specific
	Dashboard  string
	Doctor     string
	Metrics    string
	Cost       string
	
	// Universal shortcuts (global)
	Goto1      string
	Goto2      string
	Goto3      string
	Goto4      string
	Goto5      string
	Goto6      string
	Goto7      string
	Goto8      string
	Goto9      string
	
	// Movement
	Next       string
	Prev       string
	Top        string
	Bottom     string
	
	// Search
	Search     string
	CancelSearch string
	
	// Selection
	Select     string
	Deselect   string
	Toggle     string
}

// NewKeys creates a new key mapping
func NewKeys() *Keys {
	return &Keys{
		// Navigation
		Tab:        "tab",
		ShiftTab:   "shift+tab",
		Enter:      "enter",
		Esc:        "esc",
		Space:      "space",
		Up:         "up",
		Down:       "down",
		Left:       "left",
		Right:      "right",
		PageUp:     "pageup",
		PageDown:   "pagedown",
		Home:       "home",
		End:        "end",
		
		// Actions
		Refresh:   "r",
		Quit:       "q",
		Help:       "?",
		Back:       "esc",
		
		// Screen-specific
		Dashboard:  "1",
		Doctor:     "d",
		Metrics:    "m",
		Cost:       "$",
		
		// Universal shortcuts (global)
		Goto1:      "1",
		Goto2:      "2",
		Goto3:      "3",
		Goto4:      "4",
		Goto5:      "5",
		Goto6:      "6",
		Goto7:      "7",
		Goto8:      "8",
		Goto9:      "9",
		
		// Movement
		Next:      "j",
		Prev:      "k",
		Top:       "g",
		Bottom:    "G",
		
		// Search
		Search:    "/",
		CancelSearch: "esc",
		
		// Selection
		Select:    "enter",
		Deselect:  "esc",
		Toggle:    "space",
	}
}

// KeyMsg represents a key press event
type KeyMsg struct {
	Key string
}

// IsNavigationKey returns true if the key is for navigation
func (k *Keys) IsNavigationKey(key string) bool {
	switch key {
	case k.Tab, k.ShiftTab, k.Up, k.Down, k.Left, k.Right, k.PageUp, k.PageDown, k.Home, k.End:
		return true
	default:
		return false
	}
}

// IsActionKey returns true if the key is for actions
func (k *Keys) IsActionKey(key string) bool {
	switch key {
	case k.Enter, k.Space, k.Refresh, k.Quit, k.Help, k.Back:
		return true
	default:
		return false
	}
}

// IsGlobalShortcut returns true if the key is a global shortcut
func (k *Keys) IsGlobalShortcut(key string) bool {
	switch key {
	case k.Dashboard, k.Doctor, k.Metrics, k.Cost, k.Goto1, k.Goto2, k.Goto3, k.Goto4, k.Goto5, k.Goto6, k.Goto7, k.Goto8, k.Goto9:
		return true
	default:
		return false
	}
}

// ParseKey parses a key string into a KeyMsg
func ParseKey(key string) KeyMsg {
	return KeyMsg{Key: key}
}

// HandleKeyPress handles a key press and returns appropriate tea.Cmd
func HandleKeyPress(key string, model interface{}) tea.Cmd {
	// This would be implemented in each screen model
	// For now, return a no-op command
	return nil
}

// KeyConflictError represents a key conflict
type KeyConflictError struct {
	Key        string
	Current    string
	Requested  string
}

func (e *KeyConflictError) Error() string {
	return fmt.Sprintf("Key conflict: '%s' is already bound to '%s', cannot bind to '%s'", e.Key, e.Current, e.Requested)
}

// ValidateKeyMap validates that there are no conflicts in the key map
func ValidateKeyMap(keys *Keys) error {
	// Check for conflicts in global shortcuts
	globalShortcuts := map[string]string{
		keys.Dashboard:  "dashboard",
		keys.Doctor:     "doctor",
		keys.Metrics:    "metrics",
		keys.Cost:       "cost",
		keys.Goto1:      "screen1",
		keys.Goto2:      "screen2",
		keys.Goto3:      "screen3",
		keys.Goto4:      "screen4",
		keys.Goto5:      "screen5",
		keys.Goto6:      "screen6",
		keys.Goto7:      "screen7",
		keys.Goto8:      "screen8",
		keys.Goto9:      "screen9",
	}
	
	// Check for conflicts
	keyBindings := make(map[string]string)
	for key, action := range globalShortcuts {
		if existing, exists := keyBindings[key]; exists {
			return &KeyConflictError{
				Key:       key,
				Current:   existing,
				Requested: action,
			}
		}
		keyBindings[key] = action
	}
	
	return nil
}

// GetHelpText returns help text for the current key bindings
func GetHelpText(keys *Keys) string {
	return `Navigation:
  Tab/Shift+Tab    Cycle widgets
  j/k              Up/Down in lists
  g/G              Top/Bottom
  Enter            Activate/Select
  Esc              Back/Cancel

Actions:
  r                Refresh
  q                Quit
  ?                Help

Shortcuts:
  1-9              Switch screen
  d                Doctor view
  m                Metrics dashboard
  $                Cost dashboard

Search:
  /                Search
  Esc              Cancel search`
}