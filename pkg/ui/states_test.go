package ui

import (
	"errors"
	"testing"
)

// ─── State enum ───────────────────────────────────────────────────────────────

func TestState_String(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateLoading, "loading"},
		{StateSuccess, "success"},
		{StateError, "error"},
		{StateEmpty, "empty"},
		{StatePartial, "partial"},
		{State(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", int(tt.state), got, tt.want)
		}
	}
}

// ─── StateMachine: NewStateMachine ────────────────────────────────────────────

func TestNewStateMachine_InitialState(t *testing.T) {
	sm := NewStateMachine()
	if sm.Current() != StateLoading {
		t.Errorf("NewStateMachine() initial state = %v, want StateLoading", sm.Current())
	}
	if sm.Error() != "" {
		t.Errorf("NewStateMachine() initial errMsg = %q, want empty", sm.Error())
	}
	if len(sm.PartialMessages()) != 0 {
		t.Errorf("NewStateMachine() initial partialMsgs = %v, want empty", sm.PartialMessages())
	}
}

// ─── StateMachine: Transition (all 5 states) ─────────────────────────────────

func TestStateMachine_Transition_ToSuccess(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: false, Data: "test"})
	if sm.Current() != StateSuccess {
		t.Errorf("after DataLoadedMsg{Empty:false} → %v, want StateSuccess", sm.Current())
	}
}

func TestStateMachine_Transition_ToEmpty(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: true})
	if sm.Current() != StateEmpty {
		t.Errorf("after DataLoadedMsg{Empty:true} → %v, want StateEmpty", sm.Current())
	}
}

func TestStateMachine_Transition_ToError(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataErrorMsg{Err: errors.New("file not found")})
	if sm.Current() != StateError {
		t.Errorf("after DataErrorMsg → %v, want StateError", sm.Current())
	}
	if sm.Error() != "file not found" {
		t.Errorf("errMsg = %q, want %q", sm.Error(), "file not found")
	}
}

func TestStateMachine_Transition_ToPartial(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(PartialDataMsg{Success: "data", Missing: []string{"widget_a", "widget_b"}})
	if sm.Current() != StatePartial {
		t.Errorf("after PartialDataMsg → %v, want StatePartial", sm.Current())
	}
	if len(sm.PartialMessages()) != 2 {
		t.Errorf("partialMsgs = %v, want 2 items", sm.PartialMessages())
	}
}

func TestStateMachine_Transition_LoadingMsg(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: false})
	sm.Transition(LoadingMsg{Reason: "refreshing"})
	if sm.Current() != StateLoading {
		t.Errorf("after LoadingMsg → %v, want StateLoading", sm.Current())
	}
}

// ─── StateMachine: enforced transition rules (§16.4) ─────────────────────────

func TestStateMachine_Enforced_LoadingToLoading_Invalid(t *testing.T) {
	// INV-PROMPT-MVP-006 §16.4: Never Loading → Loading
	sm := NewStateMachine() // starts in StateLoading
	sm.Transition(LoadingMsg{Reason: "again"})
	if sm.Current() != StateLoading {
		t.Errorf("Loading → Loading should be ignored, got %v", sm.Current())
	}
}

func TestStateMachine_Enforced_SuccessToEmpty_Invalid(t *testing.T) {
	// INV-PROMPT-MVP-006 §16.4: Never Success → Empty (data doesn't vanish)
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: false}) // → Success
	sm.Transition(DataLoadedMsg{Empty: true})  // → Empty: INVALID
	if sm.Current() != StateSuccess {
		t.Errorf("Success → Empty should be ignored, got %v", sm.Current())
	}
}

func TestStateMachine_Enforced_ErrorToSuccess_Invalid(t *testing.T) {
	// INV-PROMPT-MVP-006 §16.4: Never Error → Success without retry
	sm := NewStateMachine()
	sm.Transition(DataErrorMsg{Err: errors.New("fail")}) // → Error
	sm.Transition(DataLoadedMsg{Empty: false})           // → Success: INVALID
	if sm.Current() != StateError {
		t.Errorf("Error → Success should be ignored, got %v", sm.Current())
	}
}

func TestStateMachine_Enforced_PartialToError_Invalid(t *testing.T) {
	// INV-PROMPT-MVP-006 §16.4: Never Partial → Error
	sm := NewStateMachine()
	sm.Transition(PartialDataMsg{Success: "ok", Missing: []string{"x"}}) // → Partial
	sm.Transition(DataErrorMsg{Err: errors.New("fail")})                 // → Error: INVALID
	if sm.Current() != StatePartial {
		t.Errorf("Partial → Error should be ignored, got %v", sm.Current())
	}
}

// ─── StateMachine: valid retry/refresh paths ─────────────────────────────────

func TestStateMachine_ErrorToLoading_ViaRetry(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataErrorMsg{Err: errors.New("fail")}) // → Error
	sm.Transition(RetryMsg{Target: "dashboard"})         // → Loading (valid)
	if sm.Current() != StateLoading {
		t.Errorf("Error → Loading via RetryMsg → %v, want StateLoading", sm.Current())
	}
	if sm.Error() != "" {
		t.Errorf("errMsg should be cleared after retry, got %q", sm.Error())
	}
}

func TestStateMachine_EmptyToLoading_ViaRefresh(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: true}) // → Empty
	sm.Transition(RefreshMsg{Force: false})   // → Loading (valid)
	if sm.Current() != StateLoading {
		t.Errorf("Empty → Loading via RefreshMsg → %v, want StateLoading", sm.Current())
	}
}

func TestStateMachine_PartialToLoading_ViaRetry(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(PartialDataMsg{Success: "ok", Missing: []string{"x"}}) // → Partial
	sm.Transition(RetryMsg{Target: "widget_a"})                          // → Loading (valid)
	if sm.Current() != StateLoading {
		t.Errorf("Partial → Loading via RetryMsg → %v, want StateLoading", sm.Current())
	}
	if len(sm.PartialMessages()) != 0 {
		t.Errorf("partialMsgs should be cleared after retry, got %v", sm.PartialMessages())
	}
}

func TestStateMachine_EmptyToEmpty_ViaRefresh(t *testing.T) {
	// §16.4: Empty → Empty is allowed (still no data after refresh)
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: true}) // → Empty
	sm.Transition(RefreshMsg{Force: true})    // → Loading
	sm.Transition(DataLoadedMsg{Empty: true}) // → Empty (valid)
	if sm.Current() != StateEmpty {
		t.Errorf("Empty → Empty via refresh cycle → %v, want StateEmpty", sm.Current())
	}
}

func TestStateMachine_ErrorToError_ViaRetryFailure(t *testing.T) {
	// §16.4: Error → Error is allowed (retry failure)
	sm := NewStateMachine()
	sm.Transition(DataErrorMsg{Err: errors.New("first")}) // → Error
	sm.Transition(RetryMsg{Target: "dashboard"})          // → Loading
	sm.Transition(DataErrorMsg{Err: errors.New("second")}) // → Error (valid)
	if sm.Current() != StateError {
		t.Errorf("Error → Error via retry failure → %v, want StateError", sm.Current())
	}
	if sm.Error() != "second" {
		t.Errorf("errMsg = %q, want %q", sm.Error(), "second")
	}
}

func TestStateMachine_UnknownMessage_Ignored(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(DataLoadedMsg{Empty: false}) // → Success
	sm.Transition("some random string")        // unknown message
	if sm.Current() != StateSuccess {
		t.Errorf("unknown message should be ignored, got %v", sm.Current())
	}
}

// ─── ValidateTransition ──────────────────────────────────────────────────────

func TestValidateTransition_Valid(t *testing.T) {
	valid := []struct{ from, to State }{
		// Loading → {Success, Error, Empty, Partial}
		{StateLoading, StateSuccess},
		{StateLoading, StateError},
		{StateLoading, StateEmpty},
		{StateLoading, StatePartial},
		// Success → {Loading, Error, Partial}
		{StateSuccess, StateLoading},
		{StateSuccess, StateError},
		{StateSuccess, StatePartial},
		// Error → {Loading, Error}
		{StateError, StateLoading},
		{StateError, StateError},
		// Empty → {Loading, Success, Empty}
		{StateEmpty, StateLoading},
		{StateEmpty, StateSuccess},
		{StateEmpty, StateEmpty},
		// Partial → {Loading, Success, Partial}
		{StatePartial, StateLoading},
		{StatePartial, StateSuccess},
		{StatePartial, StatePartial},
	}
	for _, tt := range valid {
		if err := ValidateTransition(tt.from, tt.to); err != nil {
			t.Errorf("ValidateTransition(%v, %v) = %v, want nil", tt.from, tt.to, err)
		}
	}
}

func TestValidateTransition_Invalid(t *testing.T) {
	invalid := []struct{ from, to State }{
		{StateLoading, StateLoading},  // no infinite spinners
		{StateSuccess, StateEmpty},    // data doesn't vanish
		{StateSuccess, StateSuccess},  // no redundant transition
		{StateError, StateSuccess},    // never Error → Success without retry
		{StateError, StateEmpty},      // invalid
		{StateError, StatePartial},    // invalid
		{StateEmpty, StateError},      // invalid
		{StateEmpty, StatePartial},    // invalid
		{StatePartial, StateError},    // never Partial → Error
		{StatePartial, StateEmpty},    // invalid
	}
	for _, tt := range invalid {
		err := ValidateTransition(tt.from, tt.to)
		if err == nil {
			t.Errorf("ValidateTransition(%v, %v) = nil, want error", tt.from, tt.to)
		}
		var ste *StateTransitionError
		if !errors.As(err, &ste) {
			t.Errorf("ValidateTransition(%v, %v) error type = %T, want *StateTransitionError", tt.from, tt.to, err)
		}
	}
}

func TestValidateTransition_UnknownFromState(t *testing.T) {
	err := ValidateTransition(State(99), StateSuccess)
	if err == nil {
		t.Fatal("ValidateTransition with unknown from-state should return error")
	}
	var ste *StateTransitionError
	if !errors.As(err, &ste) {
		t.Fatalf("error type = %T, want *StateTransitionError", err)
	}
}

// ─── State helpers ───────────────────────────────────────────────────────────

func TestStateMachine_IsHelpers(t *testing.T) {
	sm := NewStateMachine()
	if !sm.IsLoading() {
		t.Error("IsLoading() should be true initially")
	}

	sm.Transition(DataLoadedMsg{Empty: false})
	if !sm.IsSuccess() {
		t.Error("IsSuccess() should be true after DataLoadedMsg{Empty:false}")
	}

	sm.Transition(RefreshMsg{})
	sm.Transition(DataErrorMsg{Err: errors.New("x")})
	if !sm.IsError() {
		t.Error("IsError() should be true after DataErrorMsg")
	}

	sm.Transition(RetryMsg{})
	sm.Transition(DataLoadedMsg{Empty: true})
	if !sm.IsEmpty() {
		t.Error("IsEmpty() should be true after DataLoadedMsg{Empty:true}")
	}

	sm.Transition(RefreshMsg{})
	sm.Transition(PartialDataMsg{Missing: []string{"w"}})
	if !sm.IsPartial() {
		t.Error("IsPartial() should be true after PartialDataMsg")
	}
}

// ─── ScreenBase ──────────────────────────────────────────────────────────────

func newTestScreen() *ScreenBase {
	theme := NewTheme()
	keys := NewKeys()
	a11y := NewA11y(NewAccessibilityConfig(), theme, keys)
	budgets := NewPerformanceBudget()
	perf := NewPerfTracker(budgets)
	return NewScreen(theme, keys, a11y, perf)
}

func TestScreenBase_NewScreen(t *testing.T) {
	s := newTestScreen()
	if s.StateMachine == nil {
		t.Fatal("ScreenBase.StateMachine should not be nil")
	}
	if s.Theme == nil {
		t.Fatal("ScreenBase.Theme should not be nil")
	}
	if s.Keys == nil {
		t.Fatal("ScreenBase.Keys should not be nil")
	}
	if s.ScreenReader == nil {
		t.Fatal("ScreenBase.ScreenReader should not be nil")
	}
	if s.PerfTracker == nil {
		t.Fatal("ScreenBase.PerfTracker should not be nil")
	}
	if s.GetCurrentState() != StateLoading {
		t.Errorf("new screen state = %v, want StateLoading", s.GetCurrentState())
	}
}

func TestScreenBase_Transition(t *testing.T) {
	s := newTestScreen()
	s.Transition(DataLoadedMsg{Empty: false})
	if s.GetCurrentState() != StateSuccess {
		t.Errorf("after Transition → %v, want StateSuccess", s.GetCurrentState())
	}
	if s.GetStateString() != "success" {
		t.Errorf("GetStateString() = %q, want %q", s.GetStateString(), "success")
	}
}

func TestScreenBase_IsState(t *testing.T) {
	s := newTestScreen()
	if !s.IsState(StateLoading) {
		t.Error("IsState(StateLoading) should be true initially")
	}
	s.Transition(DataLoadedMsg{Empty: false})
	if !s.IsState(StateSuccess) {
		t.Error("IsState(StateSuccess) should be true after data load")
	}
}

func TestScreenBase_GetErrorMessage(t *testing.T) {
	s := newTestScreen()
	s.Transition(DataErrorMsg{Err: errors.New("permission denied")})
	if s.GetErrorMessage() != "permission denied" {
		t.Errorf("GetErrorMessage() = %q, want %q", s.GetErrorMessage(), "permission denied")
	}
}

func TestScreenBase_GetPartialMessages(t *testing.T) {
	s := newTestScreen()
	s.Transition(PartialDataMsg{Missing: []string{"cross_repo", "doctor"}})
	msgs := s.GetPartialMessages()
	if len(msgs) != 2 {
		t.Errorf("GetPartialMessages() = %v, want 2 items", msgs)
	}
}

func TestScreenBase_GetStateDescription(t *testing.T) {
	s := newTestScreen()

	if d := s.GetStateDescription(); d != "Loading..." {
		t.Errorf("loading description = %q, want %q", d, "Loading...")
	}

	s.Transition(DataLoadedMsg{Empty: false})
	if d := s.GetStateDescription(); d != "Success" {
		t.Errorf("success description = %q, want %q", d, "Success")
	}

	s.Transition(RefreshMsg{})
	s.Transition(DataErrorMsg{Err: errors.New("boom")})
	if d := s.GetStateDescription(); d != "Error: boom" {
		t.Errorf("error description = %q, want %q", d, "Error: boom")
	}

	s.Transition(RetryMsg{})
	s.Transition(DataLoadedMsg{Empty: true})
	if d := s.GetStateDescription(); d != "Empty" {
		t.Errorf("empty description = %q, want %q", d, "Empty")
	}

	s.Transition(RefreshMsg{})
	s.Transition(PartialDataMsg{Missing: []string{"x"}})
	if d := s.GetStateDescription(); d != "Partial" {
		t.Errorf("partial description = %q, want %q", d, "Partial")
	}
}

func TestScreenBase_GetStateColor(t *testing.T) {
	s := newTestScreen()
	colors := map[State]string{
		StateLoading: "muted",
		StateSuccess: "secondary",
		StateError:   "error",
		StateEmpty:   "muted",
		StatePartial: "warning",
	}

	for state, wantColor := range colors {
		s.StateMachine.current = state // прямая установка для теста
		if got := s.GetStateColor(); got != wantColor {
			t.Errorf("GetStateColor(%v) = %q, want %q", state, got, wantColor)
		}
	}
}

func TestScreenBase_GetStateIcon(t *testing.T) {
	s := newTestScreen()
	icons := map[State]string{
		StateLoading: "⠋",
		StateSuccess: "✓",
		StateError:   "✗",
		StateEmpty:   "○",
		StatePartial: "⚠",
	}

	for state, wantIcon := range icons {
		s.StateMachine.current = state
		if got := s.GetStateIcon(); got != wantIcon {
			t.Errorf("GetStateIcon(%v) = %q, want %q", state, got, wantIcon)
		}
	}
}

// ─── StateTransitionError ────────────────────────────────────────────────────

func TestStateTransitionError_Error(t *testing.T) {
	e := &StateTransitionError{
		FromState: StateSuccess,
		ToState:   StateEmpty,
		Message:   "invalid transition",
	}
	want := "invalid transition from success to empty: invalid transition"
	if got := e.Error(); got != want {
		t.Errorf("StateTransitionError.Error() = %q, want %q", got, want)
	}
}
