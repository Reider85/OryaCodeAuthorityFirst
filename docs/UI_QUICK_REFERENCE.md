# UI Self-Review Quick Reference

**Cheat sheet for CodeAuthorityFirst UI development**

*Version: 0.4.1-mvp.1*  
*Quick reference for daily UI development*

---

## 🚀 Performance Requirements (Must Pass)

| Requirement | Target | Test |
|-------------|--------|------|
| **First paint** | <100ms | Time from init to first output |
| **Frame render** | <16ms | Scroll/redraw performance |
| **Keystroke response** | <50ms | Simple action response time |
| **Widget render** | <50ms | Individual component rendering |
| **Memory usage** | <50MB | Peak memory consumption |

**Quick Checks:**
- [ ] Skeleton appears immediately
- [ ] No blocking operations on init
- [ ] Virtualized lists for >100 items
- [ ] Cached layouts and styles
- [ ] Pre-computed lipgloss styles

---

## ♿ Accessibility Requirements (Must Pass)

| Requirement | Test |
|-------------|------|
| **Focus visibility** | ▸ marker shows focus |
| **Symbol-text pairing** | ✓⚠✗ always with text |
| **Keyboard-only** | 100% coverage without mouse |
| **Help overlay** | ? shows shortcuts |
| **High-contrast** | Works in WCAG-AAA mode |

**Quick Checks:**
- [ ] All visual elements have text representation
- [ ] No color-only indicators
- [ ] Tab navigation works everywhere
- [ ] Help overlay available on all screens
- [ ] High-contrast mode functional

---

## ⌨️ Universal Keyboard Map (Memorize)

### Global Keys
- `?` - Help overlay
- `q/Ctrl+C` - Quit
- `Esc` - Back/Cancel

### Navigation
- `Tab/Shift+Tab` - Cycle widgets
- `j/k/↑/↓` - Move in lists
- `g/G` - Top/Bottom

### Actions
- `Enter` - Activate/Select
- `Space` - Toggle
- `r` - Refresh

### Global Shortcuts
- `D` - Doctor view
- `M` - Metrics dashboard
- `$` - Cost dashboard

---

## 🎨 Visual Design System

### Colors
```go
Primary:    #61AFEF (blue)
Secondary:  #98C379 (green) 
Warning:    #E5C07B (yellow)
Error:      #E06C75 (red)
Muted:      #5C6370 (gray)
Accent:     #C678DD (purple)
Background: #282C34 (dark)
```

### Typography
- **H1**: `# Title` (Bold + Accent)
- **H2**: `## Section` (Bold + Primary)
- **Body**: Regular text
- **Code**: `inline code` (SarasaMonoSC)

### High Contrast
```go
Primary:    #FFFFFF (white on black)
Secondary:  #00FF00 (bright green)
Warning:    #FFFF00 (bright yellow)
Error:      #FF0000 (bright red)
```

---

## 📐 Screen Template (Follow Exactly)

```go
// pkg/ui/screens/<name>.go
// §<N> <Title> (caf <command>)
//
// Purpose: <one sentence>
// Layout: <ASCII mockup>
// Components: <list>
// Keyboard map: <table>
// Performance budget: <table>
// Edge cases: <list>
// Accessibility: <notes>

package screens

type <Name>Screen struct {
    ui.StateMachine
    
    // Data
    Data       []DataType
    Loading    bool
    Error      error
    
    // Focus
    Current    int
    ShowHelp   bool
    
    // Performance
    Cache      *WidgetCache
}

func (s *<Name>Screen) Init() tea.Cmd {
    return tea.Batch(
        s.loadData(),
        components.ShowSkeleton(),
    )
}

func (s *<Name>Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return s.handleKey(msg)
    case DataLoadedMsg:
        s.Data = msg.Data
        s.Loading = false
        return s, nil
    }
    return s, nil
}

func (s *<Name>Screen) View() string {
    if s.Loading {
        return components.SkeletonView()
    }
    
    var builder strings.Builder
    builder.WriteString(s.renderHeader())
    builder.WriteString(s.renderContent())
    builder.WriteString(s.renderFooter())
    
    return builder.String()
}
```

---

## 🚀 UI-Velocity Patterns (Implement)

### 1. Skeleton + SWR Pattern
```go
// Init shows skeleton immediately
func (s *Screen) Init() tea.Cmd {
    return tea.Batch(
        s.loadData(),  // Async
        components.ShowSkeleton(),
    )
}

// Update swaps skeleton with real data
func (s *Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case DataLoadedMsg:
        s.Data = msg.Data
        s.Loading = false
        return s, nil
    }
    return s, nil
}
```

### 2. Virtualized Lists
```go
// Only render visible items
func (s *Screen) renderList() string {
    const visible = 10
    start := max(0, s.Current-visible/2)
    end := min(len(s.Items), start+visible)
    
    for i := start; i < end; i++ {
        // Render only visible items
    }
}
```

### 3. Cached Styles
```go
// Pre-compute styles at startup
func init() {
    ui.HeaderStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#61AFEF"))
}

// Use cached styles
func (s *Screen) View() string {
    return ui.HeaderStyle.Render("Title")
}
```

---

## 🧪 Testing Checklist (Run Tests)

### Performance Tests
```bash
# First paint test
go test -run TestFirstPaint -v

# Frame render test  
go test -run TestFrameRender -v

# Memory usage test
go test -run TestMemoryUsage -v
```

### Accessibility Tests
```bash
# Focus visibility test
go test -run TestFocusVisibility -v

# Symbol-text pairing test
go test -run TestSymbolTextPairing -v

# High-contrast test
go test -run TestHighContrast -v
```

### Integration Tests
```bash
# Full screen flow test
go test -run TestScreenIntegration -v

# Keyboard navigation test
go test -run TestKeyboardNavigation -v
```

---

## 🎯 Daily Development Workflow

### 1. Before Coding
- [ ] Review UI_SELF_REVIEW_CHECKLIST.md
- [ ] Check UI_IMPLEMENTATION_EXAMPLES.md
- [ ] Verify performance budgets
- [ ] Review accessibility requirements

### 2. During Coding
- [ ] Follow screen template exactly
- [ ] Implement skeleton + SWR pattern
- [ ] Use universal keyboard map
- [ ] Apply visual design system
- [ ] Add performance optimizations

### 3. After Coding
- [ ] Run performance tests
- [ ] Run accessibility tests
- [ ] Test keyboard navigation
- [ ] Check edge cases
- [ ] Update documentation

### 4. Before PR
- [ ] Complete final checklist
- [ ] Test on different terminal sizes
- [ ] Verify high-contrast mode
- [ ] Check screen reader support
- [ ] Review performance benchmarks

---

## ⚠️ Common Pitfalls (Avoid)

### Performance Issues
- ❌ Heavy data loading on init
- ❌ No virtualization for large lists
- ❌ Expensive string concatenation
- ❌ Missing skeleton placeholders

### Accessibility Issues
- ❌ Color-only status indicators
- ❌ Missing keyboard shortcuts
- ❌ No help overlay
- ❌ Poor focus visibility

### Code Issues
- ❌ Missing screen template sections
- ❌ Inconsistent keyboard shortcuts
- ❌ No performance documentation
- ❌ Missing error handling

### Testing Issues
- ❌ No performance tests
- ❌ Missing accessibility tests
- ❌ No integration tests
- ❌ No edge case testing

---

## 📞 Quick Help

### When Stuck...
1. **Check examples**: `docs/UI_IMPLEMENTATION_EXAMPLES.md`
2. **Review checklist**: `docs/UI_SELF_REVIEW_CHECKLIST.md`
3. **Run tests**: `go test ./pkg/ui/screens/...`
4. **Ask team**: Reference specific INV-UI requirement

### Performance Problems
- Use `time.Now()` to measure
- Check for blocking operations
- Implement caching strategies
- Use virtualization for large data

### Accessibility Problems
- Test with screen reader
- Verify keyboard navigation
- Check color-blind mode
- Ensure text representation

### Code Organization
- Follow screen template exactly
- Use existing components
- Document performance budgets
- Test all edge cases

---

## ✅ Final Verification

Before submitting PR, confirm:

- [ ] All performance requirements met
- [ ] Accessibility compliance verified
- [ ] Universal keyboard map working
- [ ] Visual design system applied
- [ ] UI-velocity patterns implemented
- [ ] Tests passing
- [ ] Documentation updated
- [ ] Edge cases handled

**Manual Review Required**: This checklist must be completed manually for each UI change.

---

*Quick reference for CodeAuthorityFirst UI development - Keep this handy while coding!*