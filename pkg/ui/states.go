package ui

import (
	"fmt"
	"time"
)

// State represents the UI state
// INV-PROMPT-MVP-006: Every screen implements ALL 5 mandatory UI states
type State int

const (
	StateLoading State = iota
	StateSuccess
	StateError
	StateEmpty
	StatePartial
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateLoading:
		return "loading"
	case StateSuccess:
		return "success"
	case StateError:
		return "error"
	case StateEmpty:
		return "empty"
	case StatePartial:
		return "partial"
	default:
		return "unknown"
	}
}

// StateMachine governs transitions for every screen
// INV-PROMPT-MVP-006: Every screen embeds StateMachine and calls Transition() on data events
type StateMachine struct {
	current State
	errMsg  string
	partialMsgs []string  // for StatePartial — what's missing
}

// NewStateMachine creates a new state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		current: StateLoading,
		errMsg:  "",
		partialMsgs: []string{},
	}
}

// Current returns the current state
func (sm *StateMachine) Current() State {
	return sm.current
}

// Error returns the error message if in error state
func (sm *StateMachine) Error() string {
	return sm.errMsg
}

// PartialMessages returns the partial state messages
func (sm *StateMachine) PartialMessages() []string {
	return sm.partialMsgs
}

// Transition transitions to a new state based on a message
// INV-PROMPT-MVP-006: Loading → {Success, Error, Empty, Partial}. Never Loading → Loading.
// INV-PROMPT-MVP-006: Success → {Loading (on refresh), Error (on refresh failure), Partial (on partial refresh failure)}.
// INV-PROMPT-MVP-006: Error → {Loading (on retry), Error (on retry failure)}.
// INV-PROMPT-MVP-006: Empty → {Loading (on refresh), Success (on data arrival)}.
// INV-PROMPT-MVP-006: Partial → {Loading (on retry), Success (on full refresh), Partial (on retry of failed widget only)}.
func (sm *StateMachine) Transition(msg interface{}) {
	switch msg := msg.(type) {
	case LoadingMsg:
		sm.current = StateLoading
	case DataLoadedMsg:
		if msg.Empty {
			sm.current = StateEmpty
		} else {
			sm.current = StateSuccess
		}
	case DataErrorMsg:
		sm.current = StateError
		sm.errMsg = msg.Err.Error()
	case PartialDataMsg:
		sm.current = StatePartial
		sm.partialMsgs = msg.Missing
	case RetryMsg:
		// Retry transitions depend on current state
		switch sm.current {
		case StateError:
			sm.current = StateLoading
		case StatePartial:
			sm.current = StateLoading
		case StateEmpty:
			sm.current = StateLoading
		}
	case RefreshMsg:
		// Refresh transitions
		switch sm.current {
		case StateSuccess:
			sm.current = StateLoading
		case StateError:
			sm.current = StateLoading
		case StatePartial:
			sm.current = StateLoading
		case StateEmpty:
			sm.current = StateLoading
		}
	}
}

// IsLoading returns true if the current state is loading
func (sm *StateMachine) IsLoading() bool {
	return sm.current == StateLoading
}

// IsSuccess returns true if the current state is success
func (sm *StateMachine) IsSuccess() bool {
	return sm.current == StateSuccess
}

// IsError returns true if the current state is error
func (sm *StateMachine) IsError() bool {
	return sm.current == StateError
}

// IsEmpty returns true if the current state is empty
func (sm *StateMachine) IsEmpty() bool {
	return sm.current == StateEmpty
}

// IsPartial returns true if the current state is partial
func (sm *StateMachine) IsPartial() bool {
	return sm.current == StatePartial
}

// StateMessages represents messages that can trigger state transitions
type StateMessages interface {
	LoadingMsg() LoadingMsg
	DataLoadedMsg(DataLoadedMsg) 
	DataErrorMsg(DataErrorMsg)
	PartialDataMsg(PartialDataMsg)
	RetryMsg() RetryMsg
	RefreshMsg() RefreshMsg
}

// LoadingMsg represents a loading message
type LoadingMsg struct {
	Reason string
}

// DataLoadedMsg represents a data loaded message
type DataLoadedMsg struct {
	Empty bool
	Data  interface{}
}

// DataErrorMsg represents a data error message
type DataErrorMsg struct {
	Err error
}

// PartialDataMsg represents a partial data message
type PartialDataMsg struct {
	Success interface{}
	Missing []string
}

// RetryMsg represents a retry message
type RetryMsg struct {
	Target string
}

// RefreshMsg represents a refresh message
type RefreshMsg struct {
	Force bool
}

// Screen represents a screen with state management
type Screen struct {
	StateMachine *StateMachine
	PerfTracker  *PerfTracker
	Theme        *Theme
	Keys         *Keys
}

// NewScreen creates a new screen
func NewScreen(theme *Theme, keys *Keys, perfTracker *PerfTracker) *Screen {
	return &Screen{
		StateMachine: NewStateMachine(),
		PerfTracker:  perfTracker,
		Theme:        theme,
		Keys:         keys,
	}
}

// Transition transitions the screen state
func (s *Screen) Transition(msg interface{}) {
	s.StateMachine.Transition(msg)
}

// GetCurrentState returns the current state
func (s *Screen) GetCurrentState() State {
	return s.StateMachine.Current()
}

// GetStateString returns the current state as string
func (s *Screen) GetStateString() string {
	return s.StateMachine.Current().String()
}

// IsState returns true if the screen is in the specified state
func (s *Screen) IsState(state State) bool {
	return s.StateMachine.Current() == state
}

// GetErrorMessage returns the error message if in error state
func (s *Screen) GetErrorMessage() string {
	return s.StateMachine.Error()
}

// GetPartialMessages returns the partial state messages
func (s *Screen) GetPartialMessages() []string {
	return s.StateMachine.PartialMessages()
}

// StateTransitionError represents an error in state transitions
type StateTransitionError struct {
	FromState State
	ToState   State
	Message   string
}

func (e *StateTransitionError) Error() string {
	return fmt.Sprintf("invalid transition from %s to %s: %s", e.FromState, e.ToState, e.Message)
}

// ValidateTransition validates a state transition
func ValidateTransition(from, to State) error {
	// Define valid transitions
	validTransitions := map[State][]State{
		StateLoading:  {StateSuccess, StateError, StateEmpty, StatePartial},
		StateSuccess: {StateLoading, StateError, StatePartial},
		StateError:   {StateLoading},
		StateEmpty:   {StateLoading, StateSuccess},
		StatePartial: {StateLoading, StateSuccess, StatePartial},
	}
	
	// Check if transition is valid
	validToStates, exists := validTransitions[from]
	if !exists {
		return &StateTransitionError{
			FromState: from,
			ToState:   to,
			Message:   "unknown from state",
		}
	}
	
	for _, validTo := range validToStates {
		if to == validTo {
			return nil
		}
	}
	
	return &StateTransitionError{
		FromState: from,
		ToState:   to,
		Message:   "invalid transition",
	}
}

// GetStateDescription returns a description of the current state
func (s *Screen) GetStateDescription() string {
	switch s.StateMachine.Current() {
	case StateLoading:
		return "Loading..."
	case StateSuccess:
		return "Success"
	case StateError:
		return fmt.Sprintf("Error: %s", s.StateMachine.Error())
	case StateEmpty:
		return "Empty"
	case StatePartial:
		return "Partial"
	default:
		return "Unknown"
	}
}

// GetStateColor returns the color for the current state
func (s *Screen) GetStateColor() string {
	switch s.StateMachine.Current() {
	case StateLoading:
		return "muted"
	case StateSuccess:
		return "secondary"
	case StateError:
		return "error"
	case StateEmpty:
		return "muted"
	case StatePartial:
		return "warning"
	default:
		return "muted"
	}
}

// GetStateIcon returns the icon for the current state
func (s *Screen) GetStateIcon() string {
	switch s.StateMachine.Current() {
	case StateLoading:
		return "⠋"
	case StateSuccess:
		return "✓"
	case StateError:
		return "✗"
	case StateEmpty:
		return "○"
	case StatePartial:
		return "⚠"
	default:
		return "?"
	}
}