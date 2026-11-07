package time

// Source of truth: domain-model/measurement/time/datetime-range.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidDateTimeRange is wrapped by every DateTimeRange validation error.
	ErrInvalidDateTimeRange = errorer.New("invalid datetime range")
	ErrStart                = errorer.New("start")
	ErrEnd                  = errorer.New("end")
	ErrOrder                = errorer.New("start must not be after end")

	errStartFn = errorer.Func(ErrInvalidDateTimeRange, ErrStart)
	errEndFn   = errorer.Func(ErrInvalidDateTimeRange, ErrEnd)
	errOrderFn = errorer.Func(ErrInvalidDateTimeRange, ErrOrder)
)

// DateTimeRange is the range of DateTimes from Start to End, inclusive.
type DateTimeRange struct {
	Start DateTime `json:"start"`
	End   DateTime `json:"end"`
}

// Between reports whether dt is within the range, inclusive of Start and End.
func (this DateTimeRange) Between(dt DateTime) bool {
	return !this.Before(dt) && !this.After(dt)
}

// Before reports whether dt is strictly before Start.
func (this DateTimeRange) Before(dt DateTime) bool {
	return dt.Time().Before(this.Start.Time())
}

// After reports whether dt is strictly after End.
func (this DateTimeRange) After(dt DateTime) bool {
	return dt.Time().After(this.End.Time())
}

// Overlaps reports whether other shares at least one instant with the range.
// Bounds are inclusive, so ranges that only touch at an endpoint overlap.
func (this DateTimeRange) Overlaps(other DateTimeRange) bool {
	return !this.Before(other.End) && !this.After(other.Start)
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidDateTimeRange and the field's sentinel
// (ErrStart, ErrEnd, or ErrOrder). A Start or End failure also wraps
// ErrInvalid from the DateTime.
func (this DateTimeRange) Validate() error {
	var errs []error

	if err := this.Start.Validate(); err != nil {
		errs = append(errs, errStartFn("%w", err))
	}

	if err := this.End.Validate(); err != nil {
		errs = append(errs, errEndFn("%w", err))
	}

	if this.Start.Time().After(this.End.Time()) {
		errs = append(errs, errOrderFn("%s > %s", this.Start, this.End))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this DateTimeRange) MarshalJSON() ([]byte, error) {
	type alias DateTimeRange
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// DateTimeRange. A JSON null leaves the DateTimeRange unchanged.
func (this *DateTimeRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias DateTimeRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := DateTimeRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewRange returns a validated DateTimeRange from start to end.
func NewRange(start, end DateTime) (DateTimeRange, error) {
	r := DateTimeRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return DateTimeRange{}, err
	}

	return r, nil
}

// MustNewRange is like NewRange but panics if the DateTimeRange is invalid.
func MustNewRange(start, end DateTime) DateTimeRange {
	r, err := NewRange(start, end)

	if err != nil {
		panic(err)
	}

	return r
}
