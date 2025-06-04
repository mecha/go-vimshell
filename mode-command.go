package vimshell

import (
	t "github.com/gdamore/tcell/v2"
)

const CmdModeName = "command"

type CommandMode struct {
	ReturnMode string
}

func NewCommandMode(prevModeName string) *CommandMode {
	return &CommandMode{
		ReturnMode: prevModeName,
	}
}

func (c *CommandMode) Name() string {
	return CmdModeName
}

func (c *CommandMode) StatusStyle(style t.Style) t.Style {
	return style.Foreground(t.ColorBlack).Background(t.ColorYellow).Bold(true)
}

func (mode *CommandMode) HandleKey(shell *Shell, ev *t.EventKey) {
	key, rune := ev.Key(), ev.Rune()

	switch key {
	case t.KeyRune:
		shell.CommandLine.Feed(rune)

	case t.KeyBackspace, t.KeyBackspace2:
		shell.CommandLine.Backspace()
		if len(shell.CommandLine.Text) == 0 {
			shell.SetMode(mode.ReturnMode)
		}

	case t.KeyCtrlW:
		shell.CommandLine.DeleteWord()

	case t.KeyEnter:
		shell.CommandLine.Submit(shell)
		fallthrough
	case t.KeyEscape, t.KeyCtrlC:
		shell.CommandLine.Reset()
		shell.SetMode(mode.ReturnMode)

	case t.KeyUp:
		shell.CommandLine.HistoryOlder()
	case t.KeyDown:
		shell.CommandLine.HistoryNewer()
	case t.KeyLeft:
		shell.CommandLine.CursorLeft()
	case t.KeyRight:
		shell.CommandLine.CursorRight()
	}
}
