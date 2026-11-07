package time

// Source of truth: domain-model/measurement/time/date-range.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidDateRange is wrapped by every DateRange validation error.
	// The field sentinels (ErrStart, ErrEnd, ErrOrder) are shared with
	// DateTimeRange.
	ErrInvalidDateRange = errorer.New("invalid date range")

	errDateRangeStartFn = errorer.Func(ErrInvalidDateRange, ErrStart)
	errDateRangeEndFn   = errorer.Func(ErrInvalidDateRange, ErrEnd)
	errDateRangeOrderFn = errorer.Func(ErrInvalidDateRange, ErrOrder)
)

// DateRange is the range of Dates from Start to End, inclusive.
type DateRange struct {
	Start Date `json:"start"`
	End   Date `json:"end"`
}

// Between reports whether d is within the range, inclusive of Start and End.
func (this DateRange) Between(d Date) bool {
	return !this.Before(d) && !this.After(d)
}

// Before reports whether d is strictly before Start.
func (this DateRange) Before(d Date) bool {
	return d.Time().Before(this.Start.Time())
}

// After reports whether d is strictly after End.
func (this DateRange) After(d Date) bool {
	return d.Time().After(this.End.Time())
}

// Overlaps reports whether other shares at least one Date with the range.
// Bounds are inclusive, so ranges that only touch at an endpoint overlap.
func (this DateRange) Overlaps(other DateRange) bool {
	return !this.Before(other.End) && !this.After(other.Start)
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidDateRange and the field's sentinel (ErrStart,
// ErrEnd, or ErrOrder). A Start or End failure also wraps ErrInvalidDate.
func (this DateRange) Validate() error {
	var errs []error

	if err := this.Start.Validate(); err != nil {
		errs = append(errs, errDateRangeStartFn("%w", err))
	}

	if err := this.End.Validate(); err != nil {
		errs = append(errs, errDateRangeEndFn("%w", err))
	}

	if this.Start.Time().After(this.End.Time()) {
		errs = append(errs, errDateRangeOrderFn("%s > %s", this.Start, this.End))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this DateRange) MarshalJSON() ([]byte, error) {
	type alias DateRange
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// DateRange. A JSON null leaves the DateRange unchanged.
func (this *DateRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias DateRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := DateRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewDateRange returns a validated DateRange from start to end.
func NewDateRange(start, end Date) (DateRange, error) {
	r := DateRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return DateRange{}, err
	}

	return r, nil
}

// MustNewDateRange is like NewDateRange but panics if the DateRange is invalid.
func MustNewDateRange(start, end Date) DateRange {
	r, err := NewDateRange(start, end)

	if err != nil {
		panic(err)
	}

	return r
}
