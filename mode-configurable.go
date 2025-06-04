package vimshell

import t "github.com/gdamore/tcell/v2"

type ConfigurableMode struct {
	name   string
	Keymap map[t.Key]KeyHandlerFunc
}

func NewMode(name string) *ConfigurableMode {
	return &ConfigurableMode{
		name:   name,
		Keymap: map[t.Key]KeyHandlerFunc{},
	}
}

func (mode *ConfigurableMode) Name() string {
	return mode.name
}

func (mode *ConfigurableMode) StatusStyle(style t.Style) t.Style {
	return style.Background(t.ColorGreen).Foreground(t.ColorBlack).Bold(true)
}

func (n *ConfigurableMode) MapKey(key t.Key, fn KeyHandlerFunc) {
	n.Keymap[key] = fn
}

func (n *ConfigurableMode) HandleKey(shell *Shell, ev *t.EventKey) {
	key := ev.Key()
	fn, has := n.Keymap[key]
	if has {
		fn(shell, ev)
	}
}
