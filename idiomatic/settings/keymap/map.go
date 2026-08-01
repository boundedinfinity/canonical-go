package keymap

type Key string

type KeyMap interface {
	Run(Key) []Action
}

type Action struct {
	Name string
}
