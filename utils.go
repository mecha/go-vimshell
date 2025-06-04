package vimshell

import (
	t "github.com/gdamore/tcell/v2"
)

func writeString(scr t.Screen, x, y int, str string, combc []rune, style t.Style) {
	for i, r := range str {
		scr.SetContent(x+i, y, r, combc, style)
	}
}

func repeat(scr t.Screen, x, y int, rune rune, combc []rune, style t.Style) {
	width, _ := scr.Size()
	for eks := x; eks < width; eks++ {
		scr.SetContent(eks, y, rune, combc, style)
	}
}
