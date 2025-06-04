package vimshell

import (
	"errors"

	t "github.com/gdamore/tcell/v2"
)

type Shell struct {
	Mode      string
	Modes     map[string]Mode
	CmdMode   *CommandMode
	StatusBar *StatusBar
	Commands  map[string]CommandFunc
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
		Mode:      defModeName,
		Modes:     map[string]Mode{defModeName: defMode},
		CmdMode:   NewCommandMode(defModeName),
		Commands:  map[string]CommandFunc{},
		StatusBar: NewDefaultStatusBar(),
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
	shell.StatusBar = NewDefaultStatusBar()
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

// Adds a command to the shell.
// The name of the command must be unique. If a command already exists with the
// given name, the existing command will be replaced.
func (s *Shell) AddCommand(name string, fn CommandFunc) {
	s.Commands[name] = fn
}

// Runs a command with the given arguments.
// The outcome is stored in the status bar in the form of output text or an error.
func (s *Shell) RunCommand(command string, args []string) {
	if cmd, has := s.Commands[command]; has {
		s.StatusBar.Message, s.StatusBar.Error = cmd(args)
	} else {
		s.StatusBar.Message, s.StatusBar.Error = "", errors.New("unknown command: "+command)
	}
}

// Handles a key input event using the shell's current mode.
func (s *Shell) HandleKey(ev *t.EventKey) {
	s.CurrMode().HandleKey(s, ev)
}

// Renders the status bar
func (s *Shell) RenderStatusBar(scrn t.Screen, style t.Style) {
	s.StatusBar.Render(s, scrn, style)
}
