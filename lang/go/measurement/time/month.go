package time

// Source of truth: domain-model/measurement/time/month.md

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidMonth is returned when a value is not a valid month of the year.
var ErrInvalidMonth = errorer.New("invalid month")
var errInvalidMonthFn = errorer.Func(ErrInvalidMonth)

// Month is a month of the year with no year or day, numbered by ISO 8601
// from January (1) to December (12).
//
// It serializes as the month's English name, e.g. "January".
type Month int

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// ZeroMonth is the zero value. It is not a valid Month.
var ZeroMonth = Month(0)

// MonthsOfYear lists every Month in ISO 8601 order, January first.
var MonthsOfYear = []Month{January, February, March, April, May, June, July, August, September, October, November, December}

var monthNames = [...]string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

// Number returns the ISO 8601 month number, 1 (January) through 12 (December).
func (this Month) Number() int {
	return int(this)
}

// Month returns the month as a time.Month.
// The result for an invalid Month is undefined.
func (this Month) Month() time.Month {
	return time.Month(this)
}

// String returns the month's full English name, e.g. "January", or a
// placeholder such as "Month(0)" for an invalid value.
func (this Month) String() string {
	if this.Validate() != nil {
		return "Month(" + strconv.Itoa(int(this)) + ")"
	}

	return monthNames[this]
}

// IsZero reports whether the Month is the zero value.
func (this Month) IsZero() bool {
	return this == ZeroMonth
}

// Validate returns an error wrapping ErrInvalidMonth if the Month is not
// between January (1) and December (12).
func (this Month) Validate() error {
	if this < January || this > December {
		return errInvalidMonthFn("%d must be between 1 and 12", int(this))
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the Month serializes as
// its name. It returns the validation error for an invalid Month.
func (this Month) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *Month) UnmarshalText(b []byte) error {
	m, err := ParseMonth(string(b))

	if err != nil {
		return err
	}

	*this = m
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the Month as a JSON
// string. It returns the validation error for an invalid Month.
func (this Month) MarshalJSON() ([]byte, error) {
	text, err := this.MarshalText()

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a month name or abbreviation. A JSON null leaves the Month
// unchanged.
func (this *Month) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidMonthFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewMonth returns the month of t, in t's own location.
func NewMonth(t time.Time) Month {
	return Month(t.Month())
}

// ParseMonth returns the Month named by s. It accepts the full name or the
// three-letter abbreviation, ignoring case ("January", "jan", "JAN").
// It returns an error wrapping ErrInvalidMonth for anything else.
func ParseMonth(s string) (Month, error) {
	for _, m := range MonthsOfYear {
		name := monthNames[m]

		if strings.EqualFold(s, name) || strings.EqualFold(s, name[:3]) {
			return m, nil
		}
	}

	return ZeroMonth, errInvalidMonthFn("%q", s)
}

// MustParseMonth is like ParseMonth but panics if s is not a month of the
// year.
func MustParseMonth(s string) Month {
	m, err := ParseMonth(s)

	if err != nil {
		panic(err)
	}

	return m
}
