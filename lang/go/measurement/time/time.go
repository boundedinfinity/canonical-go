package time

// Source of truth: domain-model/measurement/time/time.md

import (
	"encoding/json"
	"time"

	"github.com/boundedinfinity/go-commoner/errorer"
)

const (
	// timeParseLayout accepts optional fractional seconds: time.Parse allows
	// them after a seconds field even when the layout has none.
	timeParseLayout  = "15:04:05Z07:00"
	timeFormatLayout = "15:04:05.999999999Z07:00"
)

// ErrInvalidTime is returned when a value is not a valid RFC 3339 full-time.
var ErrInvalidTime = errorer.New("invalid time")
var errInvalidTimeFn = errorer.Func(ErrInvalidTime)

// ZeroTime is the zero value, midnight UTC ("00:00:00Z"). It is a valid Time.
var ZeroTime = Time{}

// Time is a time of day with a time offset and no date, following the
// RFC 3339 full-time format.
//
// It is stored on the zero date (0001-01-01), so the zero value is midnight
// UTC. It serializes as a string, e.g. "09:30:00.123-05:00".
type Time time.Time

// Time returns the time of day as a time.Time on 0001-01-01, in the Time's
// own offset.
func (this Time) Time() time.Time {
	return time.Time(this)
}

// String returns the RFC 3339 full-time representation, with fractional
// seconds only when present.
func (this Time) String() string {
	return time.Time(this).Format(timeFormatLayout)
}

// IsZero reports whether the Time is the zero value, midnight UTC.
func (this Time) IsZero() bool {
	return time.Time(this).IsZero()
}

// Equal reports whether this and other are the same time of day in UTC, even
// if their time offsets differ. It wraps at midnight, so "23:00:00-05:00"
// equals "04:00:00Z".
func (this Time) Equal(other Time) bool {
	a, b := this.Time().UTC(), other.Time().UTC()

	return a.Hour() == b.Hour() &&
		a.Minute() == b.Minute() &&
		a.Second() == b.Second() &&
		a.Nanosecond() == b.Nanosecond()
}

// Validate always returns nil: every Time, including midnight, is valid.
// The format rules are enforced by ParseTime.
func (this Time) Validate() error {
	return nil
}

// MarshalText implements encoding.TextMarshaler, so the Time serializes as a string.
func (this Time) MarshalText() ([]byte, error) {
	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *Time) UnmarshalText(b []byte) error {
	t, err := ParseTime(string(b))

	if err != nil {
		return err
	}

	*this = t
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the Time as a JSON string.
func (this Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String())
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a valid RFC 3339 full-time. A JSON null leaves the Time unchanged.
func (this *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidTimeFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewTime returns the time of day of t, keeping its time offset and dropping
// the date.
func NewTime(t time.Time) Time {
	return Time(time.Date(1, time.January, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()))
}

// CurrentTime returns the current local time of day.
func CurrentTime() Time {
	return NewTime(time.Now())
}

// ParseTime returns the Time represented by s.
// It returns an error wrapping ErrInvalidTime if s is not a valid RFC 3339
// full-time.
func ParseTime(s string) (Time, error) {
	t, err := time.Parse(timeParseLayout, s)

	if err != nil {
		return ZeroTime, errInvalidTimeFn("%q: %w", s, err)
	}

	return NewTime(t), nil
}

// MustParseTime is like ParseTime but panics if s is not a valid RFC 3339 full-time.
func MustParseTime(s string) Time {
	t, err := ParseTime(s)

	if err != nil {
		panic(err)
	}

	return t
}
