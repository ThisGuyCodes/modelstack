package modelstack_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/thisguycodes/modelstack"
)

type dumbModel struct {
	view      string
	msgs      *[]tea.Msg
	cmds      *[]tea.Cmd
	initcalls *int
}

var _ tea.Model = dumbModel{}

func newDumbModel(view string) dumbModel {
	return dumbModel{
		view:      view,
		initcalls: new(int),
		msgs:      new([]tea.Msg),
		cmds:      new([]tea.Cmd),
	}
}
func (m dumbModel) Init() tea.Cmd {
	*m.initcalls++
	return nil
}
func (m dumbModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	*m.msgs = append(*m.msgs, msg)

	var cmd tea.Cmd
	if len(*m.cmds) > 0 {
		cmd = (*m.cmds)[0]
		*m.cmds = (*m.cmds)[1:]
	}
	return m, cmd
}
func (m dumbModel) View() string { return m.view }

func TestModelstack_Init(t *testing.T) {
	t.Parallel()

	ms := modelstack.New(newDumbModel("butts"))

	assert.Equal(t, "butts", ms.View())
}

func TestModelstack_Push(t *testing.T) {
	t.Parallel()

	var ms tea.Model = modelstack.New(newDumbModel("butts"))
	ms, _ = ms.Update(modelstack.PushModel{Model: newDumbModel("alsobutts")})

	assert.Equal(t, "alsobutts", ms.View())
}

func TestModelstack_Pop(t *testing.T) {
	t.Parallel()

	var ms tea.Model = modelstack.New(newDumbModel("butts"))
	ms, _ = ms.Update(modelstack.PushModel{Model: newDumbModel("alsobutts")})
	ms, _ = ms.Update(modelstack.PopModel{})

	assert.Equal(t, "butts", ms.View())
}

func TestModelstack_InitPassthrough(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")
	var ms tea.Model = modelstack.New(dm)

	assert.Equal(t, 0, *dm.initcalls)
	ms.Init()
	assert.Equal(t, 1, *dm.initcalls)
	ms.Init()
	assert.Equal(t, 2, *dm.initcalls)
}

func TestModelstack_InitOnPush(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")

	var ms tea.Model = modelstack.New(dm)

	dm2 := newDumbModel("alsobutts")

	ms, _ = ms.Update(modelstack.PushModel{Model: dm2})

	assert.Equal(t, 1, *dm2.initcalls)
}

func TestModelstack_InitPassthroughAfterPush(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")

	var ms tea.Model = modelstack.New(dm)

	dm2 := newDumbModel("alsobutts")

	ms, _ = ms.Update(modelstack.PushModel{Model: dm2})

	ms.Init()

	assert.Equal(t, 2, *dm2.initcalls)
}

func TestModelstack_ReInitOnPop(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")

	var ms tea.Model = modelstack.New(dm)

	dm2 := newDumbModel("alsobutts")

	ms, _ = ms.Update(modelstack.PushModel{Model: dm2})
	ms, _ = ms.Update(modelstack.PopModel{})

	assert.Equal(t, 1, *dm.initcalls)
}

func TestModelstack_InitPassthroughAfterPop(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")

	var ms tea.Model = modelstack.New(dm)

	dm2 := newDumbModel("alsobutts")

	ms, _ = ms.Update(modelstack.PushModel{Model: dm2})
	ms, _ = ms.Update(modelstack.PopModel{})
	ms.Init()

	assert.Equal(t, 2, *dm.initcalls)
}

func TestModelstack_PassthroughMsgs(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")

	var ms tea.Model = modelstack.New(dm)

	ms, _ = ms.Update("butts")

	assert.Equal(t, "butts", (*dm.msgs)[0])
}

func TestModelstack_PassBackCmds(t *testing.T) {
	t.Parallel()
	dm := newDumbModel("butts")
	*dm.cmds = append(*dm.cmds, func() tea.Msg { return "alsobutts" })

	var ms tea.Model = modelstack.New(dm)

	var cmd tea.Cmd
	ms, cmd = ms.Update("butts")

	assert.Equal(t, "alsobutts", cmd())
}

func TestModelstack_CacheResize(t *testing.T) {
	t.Parallel()
	dm := newDumbModel("butts")
	var ms tea.Model = modelstack.New(dm)

	size := tea.WindowSizeMsg{Height: 10, Width: 42}
	ms, _ = ms.Update(size)
	dm2 := newDumbModel("alsobutts")

	ms, _ = ms.Update(modelstack.PushModel{dm2})

	assert.Equal(t, size, (*dm2.msgs)[0])
}

func TestModelstack_PopMsgs(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")
	var ms tea.Model = modelstack.New(dm)

	dm2 := newDumbModel("alsobutts")
	ms, _ = ms.Update(modelstack.PushModel{Model: dm2})

	ms, _ = ms.Update(modelstack.PopModel{Msgs: []tea.Msg{"morebutts"}})

	// 1 because the first msg we get is a zero tea.WindowSizeMsg
	assert.Equal(t, "morebutts", (*dm.msgs)[1])
}

func TestPush(t *testing.T) {
	t.Parallel()

	dm := newDumbModel("butts")
	msg := modelstack.Push(dm)()

	assert.Equal(t, dm, msg.(modelstack.PushModel).Model)
}

func TestPop(t *testing.T) {
	t.Parallel()

	msg := modelstack.Pop("butts")

	assert.Equal(t, "butts", msg.(modelstack.PopModel).Msgs[0])
}
