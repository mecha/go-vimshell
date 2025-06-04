package vimshell

import t "github.com/gdamore/tcell/v2"

type StyleFunc func(style t.Style) t.Style

var InheritStyle = func(style t.Style) t.Style {
	return style
}
