package main

import (
	"log"
	"os"

	t "github.com/gdamore/tcell/v2"
	v "github.com/mecha/vimshell"
)

func main() {
	scrn, err := t.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err := scrn.Init(); err != nil {
		log.Fatal(err)
	}
	defer scrn.Fini()

	quit := func() {
		scrn.Fini()
		os.Exit(0)
	}

	shell := v.NewDefaultShell(quit)
	shell.StatusBar.Message = "Hello! Type :q to exit :)"
	shell.AddCommand("q", func(args []string) (string, error) {
		quit()
		return "", nil
	})

	shell.Run(scrn, func() {
		shell.RenderStatusBar(scrn, t.StyleDefault.Background(t.ColorBlack))
		if shell.Mode != "command" {
			scrn.HideCursor()
		}
	})
}
