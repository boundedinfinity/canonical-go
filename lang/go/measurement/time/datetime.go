// Package time implements the DateTime, Date, Time, and Duration data types
// and their ranges.
//
// Source of truth: domain-model/measurement/time/datetime.md
package time

import (
	"encoding/json"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

// ErrInvalid is returned when a value is not a valid RFC 3339 date-time.
var ErrInvalid = errorer.New("invalid datetime")
var errInvalidFn = errorer.Func(ErrInvalid)

// Zero is the zero value. It is never returned by New, Now, or Parse.
var Zero = DateTime{}

// DateTime is a date and time with a time offset, following the RFC 3339
// date-time format.
//
// It serializes as a string, e.g. "2026-11-07T09:30:00.123-05:00".
type DateTime time.Time

// Time returns the underlying time.Time.
func (this DateTime) Time() time.Time {
	return time.Time(this)
}

// String returns the RFC 3339 representation, with fractional seconds only
// when present.
func (this DateTime) String() string {
	return time.Time(this).Format(time.RFC3339Nano)
}

// IsZero reports whether the DateTime is the zero value.
func (this DateTime) IsZero() bool {
	return time.Time(this).IsZero()
}

// Equal reports whether this and other represent the same instant, even if
// their time offsets differ.
func (this DateTime) Equal(other DateTime) bool {
	return time.Time(this).Equal(time.Time(other))
}

// Validate returns an error wrapping ErrInvalid if the DateTime is the zero value.
func (this DateTime) Validate() error {
	if this.IsZero() {
		return errInvalidFn("must not be zero")
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler, so the DateTime serializes as a string.
func (this DateTime) MarshalText() ([]byte, error) {
	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *DateTime) UnmarshalText(b []byte) error {
	dt, err := Parse(string(b))

	if err != nil {
		return err
	}

	*this = dt
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the DateTime as a JSON string.
func (this DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String())
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a valid RFC 3339 date-time. A JSON null leaves the DateTime unchanged.
func (this *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// New returns the DateTime for t, keeping its time offset.
// It returns an error wrapping ErrInvalid if t is the zero time.
func New(t time.Time) (DateTime, error) {
	dt := DateTime(t)

	if err := dt.Validate(); err != nil {
		return Zero, err
	}

	return dt, nil
}

// MustNew is like New but panics if t is the zero time.
func MustNew(t time.Time) DateTime {
	dt, err := New(t)

	if err != nil {
		panic(err)
	}

	return dt
}

// Now returns the current local date and time.
func Now() DateTime {
	return DateTime(time.Now())
}

// Parse returns the DateTime represented by s.
// It returns an error wrapping ErrInvalid if s is not a valid RFC 3339
// date-time or is the zero time.
func Parse(s string) (DateTime, error) {
	t, err := time.Parse(time.RFC3339Nano, s)

	if err != nil {
		return Zero, errInvalidFn("%q: %w", s, err)
	}

	dt := DateTime(t)

	if err := dt.Validate(); err != nil {
		return Zero, errInvalidFn("%q: %w", s, err)
	}

	return dt, nil
}

// MustParse is like Parse but panics if s is not a valid RFC 3339 date-time.
func MustParse(s string) DateTime {
	dt, err := Parse(s)

	if err != nil {
		panic(err)
	}

	return dt
}
