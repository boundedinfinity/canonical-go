package datetime

// Source of truth: domain-model/datetime/date.md

import (
	"encoding/json"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalidDate is returned when a value is not a valid RFC 3339 full-date.
var ErrInvalidDate = errorer.New("invalid date")
var errInvalidDateFn = errorer.Func(ErrInvalidDate)

// ZeroDate is the zero value. It is never returned by NewDate, Today, or ParseDate.
var ZeroDate = Date{}

// Date is a calendar date with no time of day or time offset, following the
// RFC 3339 full-date format.
//
// It is stored at midnight UTC, so equal dates compare equal with ==.
// It serializes as a string, e.g. "2026-11-07".
type Date time.Time

// Time returns the date as a time.Time at midnight UTC.
func (this Date) Time() time.Time {
	return time.Time(this)
}

// String returns the RFC 3339 full-date representation.
func (this Date) String() string {
	return time.Time(this).Format(time.DateOnly)
}

// IsZero reports whether the Date is the zero value.
func (this Date) IsZero() bool {
	return time.Time(this).IsZero()
}

// Validate returns an error wrapping ErrInvalidDate if the Date is the zero value.
func (this Date) Validate() error {
	if this.IsZero() {
		return errInvalidDateFn("must not be zero")
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the Date serializes as a string.
func (this Date) MarshalText() ([]byte, error) {
	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *Date) UnmarshalText(b []byte) error {
	d, err := ParseDate(string(b))

	if err != nil {
		return err
	}

	*this = d
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the Date as a JSON string.
func (this Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String())
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a valid RFC 3339 full-date. A JSON null leaves the Date unchanged.
func (this *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidDateFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewDate returns the calendar date of t in t's own location, dropping the
// time of day. It returns an error wrapping ErrInvalidDate if that date is the
// zero value.
func NewDate(t time.Time) (Date, error) {
	year, month, day := t.Date()
	d := Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))

	if err := d.Validate(); err != nil {
		return ZeroDate, err
	}

	return d, nil
}

// MustNewDate is like NewDate but panics if the date is the zero value.
func MustNewDate(t time.Time) Date {
	d, err := NewDate(t)

	if err != nil {
		panic(err)
	}

	return d
}

// Today returns the current local date.
func Today() Date {
	return MustNewDate(time.Now())
}

// ParseDate returns the Date represented by s.
// It returns an error wrapping ErrInvalidDate if s is not a valid RFC 3339
// full-date or is the zero value.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)

	if err != nil {
		return ZeroDate, errInvalidDateFn("%q: %w", s, err)
	}

	d := Date(t)

	if err := d.Validate(); err != nil {
		return ZeroDate, errInvalidDateFn("%q: %w", s, err)
	}

	return d, nil
}

// MustParseDate is like ParseDate but panics if s is not a valid RFC 3339 full-date.
func MustParseDate(s string) Date {
	d, err := ParseDate(s)

	if err != nil {
		panic(err)
	}

	return d
}
