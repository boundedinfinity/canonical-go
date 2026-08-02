package keymap

type Map interface {
	Run(Key) []ApplicationMessage
}
