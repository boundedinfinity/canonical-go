package time

// Source of truth: domain-model/measurement/time/hour-of-day.md

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidHourOfDay is returned when a value is not a valid hour of the day.
var ErrInvalidHourOfDay = errorer.New("invalid hour of day")
var errInvalidHourOfDayFn = errorer.Func(ErrInvalidHourOfDay)

var hourOfDayText = regexp.MustCompile(`^[0-9]{1,2}$`)

// ZeroHourOfDay is the zero value, midnight (0). It is a valid HourOfDay.
var ZeroHourOfDay = HourOfDay(0)

// HourOfDay is an hour of the day on a 24-hour clock, from 0 (midnight to
// 1 AM) through 23 (11 PM to midnight), with no date, minutes, or offset.
//
// It serializes as a JSON number, e.g. 14.
type HourOfDay int

// Number returns the hour, 0 through 23.
func (this HourOfDay) Number() int {
	return int(this)
}

// String returns the decimal hour, e.g. "14".
func (this HourOfDay) String() string {
	return strconv.Itoa(int(this))
}

// IsZero reports whether the HourOfDay is the zero value, midnight.
func (this HourOfDay) IsZero() bool {
	return this == ZeroHourOfDay
}

// Validate returns an error wrapping ErrInvalidHourOfDay if the HourOfDay is
// not between 0 and 23.
func (this HourOfDay) Validate() error {
	if this < 0 || this > 23 {
		return errInvalidHourOfDayFn("%d must be between 0 and 23", int(this))
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the HourOfDay can be a
// map key. It returns the validation error for an invalid HourOfDay.
func (this HourOfDay) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *HourOfDay) UnmarshalText(b []byte) error {
	h, err := ParseHourOfDay(string(b))

	if err != nil {
		return err
	}

	*this = h
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the HourOfDay as a JSON
// number. It returns the validation error for an invalid HourOfDay.
func (this HourOfDay) MarshalJSON() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(int(this))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON integer from
// 0 through 23. A JSON null leaves the HourOfDay unchanged.
func (this *HourOfDay) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var n int

	if err := json.Unmarshal(data, &n); err != nil {
		return errInvalidHourOfDayFn("%w", err)
	}

	h := HourOfDay(n)

	if err := h.Validate(); err != nil {
		return err
	}

	*this = h
	return nil
}

// NewHourOfDay returns the hour of t, in t's own location.
func NewHourOfDay(t time.Time) HourOfDay {
	return HourOfDay(t.Hour())
}

// ParseHourOfDay returns the HourOfDay represented by s, a decimal hour from
// "0" through "23", with at most one leading zero ("07").
// It returns an error wrapping ErrInvalidHourOfDay for anything else.
func ParseHourOfDay(s string) (HourOfDay, error) {
	if !hourOfDayText.MatchString(s) {
		return ZeroHourOfDay, errInvalidHourOfDayFn("%q", s)
	}

	n, _ := strconv.Atoi(s)
	h := HourOfDay(n)

	if err := h.Validate(); err != nil {
		return ZeroHourOfDay, errInvalidHourOfDayFn("%q: %w", s, err)
	}

	return h, nil
}

// MustParseHourOfDay is like ParseHourOfDay but panics if s is not a valid
// hour.
func MustParseHourOfDay(s string) HourOfDay {
	h, err := ParseHourOfDay(s)

	if err != nil {
		panic(err)
	}

	return h
}
