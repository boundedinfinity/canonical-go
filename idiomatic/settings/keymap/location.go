package keymap

type Location int

type locations struct {
	General Location
	Left    Location
	Right   Location
	Numpad  Location
}

var Locations = locations{
	General: 0,
	Left:    1,
	Right:   2,
	Numpad:  3,
}
