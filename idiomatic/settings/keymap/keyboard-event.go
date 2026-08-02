package keymap

type KeyboardEvent struct {
	Code     int
	Location Location
	Alt      bool
	Ctrl     bool
	Meta     bool
	Shift    bool
}

func (this KeyboardEvent) ResolveKeyOk() (Key, bool) {
	return Keys.GetOk(this.Code)
}

func (this KeyboardEvent) ResolveCharacterOk() (Character, bool) {
	var character Character
	var found bool
	key, ok := Keys.GetOk(this.Code)

	if !ok {
		return character, found
	}

	switch key {
	case Keys.Escape:
		character = Characters.Escape
		found = true

	case Keys.Function01:
		character = Characters.Function01
		found = true
	case Keys.Function02:
		character = Characters.Function02
		found = true
	case Keys.Function03:
		character = Characters.Function03
		found = true
	case Keys.Function04:
		character = Characters.Function04
		found = true
	case Keys.Function05:
		character = Characters.Function05
		found = true
	case Keys.Function06:
		character = Characters.Function06
		found = true
	case Keys.Function07:
		character = Characters.Function07
		found = true
	case Keys.Function08:
		character = Characters.Function08
		found = true
	case Keys.Function09:
		character = Characters.Function09
		found = true
	case Keys.Function10:
		character = Characters.Function10
		found = true
	case Keys.Function11:
		character = Characters.Function11
		found = true
	case Keys.Function12:
		character = Characters.Function12
		found = true
	case Keys.Backquote:
		if this.Shift {
			character = Characters.Tilde
			found = true
		} else {
			character = Characters.GraveAccent
			found = true
		}

	case Keys.Digit1:
		if this.Shift {
			character = Characters.ExclamationMark
			found = true
		} else {
			character = Characters.Digit1
			found = true
		}
	case Keys.Digit2:
		if this.Shift {
			character = Characters.AtSign
			found = true
		} else {
			character = Characters.Digit2
			found = true
		}
	case Keys.Digit3:
		if this.Shift {
			character = Characters.NumberSign
			found = true
		} else {
			character = Characters.Digit3
			found = true
		}
	case Keys.Digit4:
		if this.Shift {
			character = Characters.DollarSign
			found = true
		} else {
			character = Characters.Digit4
			found = true
		}
	case Keys.Digit5:
		if this.Shift {
			character = Characters.PercentSign
			found = true
		} else {
			character = Characters.Digit5
			found = true
		}
	case Keys.Digit6:
		if this.Shift {
			character = Characters.Caret
			found = true
		} else {
			character = Characters.Digit6
			found = true
		}

	case Keys.Digit7:
		if this.Shift {
			character = Characters.Ampersand
			found = true
		} else {
			character = Characters.Digit7
			found = true
		}
	case Keys.Digit8:
		if this.Shift {
			character = Characters.Asterisk
			found = true
		} else {
			character = Characters.Digit8
			found = true
		}
	case Keys.Digit9:
		if this.Shift {
			character = Characters.LeftParenthesis
			found = true
		} else {
			character = Characters.Digit9
			found = true
		}
	case Keys.Digit0:
		if this.Shift {
			character = Characters.RightParenthesis
			found = true
		} else {
			character = Characters.Digit0
			found = true
		}
	case Keys.Minus:
		if this.Shift {
			character = Characters.Underscore
			found = true
		} else {
			character = Characters.Hyphen
			found = true
		}
	case Keys.Equal:
		if this.Shift {
			character = Characters.Plus
			found = true
		} else {
			character = Characters.Equal
			found = true
		}

	}
	return character, found
}
