package time

// Source of truth: domain-model/measurement/time/day-of-week.md

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidDayOfWeek is returned when a value is not a valid day of the week.
var ErrInvalidDayOfWeek = errorer.New("invalid day of week")
var errInvalidDayOfWeekFn = errorer.Func(ErrInvalidDayOfWeek)

// DayOfWeek is a day of the week with no date, numbered by ISO 8601 from
// Monday (1) to Sunday (7).
//
// It serializes as the day's English name, e.g. "Monday".
type DayOfWeek int

const (
	Monday DayOfWeek = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// ZeroDayOfWeek is the zero value. It is not a valid DayOfWeek.
var ZeroDayOfWeek = DayOfWeek(0)

// DaysOfWeek lists every DayOfWeek in ISO 8601 order, Monday first.
var DaysOfWeek = []DayOfWeek{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}

var dayOfWeekNames = [...]string{"", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// Number returns the ISO 8601 day number, 1 (Monday) through 7 (Sunday).
func (this DayOfWeek) Number() int {
	return int(this)
}

// Weekday returns the day as a time.Weekday.
// The result for an invalid DayOfWeek is undefined.
func (this DayOfWeek) Weekday() time.Weekday {
	return time.Weekday(this % 7)
}

// String returns the day's full English name, e.g. "Monday", or a
// placeholder such as "DayOfWeek(0)" for an invalid value.
func (this DayOfWeek) String() string {
	if this.Validate() != nil {
		return "DayOfWeek(" + strconv.Itoa(int(this)) + ")"
	}

	return dayOfWeekNames[this]
}

// IsZero reports whether the DayOfWeek is the zero value.
func (this DayOfWeek) IsZero() bool {
	return this == ZeroDayOfWeek
}

// Validate returns an error wrapping ErrInvalidDayOfWeek if the DayOfWeek is
// not between Monday (1) and Sunday (7).
func (this DayOfWeek) Validate() error {
	if this < Monday || this > Sunday {
		return errInvalidDayOfWeekFn("%d must be between 1 and 7", int(this))
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the DayOfWeek serializes
// as its name. It returns the validation error for an invalid DayOfWeek.
func (this DayOfWeek) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *DayOfWeek) UnmarshalText(b []byte) error {
	d, err := ParseDayOfWeek(string(b))

	if err != nil {
		return err
	}

	*this = d
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the DayOfWeek as a JSON
// string. It returns the validation error for an invalid DayOfWeek.
func (this DayOfWeek) MarshalJSON() ([]byte, error) {
	text, err := this.MarshalText()

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a day name or abbreviation. A JSON null leaves the DayOfWeek
// unchanged.
func (this *DayOfWeek) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidDayOfWeekFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewDayOfWeek returns the day of the week of t, in t's own location.
func NewDayOfWeek(t time.Time) DayOfWeek {
	weekday := t.Weekday()

	if weekday == time.Sunday {
		return Sunday
	}

	return DayOfWeek(weekday)
}

// ParseDayOfWeek returns the DayOfWeek named by s. It accepts the full name
// or the three-letter abbreviation, ignoring case ("Monday", "mon", "MON").
// It returns an error wrapping ErrInvalidDayOfWeek for anything else.
func ParseDayOfWeek(s string) (DayOfWeek, error) {
	for _, d := range DaysOfWeek {
		name := dayOfWeekNames[d]

		if strings.EqualFold(s, name) || strings.EqualFold(s, name[:3]) {
			return d, nil
		}
	}

	return ZeroDayOfWeek, errInvalidDayOfWeekFn("%q", s)
}

// MustParseDayOfWeek is like ParseDayOfWeek but panics if s is not a day of
// the week.
func MustParseDayOfWeek(s string) DayOfWeek {
	d, err := ParseDayOfWeek(s)

	if err != nil {
		panic(err)
	}

	return d
}
