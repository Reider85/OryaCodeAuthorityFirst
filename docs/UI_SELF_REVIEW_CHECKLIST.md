# UI Self-Review Checklist - CodeAuthorityFirst

**INV-PROMPT-MVP-010** - UI Self-Review checklist (требует ручной верификации)

*Version: 0.4.1-mvp.1*  
*Scope: MVP screens (dashboard, adapters, onboarding)*  
*Last Updated: 2026-11-06*

---

## 📋 Checklist Overview

This checklist provides manual verification procedures for CodeAuthorityFirst UI implementations. Each screen must pass all applicable checks before PR submission.

**Coverage Areas:**
- Performance requirements (INV-UI-001 through INV-UI-010)
- Accessibility compliance (WCAG-equivalent for terminal)
- Screen template compliance
- Universal keyboard map verification
- Visual design system
- UI-velocity patterns
- Testing coverage
- Code organization standards
- Edge case handling
- Documentation integration

---

## 🔍 Performance Requirements Verification

### INV-UI-001: First Paint <100ms
**Verification Steps:**
1. Build the screen: `go build -o caf ./cmd/caf`
2. Run with timing: `time -p ./caf dashboard`
3. Measure from command start to first visual output
4. **Acceptance**: Must complete in <100ms

**Common Issues:**
- Heavy data loading on screen init
- Expensive layout calculations
- Missing skeleton placeholder

### INV-UI-002: Frame Render <16ms (60fps)
**Verification Steps:**
1. Open screen in terminal
2. Trigger scroll action (j/k or ↑/↓)
3. Measure frame time with external tool
4. **Acceptance**: Scroll must be smooth, no visible lag

**Common Issues:**
- Virtualization not implemented for large lists
- Expensive string concatenation
- Missing cached layouts

### INV-UI-003: Keystroke Response <50ms
**Verification Steps:**
1. Open screen
2. Press simple keys (Tab, j, k, Enter, r)
3. Measure response time
4. **Acceptance**: Instant response, no perceptible delay

**Common Issues:**
- Async operations blocking main thread
- Expensive validation on keystroke
- Missing direct dispatch for simple actions

### INV-UI-004: Skeleton + SWR Pattern
**Verification Steps:**
1. Open screen with slow data source
2. Verify skeleton appears immediately
3. Check that real data replaces skeleton
4. **Acceptance**: No empty states, smooth data transition

**Common Issues:**
- Skeleton missing or delayed
- DataLoadedMsg not properly handled
- SWR pattern not implemented

### INV-UI-005: Lazy Screen Mounting
**Verification Steps:**
1. Navigate to screen
2. Close screen immediately after opening
3. Monitor memory and CPU usage
4. **Acceptance**: Closed screens cost zero render time

**Common Issues:**
- Screens not properly cleaned up
- Background goroutines still running
- State not properly reset

### INV-UI-006: Lipgloss Pre-computed Styles
**Verification Steps:**
1. Check `pkg/ui/theme.go` for style definitions
2. Verify styles are computed at startup
3. **Acceptance**: No style recomputation during rendering

**Common Issues:**
- Styles computed during View() calls
- Missing lipgloss.NewStyle() calls
- Dynamic style creation in hot paths

### INV-UI-007: Heatmap Render <200ms
**Verification Steps:**
1. Test with 247+ file repository
2. Measure heatmap render time
3. **Acceptance**: <200ms for initial render

**Common Issues:**
- Missing pre-computation
- Inefficient cell rendering
- No caching of layout calculations

---

## ♿ Accessibility Compliance (WCAG-equivalent for terminal)

### Perceivable: Text Representation for All Visual Elements
**Verification Steps:**
1. Check all visual elements have text equivalents
2. Verify sparklines have alt-text descriptions
3. Confirm bar charts have percentage text
4. **Acceptance**: All visual elements work without color/symbols

**Common Issues:**
- Color-only status indicators
- Symbol-only progress bars
- Missing alt-text for charts

### Operable: 100% Keyboard-Only Operation
**Verification Steps:**
1. Navigate entire screen without mouse
2. Verify all features accessible via keyboard
3. Check no mouse traps or focus issues
4. **Acceptance**: Complete keyboard coverage (INV-UI-009)

**Common Issues:**
- Mouse-only interactions
- Missing keyboard shortcuts
- Focus order issues

### Understandable: Consistent Navigation
**Verification Steps:**
1. Check universal keyboard map compliance
2. Verify help overlay (?) works on all screens
3. Confirm consistent terminology
4. **Acceptance**: Predictable navigation across screens

**Common Issues:**
- Inconsistent shortcut keys
- Missing help overlays
- Non-standard terminology

### Robust: Screen Reader Support
**Verification Steps:**
1. Test with screen reader (VoiceOver, NVDA)
2. Verify text output is readable
3. Check status-first announcement order
4. **Acceptance**: Screen reader can navigate and understand

**Common Issues:**
- Hidden text elements
- Status not announced first
- Complex markup issues

### Color-Blind Safety
**Verification Steps:**
1. Test in deuteranomaly mode
2. Verify color + symbol + text combinations
3. Check high-contrast mode works
4. **Acceptance**: All information accessible without color distinction

**Common Issues:**
- Color-only status indicators
- Poor contrast in high-contrast mode
- Missing high-contrast palette

---

## 📐 Screen Template Compliance

### Header Requirements
**Verification Steps:**
1. Check screen has proper header comment
2. Verify all sections are documented:
   - Purpose
   - Layout
   - Components
   - Keyboard map
   - Performance budget
   - Edge cases
   - Accessibility notes
3. **Acceptance**: Complete template documentation

**Common Issues:**
- Missing sections
- Incomplete documentation
- Outdated information

### Structure Requirements
**Verification Steps:**
1. Verify Init() → Update() → View() pattern
2. Check state management follows Bubble Tea ELM
3. Confirm proper tea.Cmd usage
4. **Acceptance**: Correct ELM architecture implementation

**Common Issues:**
- Missing Init() method
- Incorrect Update() implementation
- View() doing state updates

### State Management
**Verification Steps:**
1. Check ui.StateMachine usage
2. Verify state transitions are clean
3. Confirm no state leaks between screens
4. **Acceptance**: Proper state isolation and management

**Common Issues:**
- Global state pollution
- Missing state cleanup
- Incorrect state transitions

---

## ⌨️ Universal Keyboard Map Verification

### Global Keys
**Verification Steps:**
1. Verify ? (help) works on all screens
2. Check q/Ctrl+C (quit) functionality
3. Confirm Esc (back/cancel) behavior
4. **Acceptance**: All global keys work consistently

**Common Issues:**
- Inconsistent quit behavior
- Esc not working in all contexts
- Help overlay missing

### Navigation Keys
**Verification Steps:**
1. Test Tab/Shift+Tab widget cycling
2. Verify j/k/↑/↓ list navigation
3. Check g/G (top/bottom) functionality
4. **Acceptance**: Smooth, predictable navigation

**Common Issues:**
- Widget focus not cycling properly
- List navigation issues
- Missing top/bottom shortcuts

### Action Keys
**Verification Steps:**
1. Test Enter (activate) functionality
2. Verify r (refresh) behavior
3. Check 1-9 (screen switching) works
4. **Acceptance**: All actions accessible via keyboard

**Common Issues:**
- Activation not working
- Refresh not triggering
- Screen switching broken

### Global Shortcuts
**Verification Steps:**
1. Test D (doctor) from any screen
2. Verify M (metrics) accessibility
3. Check $ (cost) functionality
4. **Acceptance**: Global shortcuts work from all screens

**Common Issues:**
- Global shortcuts not registered
- Conflicts with screen-specific keys
- Missing shortcut documentation

---

## 🎨 Visual Design System

### Color Palette
**Verification Steps:**
1. Check colors defined in `pkg/ui/theme.go`
2. Verify color usage matches specification
3. Confirm fallback to 16-color mode
4. **Acceptance**: Consistent color application

**Common Issues:**
- Missing color definitions
- Inconsistent color usage
- Poor fallback support

### Typography
**Verification Steps:**
1. Verify H1 (Bold + Accent) usage
2. Check H2 (Bold + Primary) implementation
3. Confirm Body (Regular) styling
4. **Acceptance**: Consistent typography hierarchy

**Common Issues:**
- Incorrect heading styles
- Missing text styling
- Inconsistent font usage

### High-Contrast Mode
**Verification Steps:**
1. Test high-contrast mode activation
2. Verify WCAG-AAA palette compliance
3. Check all screens render correctly
4. **Acceptance**: High-contrast mode fully functional

**Common Issues:**
- Missing high-contrast palette
- Incomplete high-contrast styling
- Poor contrast ratios

---

## 🚀 UI-Velocity Patterns

### Critical CSS Analogue
**Verification Steps:**
1. Check lipgloss styles pre-computed at startup
2. Verify no style recomputation during rendering
3. **Acceptance**: Styles cached and reused efficiently

**Common Issues:**
- Styles computed during View()
- Missing lipgloss.NewStyle() calls
- Dynamic style creation

### Islands/Partial Hydration
**Verification Steps:**
1. Verify screens mount lazily
2. Check closed screens cost zero render time
3. **Acceptance**: Efficient partial hydration implementation

**Common Issues:**
- Screens not lazy-loaded
- Memory leaks from closed screens
- Expensive cleanup operations

### Skeleton/SWR Pattern
**Verification Steps:**
1. Check skeleton placeholders appear immediately
2. Verify DataLoadedMsg swaps to real data
3. **Acceptance**: Smooth data loading experience

**Common Issues:**
- Skeleton missing or delayed
- SWR pattern not implemented
- Data flickering during load

### Progressive Enhancement
**Verification Steps:**
1. Test in 16-color mode
2. Verify enhancement to 256-color
3. Check enhancement to truecolor
4. **Acceptance**: Functional at all enhancement levels

**Common Issues:**
- Missing 16-color fallback
- Inconsistent enhancement levels
- Broken progressive rendering

---

## 🧪 Testing Coverage Requirements

### Unit Tests
**Verification Steps:**
1. Check screen has *_test.go file
2. Verify snapshot tests with termenv profiles
3. Confirm Mock DataLoader interface usage
4. **Acceptance**: Complete unit test coverage

**Common Issues:**
- Missing test files
- Incomplete snapshot tests
- Missing mock interfaces

### Integration Tests
**Verification Steps:**
1. Check real .caf/ fixtures usage
2. Verify teatest integration tests
3. Confirm terminal size testing (60/100/140 cols)
4. **Acceptance**: Full integration test coverage

**Common Issues:**
- Missing fixture tests
- Incomplete teatest coverage
- No responsive layout testing

### Performance Benchmarks
**Verification Steps:**
1. Check CI benchmarks defined
2. Verify regression detection (>2x triggers alert)
3. Confirm benchmark targets are met
4. **Acceptance**: Performance enforced in CI

**Common Issues:**
- Missing benchmark definitions
- No regression detection
- Incorrect benchmark targets

### Accessibility Tests
**Verification Steps:**
1. Check focus visibility tests
2. Verify symbol-text pairing tests
3. Confirm high-contrast mode tests
4. **Acceptance**: Complete accessibility test coverage

**Common Issues:**
- Missing focus tests
- No accessibility verification
- Incomplete a11y coverage

---

## 📁 Code Organization Standards

### File Structure
**Verification Steps:**
1. Check screen in pkg/ui/screens/<name>.go
2. Verify component reuse from pkg/ui/components/
3. Confirm proper package naming
4. **Acceptance**: Correct file organization

**Common Issues:**
- Incorrect file location
- Missing component reuse
- Wrong package structure

### State Management
**Verification Steps:**
1. Check ui.StateMachine usage
2. Verify centralized state management
3. Confirm proper state isolation
4. **Acceptance**: Clean state architecture

**Common Issues:**
- Global state pollution
- Missing state management
- Poor state isolation

### Error Handling
**Verification Steps:**
1. Check consistent error states
2. Verify recovery options ([r] Retry [d] Doctor [i] Ignore [q] Quit)
3. Confirm proper error messaging
4. **Acceptance**: Robust error handling

**Common Issues:**
- Inconsistent error states
- Missing recovery options
- Poor error messaging

---

## ⚠️ Edge Case Handling

### Empty States
**Verification Steps:**
1. Check empty states have CTA (Call To Action)
2. Verify helpful error messages
3. Confirm no dead-end states
4. **Acceptance**: User always knows what to do

**Common Issues:**
- Missing CTAs in empty states
- Unhelpful error messages
- Dead-end UI states

### Loading States
**Verification Steps:**
1. Check skeleton placeholders
2. Verify progress indication
3. Confirm cancel functionality
4. **Acceptance**: Smooth loading experience

**Common Issues:**
- Missing skeleton placeholders
- No progress indication
- Loading cannot be cancelled

### Responsive Layout
**Verification Steps:**
1. Test at 60/100/140 column widths
2. Verify layout adapts properly
3. Confirm no content cutoff
4. **Acceptance**: Responsive at all sizes

**Common Issues:**
- Layout breaks at narrow widths
- Content cutoff at wide widths
- Inconsistent responsive behavior

### Error Recovery
**Verification Steps:**
1. Check error state recovery options
2. Verify graceful degradation
3. Confirm user guidance
4. **Acceptance**: Robust error recovery

**Common Issues:**
- No recovery options
- Poor error messaging
- Application crashes on errors

---

## 📚 Documentation Integration

### Screen Comments
**Verification Steps:**
1. Check header comment includes all template sections
2. Verify performance budget documentation
3. Confirm accessibility notes included
4. **Acceptance**: Complete screen documentation

**Common Issues:**
- Missing template sections
- Incomplete performance notes
- No accessibility documentation

### Performance Budgets
**Verification Steps:**
1. Check performance budgets referenced in code
2. Verify budget enforcement in tests
3. Confirm documentation matches implementation
4. **Acceptance**: Performance documented and enforced

**Common Issues:**
- Missing performance documentation
- No budget enforcement
- Documentation outdated

### Keyboard Shortcuts
**Verification Steps:**
1. Check help overlay (?) documentation
2. Verify keyboard map completeness
3. Confirm shortcut consistency
4. **Acceptance**: Complete keyboard documentation

**Common Issues:**
- Missing help documentation
- Incomplete keyboard maps
- Inconsistent shortcut documentation

---

## 🎯 Verification Workflow

### Manual Testing Procedure
1. **Screen Preparation**
   - Build the application: `go build -o caf ./cmd/caf`
   - Create test repository with sample data
   - Set up test environment

2. **Performance Testing**
   - Measure first paint time
   - Test scroll/render performance
   - Verify keystroke response
   - Check memory usage

3. **Accessibility Testing**
   - Test with screen reader
   - Verify keyboard navigation
   - Check color-blind mode
   - Test high-contrast mode

4. **Functional Testing**
   - Test all keyboard shortcuts
   - Verify edge cases
   - Check error handling
   - Test responsive layout

5. **Documentation Review**
   - Check code comments
   - Verify help documentation
   - Confirm performance budgets
   - Review accessibility notes

### Common Pitfalls
- **Performance**: Heavy data loading on screen init
- **Accessibility**: Color-only status indicators
- **Keyboard**: Missing global shortcuts
- **Layout**: Poor responsive behavior
- **State**: Memory leaks from closed screens
- **Testing**: Missing snapshot tests
- **Documentation**: Outdated performance notes

### Troubleshooting Guide
- **Slow Performance**: Check for missing skeleton, virtualization, or cached layouts
- **Keyboard Issues**: Verify universal keyboard map and global shortcuts
- **Accessibility Problems**: Check text representation and color-blind safety
- **Layout Issues**: Test responsive behavior at different terminal widths
- **State Problems**: Check for memory leaks and proper cleanup

---

## ✅ Final Checklist

Before submitting a PR, verify:

- [ ] All performance requirements met (INV-UI-001 through INV-UI-010)
- [ ] Accessibility compliance verified (WCAG-equivalent)
- [ ] Screen template fully implemented
- [ ] Universal keyboard map working
- [ ] Visual design system applied consistently
- [ ] UI-velocity patterns implemented
- [ ] Testing coverage complete
- [ ] Code organization standards followed
- [ ] Edge cases handled properly
- [ ] Documentation integrated and updated

**Manual Verification Required**: This checklist must be completed manually for each screen change. Automated tests cannot fully verify UI quality and user experience.

---

*This checklist is part of CodeAuthorityFirst's MVP scope (INV-PROMPT-MVP-010) and must be followed for all UI implementations.*