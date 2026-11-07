package time

// Source of truth: domain-model/measurement/time/clock-format.md

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidClockFormat is returned when a value is not a valid clock format.
var ErrInvalidClockFormat = errorer.New("invalid clock format")
var errInvalidClockFormatFn = errorer.Func(ErrInvalidClockFormat)

// ClockFormat is how times of day are displayed: on a 12-hour clock (2 PM)
// or a 24-hour clock (14:00).
//
// It serializes as its name, e.g. "24-hour".
type ClockFormat int

const (
	ClockFormat12Hour ClockFormat = iota + 1
	ClockFormat24Hour
)

// ZeroClockFormat is the zero value. It is not a valid ClockFormat.
var ZeroClockFormat = ClockFormat(0)

// ClockFormats lists every ClockFormat.
var ClockFormats = []ClockFormat{ClockFormat12Hour, ClockFormat24Hour}

var clockFormatNames = [...]string{"", "12-hour", "24-hour"}

// String returns the format's name, e.g. "24-hour", or a placeholder such as
// "ClockFormat(0)" for an invalid value.
func (this ClockFormat) String() string {
	if this.Validate() != nil {
		return "ClockFormat(" + strconv.Itoa(int(this)) + ")"
	}

	return clockFormatNames[this]
}

// IsZero reports whether the ClockFormat is the zero value.
func (this ClockFormat) IsZero() bool {
	return this == ZeroClockFormat
}

// Validate returns an error wrapping ErrInvalidClockFormat if the
// ClockFormat is not one of ClockFormats.
func (this ClockFormat) Validate() error {
	if this < ClockFormat12Hour || this > ClockFormat24Hour {
		return errInvalidClockFormatFn("%d is not a clock format", int(this))
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the ClockFormat
// serializes as its name. It returns the validation error for an invalid
// ClockFormat.
func (this ClockFormat) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *ClockFormat) UnmarshalText(b []byte) error {
	c, err := ParseClockFormat(string(b))

	if err != nil {
		return err
	}

	*this = c
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the ClockFormat as a JSON
// string. It returns the validation error for an invalid ClockFormat.
func (this ClockFormat) MarshalJSON() ([]byte, error) {
	text, err := this.MarshalText()

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a clock format name. A JSON null leaves the ClockFormat
// unchanged.
func (this *ClockFormat) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidClockFormatFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// ParseClockFormat returns the ClockFormat named by s, ignoring case
// ("24-hour", "24-HOUR"). It returns an error wrapping ErrInvalidClockFormat
// for anything else.
func ParseClockFormat(s string) (ClockFormat, error) {
	for _, c := range ClockFormats {
		if strings.EqualFold(s, c.String()) {
			return c, nil
		}
	}

	return ZeroClockFormat, errInvalidClockFormatFn("%q", s)
}

// MustParseClockFormat is like ParseClockFormat but panics if s is not a
// clock format.
func MustParseClockFormat(s string) ClockFormat {
	c, err := ParseClockFormat(s)

	if err != nil {
		panic(err)
	}

	return c
}
