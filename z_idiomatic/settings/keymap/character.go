package keymap

import (
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

type Character struct {
	Char        string
	Description string
	Aliases     []string
}

var Characters = characters{
	UppercaseA:       Character{Char: "A", Description: "Uppercase A"},
	UppercaseB:       Character{Char: "B", Description: "Uppercase B"},
	UppercaseC:       Character{Char: "C", Description: "Uppercase C"},
	UppercaseD:       Character{Char: "D", Description: "Uppercase D"},
	UppercaseE:       Character{Char: "E", Description: "Uppercase E"},
	UppercaseF:       Character{Char: "F", Description: "Uppercase F"},
	UppercaseG:       Character{Char: "G", Description: "Uppercase G"},
	UppercaseH:       Character{Char: "H", Description: "Uppercase H"},
	UppercaseI:       Character{Char: "I", Description: "Uppercase I"},
	UppercaseJ:       Character{Char: "J", Description: "Uppercase J"},
	UppercaseK:       Character{Char: "K", Description: "Uppercase K"},
	UppercaseL:       Character{Char: "L", Description: "Uppercase L"},
	UppercaseM:       Character{Char: "M", Description: "Uppercase M"},
	UppercaseN:       Character{Char: "N", Description: "Uppercase N"},
	UppercaseO:       Character{Char: "O", Description: "Uppercase O"},
	UppercaseP:       Character{Char: "P", Description: "Uppercase P"},
	UppercaseQ:       Character{Char: "Q", Description: "Uppercase Q"},
	UppercaseR:       Character{Char: "R", Description: "Uppercase R"},
	UppercaseS:       Character{Char: "S", Description: "Uppercase S"},
	UppercaseT:       Character{Char: "T", Description: "Uppercase T"},
	UppercaseU:       Character{Char: "U", Description: "Uppercase U"},
	UppercaseV:       Character{Char: "V", Description: "Uppercase V"},
	UppercaseW:       Character{Char: "W", Description: "Uppercase W"},
	UppercaseX:       Character{Char: "X", Description: "Uppercase X"},
	UppercaseY:       Character{Char: "Y", Description: "Uppercase Y"},
	UppercaseZ:       Character{Char: "Z", Description: "Uppercase Z"},
	LowercaseA:       Character{Char: "a", Description: "Lowercase A"},
	LowercaseB:       Character{Char: "b", Description: "Lowercase B"},
	LowercaseC:       Character{Char: "c", Description: "Lowercase C"},
	LowercaseD:       Character{Char: "d", Description: "Lowercase D"},
	LowercaseE:       Character{Char: "e", Description: "Lowercase E"},
	LowercaseF:       Character{Char: "f", Description: "Lowercase F"},
	LowercaseG:       Character{Char: "g", Description: "Lowercase G"},
	LowercaseH:       Character{Char: "h", Description: "Lowercase H"},
	LowercaseI:       Character{Char: "i", Description: "Lowercase I"},
	LowercaseJ:       Character{Char: "j", Description: "Lowercase J"},
	LowercaseK:       Character{Char: "k", Description: "Lowercase K"},
	LowercaseL:       Character{Char: "l", Description: "Lowercase L"},
	LowercaseM:       Character{Char: "m", Description: "Lowercase M"},
	LowercaseN:       Character{Char: "n", Description: "Lowercase N"},
	LowercaseO:       Character{Char: "o", Description: "Lowercase O"},
	LowercaseP:       Character{Char: "p", Description: "Lowercase P"},
	LowercaseQ:       Character{Char: "q", Description: "Lowercase Q"},
	LowercaseR:       Character{Char: "r", Description: "Lowercase R"},
	LowercaseS:       Character{Char: "s", Description: "Lowercase S"},
	LowercaseT:       Character{Char: "t", Description: "Lowercase T"},
	LowercaseU:       Character{Char: "u", Description: "Lowercase U"},
	LowercaseV:       Character{Char: "v", Description: "Lowercase V"},
	LowercaseW:       Character{Char: "w", Description: "Lowercase W"},
	LowercaseX:       Character{Char: "x", Description: "Lowercase X"},
	LowercaseY:       Character{Char: "y", Description: "Lowercase Y"},
	LowercaseZ:       Character{Char: "z", Description: "Lowercase Z"},
	Digit0:           Character{Char: "0", Description: "Digit 0"},
	Digit1:           Character{Char: "1", Description: "Digit 1"},
	Digit2:           Character{Char: "2", Description: "Digit 2"},
	Digit3:           Character{Char: "3", Description: "Digit 3"},
	Digit4:           Character{Char: "4", Description: "Digit 4"},
	Digit5:           Character{Char: "5", Description: "Digit 5"},
	Digit6:           Character{Char: "6", Description: "Digit 6"},
	Digit7:           Character{Char: "7", Description: "Digit 7"},
	Digit8:           Character{Char: "8", Description: "Digit 8"},
	Digit9:           Character{Char: "9", Description: "Digit 9"},
	Escape:           Character{Char: "Escape", Description: "Escape Key"},
	Space:            Character{Char: " ", Description: "Space"},
	ExclamationMark:  Character{Char: "!", Description: "Exclamation Mark", Aliases: []string{"Bang", "Exclamation Point"}},
	DoubleQuote:      Character{Char: "\"", Description: "Double Quote", Aliases: []string{"Quotation Mark"}},
	NumberSign:       Character{Char: "#", Description: "Number Sign", Aliases: []string{"Hash", "Pound Sign"}},
	DollarSign:       Character{Char: "$", Description: "Dollar Sign"},
	PercentSign:      Character{Char: "%", Description: "Percent Sign"},
	Ampersand:        Character{Char: "&", Description: "Ampersand", Aliases: []string{"And Sign"}},
	Apostrophe:       Character{Char: "'", Description: "Apostrophe", Aliases: []string{"Single Quote"}},
	LeftParenthesis:  Character{Char: "(", Description: "Left Parenthesis"},
	RightParenthesis: Character{Char: ")", Description: "Right Parenthesis"},
	Asterisk:         Character{Char: "*", Description: "Asterisk", Aliases: []string{"Star"}},
	Plus:             Character{Char: "+", Description: "Plus Sign"},
	Comma:            Character{Char: ",", Description: "Comma"},
	Hyphen:           Character{Char: "-", Description: "Hyphen", Aliases: []string{"Minus Sign", "Dash"}},
	Period:           Character{Char: ".", Description: "Period", Aliases: []string{"Full Stop", "Dot", "Decimal Point", "Point", "Period"}},
	Slash:            Character{Char: "/", Description: "Slash", Aliases: []string{"Forward Slash", "Solidus"}},
	Colon:            Character{Char: ":", Description: "Colon"},
	Semicolon:        Character{Char: ";", Description: "Semicolon"},
	LessThan:         Character{Char: "<", Description: "Less Than Sign", Aliases: []string{"Left Angle Bracket"}},
	Equal:            Character{Char: "=", Description: "Equals Sign"},
	GreaterThan:      Character{Char: ">", Description: "Greater Than Sign", Aliases: []string{"Right Angle Bracket"}},
	QuestionMark:     Character{Char: "?", Description: "Question Mark"},
	AtSign:           Character{Char: "@", Description: "At Sign", Aliases: []string{"Commercial At", "Email Sign"}},
	LeftBracket:      Character{Char: "[", Description: "Left Bracket", Aliases: []string{"Opening Bracket", "Left Square Bracket"}},
	Backslash:        Character{Char: "\\", Description: "Backslash", Aliases: []string{"Reverse Solidus"}},
	RightBracket:     Character{Char: "]", Description: "Right Bracket", Aliases: []string{"Closing Bracket", "Right Square Bracket"}},
	Caret:            Character{Char: "^", Description: "Caret", Aliases: []string{"Circumflex", "Hat"}},
	Underscore:       Character{Char: "_", Description: "Underscore", Aliases: []string{"Low Line"}},
	GraveAccent:      Character{Char: "`", Description: "Grave Accent", Aliases: []string{"Backtick", "Backquote"}},
	LeftBrace:        Character{Char: "{", Description: "Left Brace", Aliases: []string{"Left Curly Brace", "Left Curly Bracket"}},
	VerticalBar:      Character{Char: "|", Description: "Vertical Bar", Aliases: []string{"Pipe", "Vertical Line"}},
	RightBrace:       Character{Char: "}", Description: "Right Brace", Aliases: []string{"Right Curly Brace", "Right Curly Bracket"}},
	Tilde:            Character{Char: "~", Description: "Tilde", Aliases: []string{"Swizzle"}},
	UpArrow:          Character{Char: "↑", Description: "Up Arrow", Aliases: []string{"Arrow Up", "Up"}},
	DownArrow:        Character{Char: "↓", Description: "Down Arrow", Aliases: []string{"Arrow Down", "Down"}},
	LeftArrow:        Character{Char: "←", Description: "Left Arrow", Aliases: []string{"Arrow Left", "Left"}},
	RightArrow:       Character{Char: "→", Description: "Right Arrow", Aliases: []string{"Arrow Right", "Right"}},
	Function01:       Character{Char: "F1", Description: "Function Key 1", Aliases: []string{"F1"}},
	Function02:       Character{Char: "F2", Description: "Function Key 2", Aliases: []string{"F2"}},
	Function03:       Character{Char: "F3", Description: "Function Key 3", Aliases: []string{"F3"}},
	Function04:       Character{Char: "F4", Description: "Function Key 4", Aliases: []string{"F4"}},
	Function05:       Character{Char: "F5", Description: "Function Key 5", Aliases: []string{"F5"}},
	Function06:       Character{Char: "F6", Description: "Function Key 6", Aliases: []string{"F6"}},
	Function07:       Character{Char: "F7", Description: "Function Key 7", Aliases: []string{"F7"}},
	Function08:       Character{Char: "F8", Description: "Function Key 8", Aliases: []string{"F8"}},
	Function09:       Character{Char: "F9", Description: "Function Key 9", Aliases: []string{"F9"}},
	Function10:       Character{Char: "F10", Description: "Function Key 10", Aliases: []string{"F10"}},
	Function11:       Character{Char: "F11", Description: "Function Key 11", Aliases: []string{"F11"}},
	Function12:       Character{Char: "F12", Description: "Function Key 12", Aliases: []string{"F12"}},
}

type characters struct {
	parseMap         map[string]*Character
	UppercaseA       Character
	UppercaseB       Character
	UppercaseC       Character
	UppercaseD       Character
	UppercaseE       Character
	UppercaseF       Character
	UppercaseG       Character
	UppercaseH       Character
	UppercaseI       Character
	UppercaseJ       Character
	UppercaseK       Character
	UppercaseL       Character
	UppercaseM       Character
	UppercaseN       Character
	UppercaseO       Character
	UppercaseP       Character
	UppercaseQ       Character
	UppercaseR       Character
	UppercaseS       Character
	UppercaseT       Character
	UppercaseU       Character
	UppercaseV       Character
	UppercaseW       Character
	UppercaseX       Character
	UppercaseY       Character
	UppercaseZ       Character
	LowercaseA       Character
	LowercaseB       Character
	LowercaseC       Character
	LowercaseD       Character
	LowercaseE       Character
	LowercaseF       Character
	LowercaseG       Character
	LowercaseH       Character
	LowercaseI       Character
	LowercaseJ       Character
	LowercaseK       Character
	LowercaseL       Character
	LowercaseM       Character
	LowercaseN       Character
	LowercaseO       Character
	LowercaseP       Character
	LowercaseQ       Character
	LowercaseR       Character
	LowercaseS       Character
	LowercaseT       Character
	LowercaseU       Character
	LowercaseV       Character
	LowercaseW       Character
	LowercaseX       Character
	LowercaseY       Character
	LowercaseZ       Character
	Digit0           Character
	Digit1           Character
	Digit2           Character
	Digit3           Character
	Digit4           Character
	Digit5           Character
	Digit6           Character
	Digit7           Character
	Digit8           Character
	Digit9           Character
	Escape           Character
	Space            Character
	ExclamationMark  Character
	DoubleQuote      Character
	NumberSign       Character
	DollarSign       Character
	PercentSign      Character
	Ampersand        Character
	Apostrophe       Character
	LeftParenthesis  Character
	RightParenthesis Character
	Asterisk         Character
	Plus             Character
	Comma            Character
	Hyphen           Character
	Period           Character
	Slash            Character
	Colon            Character
	Semicolon        Character
	LessThan         Character
	Equal            Character
	GreaterThan      Character
	QuestionMark     Character
	AtSign           Character
	LeftBracket      Character
	Backslash        Character
	RightBracket     Character
	Caret            Character
	Underscore       Character
	GraveAccent      Character
	LeftBrace        Character
	VerticalBar      Character
	RightBrace       Character
	Tilde            Character
	UpArrow          Character
	DownArrow        Character
	LeftArrow        Character
	RightArrow       Character
	Function01       Character
	Function02       Character
	Function03       Character
	Function04       Character
	Function05       Character
	Function06       Character
	Function07       Character
	Function08       Character
	Function09       Character
	Function10       Character
	Function11       Character
	Function12       Character
}

func (this characters) All() []Character {
	return []Character{
		this.UppercaseA,
		this.UppercaseB,
		this.UppercaseC,
		this.UppercaseD,
		this.UppercaseE,
		this.UppercaseF,
		this.UppercaseG,
		this.UppercaseH,
		this.UppercaseI,
		this.UppercaseJ,
		this.UppercaseK,
		this.UppercaseL,
		this.UppercaseM,
		this.UppercaseN,
		this.UppercaseO,
		this.UppercaseP,
		this.UppercaseQ,
		this.UppercaseR,
		this.UppercaseS,
		this.UppercaseT,
		this.UppercaseU,
		this.UppercaseV,
		this.UppercaseW,
		this.UppercaseX,
		this.UppercaseY,
		this.UppercaseZ,
		this.LowercaseA,
		this.LowercaseB,
		this.LowercaseC,
		this.LowercaseD,
		this.LowercaseE,
		this.LowercaseF,
		this.LowercaseG,
		this.LowercaseH,
		this.LowercaseI,
		this.LowercaseJ,
		this.LowercaseK,
		this.LowercaseL,
		this.LowercaseM,
		this.LowercaseN,
		this.LowercaseO,
		this.LowercaseP,
		this.LowercaseQ,
		this.LowercaseR,
		this.LowercaseS,
		this.LowercaseT,
		this.LowercaseU,
		this.LowercaseV,
		this.LowercaseW,
		this.LowercaseX,
		this.LowercaseY,
		this.LowercaseZ,
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
		this.UpArrow,
		this.DownArrow,
		this.LeftArrow,
		this.RightArrow,
		this.Escape,
		this.Space,
		this.ExclamationMark,
		this.DoubleQuote,
		this.NumberSign,
		this.DollarSign,
		this.PercentSign,
		this.Ampersand,
		this.Apostrophe,
		this.LeftParenthesis,
		this.RightParenthesis,
		this.Asterisk,
		this.Plus,
		this.Comma,
		this.Hyphen,
		this.Period,
		this.Slash,
		this.Colon,
		this.Semicolon,
		this.LessThan,
		this.Equal,
		this.GreaterThan,
		this.QuestionMark,
		this.AtSign,
		this.LeftBracket,
		this.Backslash,
		this.RightBracket,
		this.Caret,
		this.Underscore,
		this.GraveAccent,
		this.LeftBrace,
		this.VerticalBar,
		this.RightBrace,
		this.Tilde,
	}
}

var (
	ErrCharacter         = errorer.New("character error")
	ErrCharacterNotFound = errorer.New("not found")
	errCharacterNotFound = errorer.Func(ErrCharacter, ErrCharacterNotFound)
)

func (this characters) initParseMap() {
	if this.parseMap != nil {
		return
	}

	for _, character := range this.All() {
		this.parseMap[character.Char] = &character
		for _, alias := range character.Aliases {
			this.parseMap[stringer.ToLower(alias)] = &character
		}
	}
}

func (this characters) ParseOk(input string) (Character, bool) {
	this.initParseMap()
	character, found := this.parseMap[input]

	if !found {
		character, found = this.parseMap[stringer.ToLower(input)]
	}

	return *character, found
}

func (this characters) Parse(input string) (Character, error) {
	character, ok := this.ParseOk(input)

	if !ok {
		return character, errCharacterNotFound("%s", input)
	}

	return character, nil
}
