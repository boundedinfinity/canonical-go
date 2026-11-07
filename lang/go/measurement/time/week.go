package time

// Source of truth: domain-model/measurement/time/week.md

import (
	"cmp"
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	// ErrInvalidWeek is wrapped by every Week validation error. The field
	// sentinel ErrStart is shared with the ranges.
	ErrInvalidWeek = errorer.New("invalid week")

	errWeekStartFn = errorer.Func(ErrInvalidWeek, ErrStart)
)

// DefaultWeek is the week starting on Monday. It behaves the same as the zero
// Week.
var DefaultWeek = Week{Start: optioner.Some(Monday)}

// Week is all the days of a week, in order, starting on Start. When Start is
// not present, the week starts on Monday.
//
// The operations are undefined for an invalid Week.
type Week struct {
	Start optioner.Option[DayOfWeek] `json:"start,omitzero"`
}

// First returns the first day of the week, Start or Monday when not present.
func (this Week) First() DayOfWeek {
	return this.Start.OrElse(Monday)
}

// Last returns the last day of the week, the day before First.
func (this Week) Last() DayOfWeek {
	return (this.First()+5)%7 + 1
}

// Days returns all 7 days in order, starting with First.
func (this Week) Days() []DayOfWeek {
	first := this.First()
	days := make([]DayOfWeek, 7)

	for i := range days {
		days[i] = (first-1+DayOfWeek(i))%7 + 1
	}

	return days
}

// Position returns d's place in the week, 1 (First) through 7 (Last), or 0
// if d is invalid.
func (this Week) Position(d DayOfWeek) int {
	if d.Validate() != nil {
		return 0
	}

	return int((d-this.First()+7)%7) + 1
}

// Compare returns -1, 0, or 1 as a comes before, is the same as, or comes
// after b in the week.
func (this Week) Compare(a, b DayOfWeek) int {
	return cmp.Compare(this.Position(a), this.Position(b))
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidWeek and the field's sentinel (ErrStart).
// A Start failure also wraps ErrInvalidDayOfWeek.
func (this Week) Validate() error {
	var errs []error

	if this.Start.Defined() {
		if err := this.Start.Get().Validate(); err != nil {
			errs = append(errs, errWeekStartFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Week) MarshalJSON() ([]byte, error) {
	type alias Week
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Week.
// A JSON null leaves the Week unchanged.
func (this *Week) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Week
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	week := Week(decoded)

	if err := week.Validate(); err != nil {
		return err
	}

	*this = week
	return nil
}

// NewWeek returns a validated Week starting on start. A zero start is treated
// as not present, so the week starts on Monday.
func NewWeek(start DayOfWeek) (Week, error) {
	week := Week{
		Start: optioner.OfZero(start),
	}

	if err := week.Validate(); err != nil {
		return Week{}, err
	}

	return week, nil
}

// MustNewWeek is like NewWeek but panics if the Week is invalid.
func MustNewWeek(start DayOfWeek) Week {
	week, err := NewWeek(start)

	if err != nil {
		panic(err)
	}

	return week
}
