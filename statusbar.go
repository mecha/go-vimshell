package vimshell

import (
	"strings"

	t "github.com/gdamore/tcell/v2"
)

type StatusBar struct {
	Left  []StatusBarSection
	Right []StatusBarSection
	Style t.Style
}

type StatusBarSection func(s *Shell, scrn t.Screen, style t.Style) (string, t.Style)

// Creates a new empty status bar.
func NewStatusBar() *StatusBar {
	return &StatusBar{
		Left:  make([]StatusBarSection, 0),
		Right: make([]StatusBarSection, 0),
	}
}

// Adds left sections to the status bar.
func (sb *StatusBar) AddLeftSection(sections ...StatusBarSection) {
	sb.Left = append(sb.Left, sections...)
}

// Adds right sections to the status bar.
func (sb *StatusBar) AddRightSection(sections ...StatusBarSection) {
	sb.Right = append(sb.Right, sections...)
}

// Renders the status bar on a screen with a given style.
func (sb *StatusBar) Render(shell *Shell, scrn t.Screen, y int) {
	width, _ := scrn.Size()

	xl := 0
	for _, section := range sb.Left {
		sectionText, sectionStyle := section(shell, scrn, sb.Style)
		writeString(scrn, xl, y, sectionText, nil, sectionStyle)
		xl += len(sectionText)
	}

	xr := width
	for _, section := range sb.Right {
		sectionText, sectionStyle := section(shell, scrn, sb.Style)
		writeString(scrn, xr, y, sectionText, nil, sectionStyle)
		xr -= len(sectionText)
	}

	for x := xl; x < xr; x++ {
		scrn.SetContent(x, y, ' ', nil, sb.Style)
	}
}

func NewModeSection(styles map[string]StyleFunc) StatusBarSection {
	return func(s *Shell, scrn t.Screen, style t.Style) (string, t.Style) {
		modeStatus := strings.ToUpper(s.Mode)
		if styleFunc, has := styles[s.Mode]; has {
			style = styleFunc(style)
		}
		return modeStatus, style
	}
}
