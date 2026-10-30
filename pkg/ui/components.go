package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Component represents a UI component
type Component interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Component, tea.Cmd)
	View() string
	GetID() string
	Focus() tea.Cmd
	Blur() tea.Cmd
	IsFocused() bool
}

// BaseComponent provides common functionality for all components
type BaseComponent struct {
	ID      string
	Focused bool
	Theme   *Theme
	Keys    *Keys
}

// NewBaseComponent creates a new base component
func NewBaseComponent(id string, theme *Theme, keys *Keys) *BaseComponent {
	return &BaseComponent{
		ID:     id,
		Theme:  theme,
		Keys:   keys,
	}
}

// GetID returns the component ID
func (bc *BaseComponent) GetID() string {
	return bc.ID
}

// IsFocused returns true if the component is focused
func (bc *BaseComponent) IsFocused() bool {
	return bc.Focused
}

// SetFocused sets the focused state
func (bc *BaseComponent) SetFocused(focused bool) {
	bc.Focused = focused
}

// Focus focuses the component (INV-PROMPT-MVP-007: visible focus)
func (bc *BaseComponent) Focus() tea.Cmd {
	bc.Focused = true
	return nil
}

// Blur removes focus from the component
func (bc *BaseComponent) Blur() tea.Cmd {
	bc.Focused = false
	return nil
}

// Button represents a button component
type Button struct {
	BaseComponent
	Label     string
	Action    func() tea.Cmd
	Disabled  bool
	Hover     bool
}

// NewButton creates a new button
func NewButton(id, label string, action func() tea.Cmd, theme *Theme, keys *Keys) *Button {
	return &Button{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Label:        label,
		Action:       action,
		Disabled:     false,
		Hover:        false,
	}
}

// Init initializes the button
func (b *Button) Init() tea.Cmd {
	return nil
}

// Update updates the button
func (b *Button) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !b.Focused || b.Disabled {
			return b, nil
		}
		
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			if b.Action != nil {
				return b, b.Action()
			}
		case tea.KeyTab, tea.KeyShiftTab:
			return b, nil // Let parent handle tab navigation
		}
	}
	
	return b, nil
}

// View returns the button view
func (b *Button) View() string {
	if b.Disabled {
		return b.Theme.Muted.Render(fmt.Sprintf("[ %s ]", b.Label))
	}
	
	if b.Focused {
		return b.Theme.Button.Render(fmt.Sprintf("[ %s ]", b.Label))
	}
	
	if b.Hover {
		return b.Theme.Secondary.Render(fmt.Sprintf("[ %s ]", b.Label))
	}
	
	return fmt.Sprintf("[ %s ]", b.Label)
}

// Input represents an input field component
type Input struct {
	BaseComponent
	Value     string
	Placeholder string
	Password  bool
	MaxLength int
	Err       error
}

// NewInput creates a new input field
func NewInput(id, placeholder string, theme *Theme, keys *Keys) *Input {
	return &Input{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Value:         "",
		Placeholder:  placeholder,
		Password:     false,
		MaxLength:    256,
		Err:          nil,
	}
}

// Init initializes the input
func (i *Input) Init() tea.Cmd {
	return nil
}

// Update updates the input
func (i *Input) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !i.Focused {
		return i, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return i, nil // Submit handled by parent
		case tea.KeyEsc:
			i.Value = ""
			return i, nil
		case tea.KeyBackspace:
			if len(i.Value) > 0 {
				i.Value = i.Value[:len(i.Value)-1]
			}
			return i, nil
		case tea.KeyTab, tea.KeyShiftTab:
			return i, nil // Let parent handle tab navigation
		}
		
		// Handle character input
		if msg.Type >= tea.KeySpace && msg.Type <= tea.KeyDelete {
			if len(i.Value) < i.MaxLength {
				i.Value += msg.String()
			}
		}
	}
	
	return i, nil
}

// View returns the input view
func (i *Input) View() string {
	var displayValue string
	if i.Password {
		displayValue = strings.Repeat("•", len(i.Value))
	} else {
		displayValue = i.Value
	}
	
	if i.Err != nil {
		return fmt.Sprintf("%s %s", i.Theme.Error.Render(i.Err.Error()), i.Theme.Input.Render(displayValue))
	}
	
	if len(displayValue) == 0 {
		return fmt.Sprintf("%s %s", i.Theme.Muted.Render(i.Placeholder), i.Theme.Input.Render(displayValue))
	}
	
	if i.Focused {
		return fmt.Sprintf("> %s", i.Theme.Input.Render(displayValue))
	}
	
	return fmt.Sprintf("  %s", i.Theme.Input.Render(displayValue))
}

// List represents a list component
type List struct {
	BaseComponent
	Items     []ListItem
	Selected  int
	Height    int
	Title     string
	EmptyText string
}

// ListItem represents a list item
type ListItem struct {
	ID       string
	Title    string
	Subtitle string
	Data     interface{}
}

// NewList creates a new list
func NewList(id string, height int, theme *Theme, keys *Keys) *List {
	return &List{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Items:         []ListItem{},
		Selected:     0,
		Height:       height,
		Title:        "",
		EmptyText:    "No items",
	}
}

// Init initializes the list
func (l *List) Init() tea.Cmd {
	return nil
}

// Update updates the list
func (l *List) Update(msg tea.Msg) (Component, tea.Cmd) {
	if !l.Focused {
		return l, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if l.Selected > 0 {
				l.Selected--
			}
			return l, nil
		case tea.KeyDown:
			if l.Selected < len(l.Items)-1 {
				l.Selected++
			}
			return l, nil
		case tea.KeyHome:
			l.Selected = 0
			return l, nil
		case tea.KeyEnd:
			l.Selected = len(l.Items) - 1
			return l, nil
		case tea.KeyEnter:
			if l.Selected >= 0 && l.Selected < len(l.Items) {
				return l, tea.Cmd(func() tea.Msg {
					return ListItemSelected{
						Item: l.Items[l.Selected],
						Index: l.Selected,
					}
				})
			}
		case tea.KeyTab, tea.KeyShiftTab:
			return l, nil // Let parent handle tab navigation
		}
	}
	
	return l, nil
}

// View returns the list view
func (l *List) View() string {
	if len(l.Items) == 0 {
		return fmt.Sprintf("%s\n\n%s", l.Title, l.Theme.EmptyBox.Render(l.EmptyText))
	}
	
	view := l.Theme.Title.Render(l.Title) + "\n"
	
	start := 0
	end := len(l.Items)
	
	// Calculate visible range
	if l.Height > 0 {
		if len(l.Items) > l.Height {
			start = l.Selected - l.Height/2
			if start < 0 {
				start = 0
			}
			end = start + l.Height
			if end > len(l.Items) {
				end = len(l.Items)
				start = end - l.Height
				if start < 0 {
					start = 0
				}
			}
		}
	}
	
	for i := start; i < end; i++ {
		item := l.Items[i]
		var itemStyle lipgloss.Style
		
		if i == l.Selected {
			itemStyle = l.Theme.Secondary
		} else {
			itemStyle = l.Theme.Muted
		}
		
		if item.Subtitle != "" {
			view += fmt.Sprintf("%s %s\n", itemStyle.Render(item.Title), l.Theme.Muted.Render(item.Subtitle))
		} else {
			view += fmt.Sprintf("%s\n", itemStyle.Render(item.Title))
		}
	}
	
	return view
}

// SetItems sets the list items
func (l *List) SetItems(items []ListItem) {
	l.Items = items
	if l.Selected >= len(items) {
		l.Selected = len(items) - 1
	}
	if l.Selected < 0 && len(items) > 0 {
		l.Selected = 0
	}
}

// GetSelectedItem returns the selected item
func (l *List) GetSelectedItem() (ListItem, bool) {
	if l.Selected >= 0 && l.Selected < len(l.Items) {
		return l.Items[l.Selected], true
	}
	return ListItem{}, false
}

// ListItemSelected represents a list item selection message
type ListItemSelected struct {
	Item  ListItem
	Index int
}

// Container represents a container component
type Container struct {
	BaseComponent
	Children  []Component
	Title     string
	Border    bool
	Focused   bool
}

// NewContainer creates a new container
func NewContainer(id, title string, theme *Theme, keys *Keys) *Container {
	return &Container{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Children:      []Component{},
		Title:        title,
		Border:       true,
		Focused:      false,
	}
}

// Init initializes the container
func (c *Container) Init() tea.Cmd {
	return nil
}

// Update updates the container
func (c *Container) Update(msg tea.Msg) (Component, tea.Cmd) {
	var cmds []tea.Cmd
	
	for i, child := range c.Children {
		updatedChild, cmd := child.Update(msg)
		c.Children[i] = updatedChild
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	
	return c, tea.Batch(cmds...)
}

// View returns the container view
func (c *Container) View() string {
	var content strings.Builder
	
	if c.Title != "" {
		content.WriteString(c.Theme.Title.Render(c.Title) + "\n")
	}
	
	for _, child := range c.Children {
		content.WriteString(child.View() + "\n")
	}
	
	if c.Border {
		return c.Theme.Border.Render(content.String())
	}
	
	return content.String()
}

// AddChild adds a child component
func (c *Container) AddChild(child Component) {
	c.Children = append(c.Children, child)
}

// GetFocusedChild returns the focused child
func (c *Container) GetFocusedChild() (Component, bool) {
	for _, child := range c.Children {
		if child.IsFocused() {
			return child, true
		}
	}
	return nil, false
}

// Text represents a text component
type Text struct {
	BaseComponent
	Content string
	Wrap    bool
}

// NewText creates a new text component
func NewText(id, content string, theme *Theme, keys *Keys) *Text {
	return &Text{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Content:      content,
		Wrap:         false,
	}
}

// Init initializes the text
func (t *Text) Init() tea.Cmd {
	return nil
}

// Update updates the text
func (t *Text) Update(msg tea.Msg) (Component, tea.Cmd) {
	return t, nil
}

// View returns the text view
func (t *Text) View() string {
	if t.Wrap {
		return lipgloss.NewStyle().Width(80).Render(t.Content)
	}
	return t.Content
}

// SetContent sets the text content
func (t *Text) SetContent(content string) {
	t.Content = content
}

// Spinner represents a loading spinner
type Spinner struct {
	BaseComponent
	Style     lipgloss.Style
	Animating bool
}

// NewSpinner creates a new spinner
func NewSpinner(id string, theme *Theme, keys *Keys) *Spinner {
	return &Spinner{
		BaseComponent: *NewBaseComponent(id, theme, keys),
		Style:        theme.Secondary,
		Animating:    false,
	}
}

// Init initializes the spinner
func (s *Spinner) Init() tea.Cmd {
	return s.Tick()
}

// Tick returns a tick command for the spinner
func (s *Spinner) Tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

// Update updates the spinner
func (s *Spinner) Update(msg tea.Msg) (Component, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		if s.Animating {
			return s, s.Tick()
		}
	}
	return s, nil
}

// View returns the spinner view
func (s *Spinner) View() string {
	if !s.Animating {
		return s.Theme.Muted.Render("○")
	}
	
	// Simple spinner animation
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return s.Style.Render(spinners[time.Now().UnixNano()%10])
}

// TickMsg represents a tick message for the spinner
type TickMsg struct{}

// SetAnimating sets the animation state
func (s *Spinner) SetAnimating(animating bool) {
	s.Animating = animating
}

// Focus focuses the spinner
func (s *Spinner) Focus() tea.Cmd {
	s.Focused = true
	return nil
}

// Blur blurs the spinner
func (s *Spinner) Blur() tea.Cmd {
	s.Focused = false
	return nil
}