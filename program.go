package vimshell

import (
	"github.com/gdamore/tcell/v2"
)

func (shell *Shell) Run(scrn tcell.Screen, draw func()) {
	for {
		scrn.Clear()
		draw()
		scrn.Show()

		ev := scrn.PollEvent()
		if ev == nil {
			return
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			scrn.Sync()
		case *tcell.EventKey:
			shell.HandleKey(ev)
		}
	}
}
