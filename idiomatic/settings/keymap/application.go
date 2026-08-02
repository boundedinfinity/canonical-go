package keymap

type ApplicationMessage struct {
	Name        string
	Application string
	Message     string
	Description string
}

type ApplicationEvent struct {
	Name        string
	Application string
	Event       string
	Description string
}
