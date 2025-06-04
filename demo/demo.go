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
	shell.Message = "Hello! Type :q to exit :)"
	shell.CommandLine.AddCommand("q", func(args []string) (string, error) {
		quit()
		return "", nil
	})

	shell.Run(scrn, func() {
		_, height := scrn.Size()
		shell.StatusBar.Render(shell, scrn, height-2)
		shell.CommandLine.Render(shell, scrn, height-1)

		if shell.Mode != "command" {
			scrn.HideCursor()
		}
	})
}
