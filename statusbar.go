package vimshell

import (
	"strings"

	t "github.com/gdamore/tcell/v2"
)

type StatusBar struct {
	Left         []StatusBarSection
	Right        []StatusBarSection
	ErrStyleFunc StyleFunc
	CmdPrefix    string
	Message      string
	Error        error
}

type StatusBarSection func(s *Shell, scrn t.Screen, style t.Style) (string, t.Style)

// Creates a new empty status bar.
func NewStatusBar() *StatusBar {
	return &StatusBar{
		Left:         make([]StatusBarSection, 0),
		Right:        make([]StatusBarSection, 0),
		ErrStyleFunc: InheritStyle,
		CmdPrefix:    "",
	}
}

// Creates a default vim-like status bar, with colon ":" as the command prefix,
// a bold ANSI bright red style for error mesages, and a single left section
// that shows the current mode.
func NewDefaultStatusBar() *StatusBar {
	sb := NewStatusBar()
	sb.SetCommandPrefix(":")
	sb.SetErrorStyle(func(style t.Style) t.Style {
		return style.Bold(true).Foreground(t.ColorRed)
	})
	sb.AddLeftSection(func(s *Shell, scrn t.Screen, style t.Style) (string, t.Style) {
		modeStatus := " " + strings.ToUpper(s.Mode) + " "
		return modeStatus, s.CurrMode().StatusStyle(style)
	})
	return sb
}

// Adds left sections to the status bar.
func (sb *StatusBar) AddLeftSection(sections ...StatusBarSection) {
	sb.Left = append(sb.Left, sections...)
}

// Adds right sections to the status bar.
func (sb *StatusBar) AddRightSection(sections ...StatusBarSection) {
	sb.Right = append(sb.Right, sections...)
}

// Sets the prefix for command input text. This should typically match the rune
// of the key that switches the shell to command mode.
func (sb *StatusBar) SetCommandPrefix(prefix string) {
	sb.CmdPrefix = prefix
}

// Sets the prefix for command input text. This should typically match the rune
// of the key that switches the shell to command mode.
func (sb *StatusBar) SetErrorStyle(styleFunc StyleFunc) {
	sb.ErrStyleFunc = styleFunc
}

// Renders the status bar on a screen with a given style.
func (sb *StatusBar) Render(shell *Shell, scrn t.Screen, style t.Style) {
	width, height := scrn.Size()

	xl := 0
	for _, section := range sb.Left {
		sectionText, sectionStyle := section(shell, scrn, style)
		writeString(scrn, xl, height-2, sectionText, nil, sectionStyle)
		xl += len(sectionText)
	}

	xr := width
	for _, section := range sb.Right {
		sectionText, sectionStyle := section(shell, scrn, style)
		writeString(scrn, xr, height-2, sectionText, nil, sectionStyle)
		xr -= len(sectionText)
	}

	for x := xl; x < xr; x++ {
		scrn.SetContent(x, height-2, ' ', nil, style)
		scrn.SetContent(x, height-1, ' ', nil, style)
	}

	switch {
	case shell.Mode == "command":
		text := sb.CmdPrefix + shell.CmdMode.Text
		writeString(scrn, 0, height-1, text, nil, style)

		cursorX := shell.CmdMode.Cursor + len(sb.CmdPrefix)
		scrn.ShowCursor(cursorX, height-1)

	case sb.Error != nil:
		errStyle := sb.ErrStyleFunc(style)
		writeString(scrn, 0, height-1, sb.Error.Error(), nil, errStyle)

	default:
		writeString(scrn, 0, height-1, sb.Message, nil, style)
	}
}
