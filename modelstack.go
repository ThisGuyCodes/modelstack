package modelstack

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thisguycodes/modelstack/internal/stack"
)

// PushModel is a tea.Msg that instructs a ModelStack to add a new tea.Model to
// the stack and pass control to it
type PushModel struct {
	Model tea.Model
}

// PopModel is a tea.Msg that instructs a ModelStack to remove the current
// tea.Model from the stack and pass control to the previous tea.Model
type PopModel struct {
	Msgs []tea.Msg
}

// Push is a tea.Cmd used to send a PushModel tea.Msg
func Push(m tea.Model) func() tea.Msg {
	return func() tea.Msg {
		return PushModel{
			Model: m,
		}
	}
}

// Pop is a tea.Cmd used to send a PopModel tea.Msg
func Pop(msgs ...tea.Msg) tea.Msg {
	return PopModel{
		Msgs: msgs,
	}
}

// New creates a new Modelstack with an initial tea.Model
func New(m tea.Model) ModelStack {
	return ModelStack{
		current: m,
		stack:   stack.New[tea.Model](),
	}
}

// ModelStack is a tea.Model that listens for PushModel and PopModel tea.Msg's
// to switch control between several tea.Model's
type ModelStack struct {
	current    tea.Model
	lastResize tea.WindowSizeMsg
	stack      *stack.Stack[tea.Model]
}

// View fulfills the tea.Model interface, rendering the currently active
// tea.Model on its stack
func (m ModelStack) View() string {
	return m.current.View()
}

func (m *ModelStack) updateCurrent(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return cmd
}

// Update fulfills the tea.Model interface, intercepting PushModel and PopModel
// tea.Msg's, as well as information needed to initialize models as they switch
// control (e.g. the most recent tea.WindowSizeMsg)
func (m ModelStack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastResize = msg
	case PushModel:
		m.stack.Push(m.current)
		m.current = msg.Model
		cmd := m.current.Init()
		cmd2 := m.updateCurrent(m.lastResize)
		return m, tea.Batch(cmd, cmd2)
	case PopModel:
		m.current = m.stack.Pop().Value
		cmds := make(tea.BatchMsg, 2+len(msg.Msgs))
		cmds[0] = m.current.Init()
		cmds[1] = m.updateCurrent(m.lastResize)
		for i, msg := range msg.Msgs {
			cmds[i+2] = m.updateCurrent(msg)
		}
		return m, tea.Batch(cmds...)
	}

	cmd := m.updateCurrent(msg)

	return m, cmd
}

// Init fulfills the tea.Model interface, currently it just calls Init() on the
// current tea.Model on the stack
func (m ModelStack) Init() tea.Cmd {
	return m.current.Init()
}
