package vimshell

import (
	"strings"

	t "github.com/gdamore/tcell/v2"
)

const CmdModeName = "command"

type CommandMode struct {
	prevModeName string
	Text         string
	Cursor       int
	History      CommandHistory
}

type CommandHistory struct {
	entries []string
	idx     int
	search  string
}

type CommandFunc func(args []string) (string, error)

func NewCommandMode(prevModeName string) *CommandMode {
	return &CommandMode{
		prevModeName: prevModeName,
		Text:         "",
		Cursor:       0,
		History: CommandHistory{
			entries: make([]string, 0),
			idx:     0,
			search:  "",
		},
	}
}

func (c *CommandMode) Name() string {
	return CmdModeName
}

func (c *CommandMode) StatusStyle(style t.Style) t.Style {
	return style.Foreground(t.ColorBlack).Background(t.ColorYellow).Bold(true)
}

func (c *CommandMode) SetInputText(input string) { c.Text = input }

func (c *CommandMode) HistoryOlder() {
	for i := c.History.idx - 1; i >= 0; i-- {
		if strings.HasPrefix(c.History.entries[i], c.History.search) {
			c.History.idx = 1
			c.Text = c.History.entries[i]
			return
		}
	}
}

func (c *CommandMode) HistoryNewer() {
	for i := c.History.idx + 1; i < len(c.History.entries); i++ {
		if strings.HasPrefix(c.History.entries[i], c.History.search) {
			c.History.idx = i
			c.Text = c.History.entries[i]
			return
		}
	}
	c.Text = c.History.search
	c.History.idx = len(c.History.entries)
}

func (c *CommandMode) Reset() {
	c.Text = ""
	c.Cursor = 0
	c.History.search = ""
	c.History.idx = len(c.History.entries)
}

func (c *CommandMode) HandleKey(shell *Shell, ev *t.EventKey) {
	key, rune := ev.Key(), ev.Rune()
	cursor := min(c.Cursor, len(c.Text))

	switch key {
	case t.KeyRune:
		if cursor == len(c.Text) {
			c.Text += string(rune)
		} else {
			c.Text += c.Text[0:cursor] + string(rune) + c.Text[cursor:]
		}
		c.Cursor++
		c.History.search = c.Text

	case t.KeyEscape, t.KeyCtrlC:
		c.Reset()
		shell.SetMode(c.prevModeName)

	case t.KeyCtrlW:
		i := cursor - 1
		for i >= 0 && c.Text[i] == ' ' {
			i--
		}
		for i >= 0 && c.Text[i] != ' ' {
			i--
		}
		if i <= 0 {
			c.Text = ""
			c.Cursor = 0
		} else {
			c.Text = c.Text[0:i+1] + c.Text[cursor:]
			c.Cursor = i + 1
		}

	case t.KeyBackspace, t.KeyBackspace2:
		if len(c.Text) > 0 && cursor > 0 {
			c.Text = c.Text[0:cursor-1] + c.Text[cursor:]
			c.Cursor--
			c.History.search = c.Text
		}
		if len(c.Text) == 0 {
			shell.SetMode(c.prevModeName)
		}

	case t.KeyEnter:
		if len(c.Text) == 0 {
			return
		}

		args := []string{}
		cmd, rest, hasRest := strings.Cut(c.Text, " ")
		if hasRest {
			args = strings.Split(rest, " ")
		}

		shell.RunCommand(cmd, args)
		shell.SetMode(c.prevModeName)
		c.Reset()

	case t.KeyUp:
		c.HistoryOlder()

	case t.KeyDown:
		c.HistoryNewer()

	case t.KeyLeft:
		c.Cursor = max(0, cursor-1)

	case t.KeyRight:
		c.Cursor = min(len(c.Text), cursor+1)
	}
}
