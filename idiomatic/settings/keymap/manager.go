package keymap

func NewManager() *Manager {
	return &Manager{}
}

type Manager struct {
}

func (this Manager) KeyboardEvent(event KeyboardEvent) []ApplicationCommand {
	var commands []ApplicationCommand

	return commands
}
