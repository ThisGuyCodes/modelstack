/*
Package modelstack impliments a [tea.Model] (from [Bubbletea]) that listens for
[PushModel] and [PopModel] [tea.Msg]-s to switch control between several
[tea.Model]-s, including some initialization maintinance.

This allows easier re-use of [tea.Model]-s 'within' each other. However it
breaks some assumptions about how your model will be used by bubbletea. The
noteable behavioral differences are listed here:

  - .Init() is called on your [tea.Model] every time it (re)gains control (not
    just once).
  - While another [tea.Model] is in control, your [tea.Model] will not receive
    any [tea.Msg]-s.
  - When a [tea.Cmd] resolves, its [tea.Msg] will be sent to the currently
    active [tea.Model], regardless of which [tea.Model] sent the [tea.Cmd].

For the most part, things should "just work", but keep a few things in mind
when designing [tea.Model]-s you want to be compatible with this package:

  - Be sure to initialize ticks in .Init()
  - Wait for [tea.Cmd]-s you want to resolve before passing control.
    [tea.Sequence] can be useful for this.
  - Be able to handle [tea.Msg]-s going missing if it's possible for control to
    be lost mid-processing of a [tea.Cmd].
  - Remember the model on the stack is the same model, not a re-initialized one
    (though .Init() is called again).

[Bubbletea]: https://github.com/charmbracelet/bubbletea
*/
package modelstack
