package modelstack

import (
	"context"
	"log/slog"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/thisguycodes/modelstack/internal/stack"
)

// Option is used to set options when initializing a ModelStack
type Option func(*ModelStack)

// WithSlogger attaches a *slog.Logger to a ModelStack
func WithSlogger(l *slog.Logger) Option {
	return func(ms *ModelStack) {
		ms.l = l
	}
}

// PushModel is a tea.Msg that instructs a ModelStack to add a new tea.Model to
// the stack and pass control to it
type PushModel struct {
	Model tea.Model
}

// PopModel is a tea.Msg that instructs a ModelStack to remove the current
// tea.Model from the stack and pass control to the previous tea.Model
type PopModel struct {
	// Msgs are tea.Msg-s that will be passed to the parent model synchronously
	// when it regains control (before its first call to .View(), after calling
	// .Init())
	Msgs []tea.Msg
}

// Push returns s a tea.Cmd used to send a PushModel tea.Msg
func Push(m tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModel{
			Model: m,
		}
	}
}

// Pop returns a tea.Cmd used to send a PopModel tea.Msg. Optionally you can
// give it tea.Msg-s to pass to the parent tea.Model.
func Pop(msgs ...tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return PopModel{
			Msgs: msgs,
		}
	}
}

// New creates a new Modelstack with an initial tea.Model
func New(m tea.Model, opts ...Option) ModelStack {
	ms := ModelStack{
		current: m,
		stack:   stack.New[tea.Model](),
	}
	for _, opt := range opts {
		opt(&ms)
	}
	return ms
}

// ModelStack is a tea.Model that listens for PushModel and PopModel tea.Msg's
// to switch control between several tea.Model's
type ModelStack struct {
	current    tea.Model
	lastResize tea.WindowSizeMsg
	stack      *stack.Stack[tea.Model]
	l          *slog.Logger
}

// View fulfills the tea.Model interface, rendering the currently active
// tea.Model on its stack
func (m ModelStack) View() string {
	m.l.LogAttrs(context.TODO(), slog.LevelDebug, ".View() called")

	return m.current.View()
}

func (m *ModelStack) updateCurrent(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return cmd
}

func lazyLogType[T any](val T) slogType[T] {
	return slogType[T]{
		val: val,
	}
}

type slogType[T any] struct {
	val T
}

func (st slogType[T]) LogValue() slog.Value {
	t := reflect.TypeOf(st.val)
	return slog.StringValue(t.String())
}

// Update fulfills the tea.Model interface, intercepting PushModel and PopModel
// tea.Msg's, as well as information needed to initialize models as they switch
// control (e.g. the most recent tea.WindowSizeMsg)
func (m ModelStack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.l.LogAttrs(context.TODO(), slog.LevelDebug, ".Update() called",
		slog.Any("msg", msg),
		slog.Any("msg.(type)", lazyLogType(msg)),
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastResize = msg
	case PushModel:
		slog.Debug("pushing to stack")
		m.stack.Push(m.current)
		m.current = msg.Model
		cmd := m.current.Init()
		cmd2 := m.updateCurrent(m.lastResize)
		return m, tea.Batch(cmd, cmd2)
	case PopModel:
		slog.Debug("popping off stack")
		m.current = m.stack.Pop().Value
		cmds := make(tea.BatchMsg, 2+len(msg.Msgs))
		cmds[0] = m.current.Init()
		cmds[1] = m.updateCurrent(m.lastResize)
		for i, msg := range msg.Msgs {
			cmds[i+2] = m.updateCurrent(msg)
		}
		return m, tea.Batch(cmds...)
	}

	slog.Debug("passing through message")
	cmd := m.updateCurrent(msg)

	return m, cmd
}

// Init fulfills the tea.Model interface, currently it just calls Init() on the
// current tea.Model on the stack
func (m ModelStack) Init() tea.Cmd {
	slog.Debug(".Init() called")
	return m.current.Init()
}
