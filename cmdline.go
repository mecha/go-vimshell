package vimshell

import (
	"errors"
	"strings"

	t "github.com/gdamore/tcell/v2"
)

type CommandLine struct {
	Commands      map[string]CommandFunc
	Prefix        string
	Text          string
	Cursor        int
	History       []string
	BaseStyle     t.Style
	CmdStyle      t.Style
	MsgStyle      t.Style
	ErrStyle      t.Style
	historyIdx    int
	historySearch string
}

type CommandFunc func(args []string) (string, error)

func NewCmdLine() *CommandLine {
	return &CommandLine{
		Commands:      map[string]CommandFunc{},
		Text:          "",
		Cursor:        0,
		History:       []string{},
		BaseStyle:     t.StyleDefault,
		CmdStyle:      t.StyleDefault,
		MsgStyle:      t.StyleDefault,
		ErrStyle:      t.StyleDefault.Foreground(t.ColorRed).Bold(true),
		historyIdx:    0,
		historySearch: "",
	}
}

func (cl *CommandLine) AddCommand(name string, fn CommandFunc) {
	cl.Commands[name] = fn
}

func (cl *CommandLine) GetCommand(name string, args []string) func(shell *Shell) {
	return func(shell *Shell) {
		if cmd, has := cl.Commands[name]; has {
			shell.Message, shell.Error = cmd(args)
		} else {
			shell.Message, shell.Error = "", errors.New("unknown command: "+name)
		}
	}
}

func (cl *CommandLine) CursorLeft()  { cl.Cursor = max(cl.Cursor-1, 0) }
func (cl *CommandLine) CursorRight() { cl.Cursor = min(cl.Cursor+1, len(cl.Text)) }

func (cl *CommandLine) Feed(rune rune) {
	cursor := min(cl.Cursor, len(cl.Text))
	if cursor == len(cl.Text) {
		cl.Text += string(rune)
	} else {
		cl.Text += cl.Text[0:cursor] + string(rune) + cl.Text[cursor:]
	}
	cl.Cursor++
	cl.historySearch = cl.Text
}

func (cl *CommandLine) Backspace() {
	cursor := min(cl.Cursor, len(cl.Text))

	if len(cl.Text) > 0 && cursor > 0 {
		cl.Text = cl.Text[0:cursor-1] + cl.Text[cursor:]
		cl.Cursor--
		cl.historySearch = cl.Text
	}
}

func (cl *CommandLine) DeleteWord() {
	cursor := min(cl.Cursor, len(cl.Text))

	i := cursor - 1
	for i >= 0 && cl.Text[i] == ' ' {
		i--
	}
	for i >= 0 && cl.Text[i] != ' ' {
		i--
	}
	if i <= 0 {
		cl.Text = ""
		cl.Cursor = 0
	} else {
		cl.Text = cl.Text[0:i+1] + cl.Text[cursor:]
		cl.Cursor = i + 1
	}
}

func (cl *CommandLine) Parse() (string, []string) {
	if len(cl.Text) == 0 {
		return "", nil
	}

	args := []string{}
	cmd, rest, hasRest := strings.Cut(cl.Text, " ")
	if hasRest {
		args = strings.Split(rest, " ")
	}

	return cmd, args
}

func (cl *CommandLine) Submit(shell *Shell) {
	cmd, args := cl.Parse()
	if cmd != "" {
		cl.GetCommand(cmd, args)(shell)
	}
	cl.Reset()
}

func (cl *CommandLine) Reset() {
	cl.Text = ""
	cl.Cursor = 0
	cl.historySearch = ""
	cl.historyIdx = len(cl.History)
}

func (cl *CommandLine) HistoryOlder() {
	for i := cl.historyIdx - 1; i >= 0; i-- {
		if strings.HasPrefix(cl.History[i], cl.historySearch) {
			cl.historyIdx = 1
			cl.Text = cl.History[i]
			return
		}
	}
}

func (cl *CommandLine) HistoryNewer() {
	for i := cl.historyIdx + 1; i < len(cl.History); i++ {
		if strings.HasPrefix(cl.History[i], cl.historySearch) {
			cl.historyIdx = i
			cl.Text = cl.History[i]
			return
		}
	}
	cl.Text = cl.historySearch
	cl.historyIdx = len(cl.History)
}

func (cl *CommandLine) Render(shell *Shell, scrn t.Screen, y int) {
	switch {
	case shell.Mode == "command":
		cmdText := cl.Prefix + cl.Text
		writeString(scrn, 0, y, cmdText, nil, cl.CmdStyle)
		repeat(scrn, len(cmdText), y, ' ', nil, cl.BaseStyle)

		cursorX := cl.Cursor + len(cl.Prefix)
		scrn.ShowCursor(cursorX, y)

	case shell.Error != nil:
		errMessage := shell.Error.Error()
		writeString(scrn, 0, y, errMessage, nil, cl.ErrStyle)
		repeat(scrn, len(errMessage), y, ' ', nil, cl.BaseStyle)

	default:
		writeString(scrn, 0, y, shell.Message, nil, cl.MsgStyle)
		repeat(scrn, len(shell.Message), y, ' ', nil, cl.BaseStyle)
	}
}
