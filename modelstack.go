package modelstack

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/thisguycodes/modelstack/internal/stack"
)

type PushModel struct {
	Model tea.Model
}

type PopModel struct {
	Msgs []tea.Msg
}

func Push(m tea.Model) func() tea.Msg {
	return func() tea.Msg {
		return PushModel{
			Model: m,
		}
	}
}

func Pop(msgs ...tea.Msg) tea.Msg {
	return PopModel{
		Msgs: msgs,
	}
}

func New(m tea.Model) ModelStack {
	return ModelStack{
		current: m,
		stack:   stack.New[tea.Model](),
	}
}

type ModelStack struct {
	current    tea.Model
	lastResize tea.WindowSizeMsg
	stack      *stack.Stack[tea.Model]
}

func (m ModelStack) View() string {
	return m.current.View()
}

func (m *ModelStack) updateCurrent(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return cmd
}

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

func (m ModelStack) Init() tea.Cmd {
	return m.current.Init()
}
