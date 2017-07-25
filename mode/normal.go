package mode

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/linde12/kod/commands"
	"github.com/linde12/kod/editor"
)

type NormalMode struct {
	editor *editor.Editor
	count  string
}

func NewNormalMode(e *editor.Editor) *NormalMode {
	return &NormalMode{e, ""}
}

func (m *NormalMode) OnKey(ev *tcell.EventKey) {
	editor := m.editor
	r := ev.Rune()
	if (r > '0' && r < '9') || (r == '0' && len(m.count) > 0) {
		m.count += string(r)
		return
	}

	// TODO: Handle error
	nrepeat, _ := strconv.Atoi(m.count)
	// Execute command at least once
	if nrepeat == 0 {
		nrepeat = 1
	}

	if ev.Key() == tcell.KeyRune {
		switch r {
		case 'A':
			editor.Commands <- commands.MoveEOL{}
			editor.SetMode(NewInsertMode(editor))
		case 'h':
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveLeft}, nrepeat}
		case 'k':
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveUp}, nrepeat}
		case 'l':
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveRight}, nrepeat}
		case 'j':
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveDown}, nrepeat}
		case 'i':
			editor.SetMode(NewInsertMode(editor))
		}
	} else {
		switch ev.Key() {
		case tcell.KeyLeft:
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveLeft}, nrepeat}
		case tcell.KeyUp:
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveUp}, nrepeat}
		case tcell.KeyRight:
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveRight}, nrepeat}
		case tcell.KeyDown:
			editor.Commands <- commands.Repeat{commands.MoveRune{Dir: commands.MoveDown}, nrepeat}
		}
	}

	m.count = ""
}