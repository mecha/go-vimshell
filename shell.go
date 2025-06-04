package vimshell

import (
	"errors"

	t "github.com/gdamore/tcell/v2"
)

type Shell struct {
	Mode        string
	Modes       map[string]Mode
	CmdMode     *CommandMode
	StatusBar   *StatusBar
	CommandLine *CommandLine
	Message     string
	Error       error
}

type Mode interface {
	Name() string
	StatusStyle(style t.Style) t.Style
	HandleKey(shell *Shell, ev *t.EventKey)
}

type KeyHandlerFunc func(shell *Shell, ev *t.EventKey)

// Creates a new shell with the given default mode, and a command mode with
// name "command". For a typical setup, see NewDefaultShell().
func NewShell(defMode Mode) *Shell {
	defModeName := defMode.Name()
	shell := &Shell{
		Mode:        defModeName,
		Modes:       map[string]Mode{defModeName: defMode},
		CmdMode:     NewCommandMode(defModeName),
		CommandLine: NewCmdLine(),
		StatusBar:   NewStatusBar(),
	}
	shell.AddMode(shell.CmdMode)
	return shell
}

// Creates a new shell with vim-like defaults. A configurable Normal mode is
// used as the default mode, which the Command mode will return to.
func NewDefaultShell(quitFn func()) *Shell {
	nrmMode := NewMode("normal")
	nrmMode.Keymap[t.KeyRune] = func(shell *Shell, ev *t.EventKey) {
		if ev.Rune() == ':' {
			shell.Mode = "command"
		}
	}
	shell := NewShell(nrmMode)
	shell.StatusBar = NewStatusBar()
	shell.StatusBar.AddLeftSection(NewModeSection(map[string]StyleFunc{}))
	return shell
}

// The shell's current mode.
func (s *Shell) CurrMode() Mode { return s.Modes[s.Mode] }

// Changes the shell's current mode.
func (s *Shell) SetMode(name string) error {
	if _, has := s.Modes[name]; !has {
		return errors.New("invalid mode")
	}
	s.Mode = name
	return nil
}

// Adds a new mode to the shell.
// Only one mode may be registered for a given name. The name of the mode is
// determined by the mode's Name() method. If a mode is added using a name that
// is already in use, the existing mode with that name will be replaced.
func (s *Shell) AddMode(mode Mode) {
	s.Modes[mode.Name()] = mode
}

// Handles a key input event using the shell's current mode.
func (s *Shell) HandleKey(ev *t.EventKey) {
	s.CurrMode().HandleKey(s, ev)
}
