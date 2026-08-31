package keymap

type Map interface {
	OnKey(Key) (Map, []ApplicationCommand)
	OnEvent(ApplicationEvent) (Map, []ApplicationCommand)
	Available() []Map
}
