package keymap

import "github.com/boundedinfinity/go-commoner/errorer"

type Key struct {
	Code int
}

func (k Key) Characters() ([]Character, bool) {
	var found []Character

	switch k.Code {
	case 27:
		found = []Character{Characters.Escape}
	case 32:
		found = []Character{Characters.Space}
	case 48:
		found = []Character{Characters.RightParenthesis, Characters.Digit0}
	case 49:
		found = []Character{Characters.ExclamationMark, Characters.Digit1}
	case 50:
		found = []Character{Characters.AtSign, Characters.Digit2}
	case 51:
		found = []Character{Characters.NumberSign, Characters.Digit3}
	case 52:
		found = []Character{Characters.DollarSign, Characters.Digit4}
	case 53:
		found = []Character{Characters.PercentSign, Characters.Digit5}
	case 54:
		found = []Character{Characters.Caret, Characters.Digit6}
	case 55:
		found = []Character{Characters.Ampersand, Characters.Digit7}
	case 56:
		found = []Character{Characters.Asterisk, Characters.Digit8}
	case 57:
		found = []Character{Characters.LeftParenthesis, Characters.Digit9}
	case 65:
		found = []Character{Characters.UppercaseA, Characters.LowercaseA}
	case 66:
		found = []Character{Characters.UppercaseB, Characters.LowercaseB}
	case 67:
		found = []Character{Characters.UppercaseC, Characters.LowercaseC}
	case 68:
		found = []Character{Characters.UppercaseD, Characters.LowercaseD}
	case 69:
		found = []Character{Characters.UppercaseE, Characters.LowercaseE}
	case 70:
		found = []Character{Characters.UppercaseF, Characters.LowercaseF}
	case 71:
		found = []Character{Characters.UppercaseG, Characters.LowercaseG}
	case 72:
		found = []Character{Characters.UppercaseH, Characters.LowercaseH}
	case 73:
		found = []Character{Characters.UppercaseI, Characters.LowercaseI}
	case 74:
		found = []Character{Characters.UppercaseJ, Characters.LowercaseJ}
	case 75:
		found = []Character{Characters.UppercaseK, Characters.LowercaseK}
	case 76:
		found = []Character{Characters.UppercaseL, Characters.LowercaseL}
	case 77:
		found = []Character{Characters.UppercaseM, Characters.LowercaseM}
	case 78:
		found = []Character{Characters.UppercaseN, Characters.LowercaseN}
	case 79:
		found = []Character{Characters.UppercaseO, Characters.LowercaseO}
	case 80:
		found = []Character{Characters.UppercaseP, Characters.LowercaseP}
	case 81:
		found = []Character{Characters.UppercaseQ, Characters.LowercaseQ}
	case 82:
		found = []Character{Characters.UppercaseR, Characters.LowercaseR}
	case 83:
		found = []Character{Characters.UppercaseS, Characters.LowercaseS}
	case 84:
		found = []Character{Characters.UppercaseT, Characters.LowercaseT}
	case 85:
		found = []Character{Characters.UppercaseU, Characters.LowercaseU}
	case 86:
		found = []Character{Characters.UppercaseV, Characters.LowercaseV}
	case 87:
		found = []Character{Characters.UppercaseW, Characters.LowercaseW}
	case 88:
		found = []Character{Characters.UppercaseX, Characters.LowercaseX}
	case 89:
		found = []Character{Characters.UppercaseY, Characters.LowercaseY}
	case 90:
		found = []Character{Characters.UppercaseZ, Characters.LowercaseZ}
	case 96:
		found = []Character{Characters.Digit0}
	case 97:
		found = []Character{Characters.Digit1}
	case 98:
		found = []Character{Characters.Digit2}
	case 99:
		found = []Character{Characters.Digit3}
	case 100:
		found = []Character{Characters.Digit4}
	case 101:
		found = []Character{Characters.Digit5}
	case 102:
		found = []Character{Characters.Digit6}
	case 103:
		found = []Character{Characters.Digit7}
	case 104:
		found = []Character{Characters.Digit8}
	case 105:
		found = []Character{Characters.Digit9}
	case 106:
		found = []Character{Characters.Asterisk}
	case 107:
		found = []Character{Characters.Plus}
	case 109:
		found = []Character{Characters.Hyphen}
	case 110:
		found = []Character{Characters.Period}
	case 111:
		found = []Character{Characters.Slash}
	case 186:
		found = []Character{Characters.Colon, Characters.Semicolon}
	case 187:
		found = []Character{Characters.Plus, Characters.Equal}
	case 188:
		found = []Character{Characters.LessThan, Characters.Comma}
	case 189:
		found = []Character{Characters.Underscore, Characters.Hyphen}
	case 190:
		found = []Character{Characters.GreaterThan, Characters.Period}
	case 191:
		found = []Character{Characters.QuestionMark, Characters.Slash}
	case 192:
		found = []Character{Characters.Tilde, Characters.GraveAccent}
	case 219:
		found = []Character{Characters.LeftBrace, Characters.LeftBracket}
	case 220:
		found = []Character{Characters.VerticalBar, Characters.Backslash}
	case 221:
		found = []Character{Characters.RightBrace, Characters.RightBracket}
	case 222:
		found = []Character{Characters.DoubleQuote, Characters.Apostrophe}
	default:
		found = []Character{}
	}

	return found, len(found) > 0
}

var Keys = keys{
	Backspace:      Key{Code: 8},
	Tab:            Key{Code: 9},
	Enter:          Key{Code: 13},
	Shift:          Key{Code: 16},
	Control:        Key{Code: 17},
	Alt:            Key{Code: 18},
	Pause:          Key{Code: 19},
	CapsLock:       Key{Code: 20},
	Escape:         Key{Code: 27},
	Space:          Key{Code: 32},
	PageUp:         Key{Code: 33},
	PageDown:       Key{Code: 34},
	End:            Key{Code: 35},
	Home:           Key{Code: 36},
	ArrowLeft:      Key{Code: 37},
	ArrowUp:        Key{Code: 38},
	ArrowRight:     Key{Code: 39},
	ArrowDown:      Key{Code: 40},
	PrintScreen:    Key{Code: 44},
	Insert:         Key{Code: 45},
	Delete:         Key{Code: 46},
	Digit0:         Key{Code: 48},
	Digit1:         Key{Code: 49},
	Digit2:         Key{Code: 50},
	Digit3:         Key{Code: 51},
	Digit4:         Key{Code: 52},
	Digit5:         Key{Code: 53},
	Digit6:         Key{Code: 54},
	Digit7:         Key{Code: 55},
	Digit8:         Key{Code: 56},
	Digit9:         Key{Code: 57},
	KeyA:           Key{Code: 65},
	KeyB:           Key{Code: 66},
	KeyC:           Key{Code: 67},
	KeyD:           Key{Code: 68},
	KeyE:           Key{Code: 69},
	KeyF:           Key{Code: 70},
	KeyG:           Key{Code: 71},
	KeyH:           Key{Code: 72},
	KeyI:           Key{Code: 73},
	KeyJ:           Key{Code: 74},
	KeyK:           Key{Code: 75},
	KeyL:           Key{Code: 76},
	KeyM:           Key{Code: 77},
	KeyN:           Key{Code: 78},
	KeyO:           Key{Code: 79},
	KeyP:           Key{Code: 80},
	KeyQ:           Key{Code: 81},
	KeyR:           Key{Code: 82},
	KeyS:           Key{Code: 83},
	KeyT:           Key{Code: 84},
	KeyU:           Key{Code: 85},
	KeyV:           Key{Code: 86},
	KeyW:           Key{Code: 87},
	KeyX:           Key{Code: 88},
	KeyY:           Key{Code: 89},
	KeyZ:           Key{Code: 90},
	Numpad0:        Key{Code: 96},
	Numpad1:        Key{Code: 97},
	Numpad2:        Key{Code: 98},
	Numpad3:        Key{Code: 99},
	Numpad4:        Key{Code: 100},
	Numpad5:        Key{Code: 101},
	Numpad6:        Key{Code: 102},
	Numpad7:        Key{Code: 103},
	Numpad8:        Key{Code: 104},
	Numpad9:        Key{Code: 105},
	NumpadMultiply: Key{Code: 106},
	NumpadAdd:      Key{Code: 107},
	NumpadSubtract: Key{Code: 109},
	NumpadDecimal:  Key{Code: 110},
	NumpadDivide:   Key{Code: 111},
	Function01:     Key{Code: 112},
	Function02:     Key{Code: 113},
	Function03:     Key{Code: 114},
	Function04:     Key{Code: 115},
	Function05:     Key{Code: 116},
	Function06:     Key{Code: 117},
	Function07:     Key{Code: 118},
	Function08:     Key{Code: 119},
	Function09:     Key{Code: 120},
	Function10:     Key{Code: 121},
	Function11:     Key{Code: 122},
	Function12:     Key{Code: 123},
	NumLock:        Key{Code: 144},
	ScrollLock:     Key{Code: 145},
	Semicolon:      Key{Code: 186},
	Equal:          Key{Code: 187},
	Comma:          Key{Code: 188},
	Minus:          Key{Code: 189},
	Period:         Key{Code: 190},
	Slash:          Key{Code: 191},
	Backquote:      Key{Code: 192},
	BracketLeft:    Key{Code: 219},
	Backslash:      Key{Code: 220},
	BracketRight:   Key{Code: 221},
	Quote:          Key{Code: 222},
}

type keys struct {
	Backspace      Key
	Tab            Key
	Enter          Key
	Shift          Key
	Control        Key
	Alt            Key
	Pause          Key
	CapsLock       Key
	Escape         Key
	Space          Key
	PageUp         Key
	PageDown       Key
	End            Key
	Home           Key
	ArrowLeft      Key
	ArrowUp        Key
	ArrowRight     Key
	ArrowDown      Key
	PrintScreen    Key
	Insert         Key
	Delete         Key
	Digit0         Key
	Digit1         Key
	Digit2         Key
	Digit3         Key
	Digit4         Key
	Digit5         Key
	Digit6         Key
	Digit7         Key
	Digit8         Key
	Digit9         Key
	KeyA           Key
	KeyB           Key
	KeyC           Key
	KeyD           Key
	KeyE           Key
	KeyF           Key
	KeyG           Key
	KeyH           Key
	KeyI           Key
	KeyJ           Key
	KeyK           Key
	KeyL           Key
	KeyM           Key
	KeyN           Key
	KeyO           Key
	KeyP           Key
	KeyQ           Key
	KeyR           Key
	KeyS           Key
	KeyT           Key
	KeyU           Key
	KeyV           Key
	KeyW           Key
	KeyX           Key
	KeyY           Key
	KeyZ           Key
	Numpad0        Key
	Numpad1        Key
	Numpad2        Key
	Numpad3        Key
	Numpad4        Key
	Numpad5        Key
	Numpad6        Key
	Numpad7        Key
	Numpad8        Key
	Numpad9        Key
	NumpadMultiply Key
	NumpadAdd      Key
	NumpadSubtract Key
	NumpadDecimal  Key
	NumpadDivide   Key
	Function01     Key
	Function02     Key
	Function03     Key
	Function04     Key
	Function05     Key
	Function06     Key
	Function07     Key
	Function08     Key
	Function09     Key
	Function10     Key
	Function11     Key
	Function12     Key
	NumLock        Key
	ScrollLock     Key
	Semicolon      Key
	Equal          Key
	Comma          Key
	Minus          Key
	Period         Key
	Slash          Key
	Backquote      Key
	BracketLeft    Key
	Backslash      Key
	BracketRight   Key
	Quote          Key
}

func (this keys) All() []Key {
	return []Key{
		this.Backspace,
		this.Tab,
		this.Enter,
		this.Shift,
		this.Control,
		this.Alt,
		this.Pause,
		this.CapsLock,
		this.Escape,
		this.Space,
		this.PageUp,
		this.PageDown,
		this.End,
		this.Home,
		this.ArrowLeft,
		this.ArrowUp,
		this.ArrowRight,
		this.ArrowDown,
		this.PrintScreen,
		this.Insert,
		this.Delete,
		this.Digit0,
		this.Digit1,
		this.Digit2,
		this.Digit3,
		this.Digit4,
		this.Digit5,
		this.Digit6,
		this.Digit7,
		this.Digit8,
		this.Digit9,
		this.KeyA,
		this.KeyB,
		this.KeyC,
		this.KeyD,
		this.KeyE,
		this.KeyF,
		this.KeyG,
		this.KeyH,
		this.KeyI,
		this.KeyJ,
		this.KeyK,
		this.KeyL,
		this.KeyM,
		this.KeyN,
		this.KeyO,
		this.KeyP,
		this.KeyQ,
		this.KeyR,
		this.KeyS,
		this.KeyT,
		this.KeyU,
		this.KeyV,
		this.KeyW,
		this.KeyX,
		this.KeyY,
		this.KeyZ,
		this.Numpad0,
		this.Numpad1,
		this.Numpad2,
		this.Numpad3,
		this.Numpad4,
		this.Numpad5,
		this.Numpad6,
		this.Numpad7,
		this.Numpad8,
		this.Numpad9,
		this.NumpadMultiply,
		this.NumpadAdd,
		this.NumpadSubtract,
		this.NumpadDecimal,
		this.NumpadDivide,
		this.Function01,
		this.Function02,
		this.Function03,
		this.Function04,
		this.Function05,
		this.Function06,
		this.Function07,
		this.Function08,
		this.Function09,
		this.Function10,
		this.Function11,
		this.Function12,
		this.NumLock,
		this.ScrollLock,
		this.Semicolon,
		this.Equal,
		this.Comma,
		this.Minus,
		this.Period,
		this.Slash,
		this.Backquote,
		this.BracketLeft,
		this.Backslash,
		this.BracketRight,
		this.Quote,
	}
}

func (this keys) GetOk(code int) (Key, bool) {
	for _, key := range this.All() {
		if key.Code == code {
			return key, true
		}
	}

	return Key{}, false
}

var (
	ErrKey   = errorer.New("key error")
	errKeyFn = errorer.Func(ErrKey)
)

func (this keys) GetError(code int) (Key, error) {
	key, ok := this.GetOk(code)

	if !ok {
		return Key{}, errKeyFn("code not found: %d", code)
	}

	return key, nil
}
