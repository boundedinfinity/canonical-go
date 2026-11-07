package datetime

// Source of truth: domain-model/datetime/time-range.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidTimeRange is wrapped by every TimeRange validation error.
	// The field sentinels (ErrStart, ErrEnd, ErrOrder) are shared with
	// DateTimeRange.
	ErrInvalidTimeRange = errorer.New("invalid time range")

	errTimeRangeStartFn = errorer.Func(ErrInvalidTimeRange, ErrStart)
	errTimeRangeEndFn   = errorer.Func(ErrInvalidTimeRange, ErrEnd)
	errTimeRangeOrderFn = errorer.Func(ErrInvalidTimeRange, ErrOrder)
)

// TimeRange is the range of Times from Start to End, inclusive. It does not
// wrap past midnight.
//
// Comparisons use each Time's instant on 0001-01-01 with its offset applied,
// so "18:00:00-05:00" is before "20:00:00-05:00".
type TimeRange struct {
	Start Time `json:"start"`
	End   Time `json:"end"`
}

// Between reports whether t is within the range, inclusive of Start and End.
func (this TimeRange) Between(t Time) bool {
	return !this.Before(t) && !this.After(t)
}

// Before reports whether t is strictly before Start.
func (this TimeRange) Before(t Time) bool {
	return t.Time().Before(this.Start.Time())
}

// After reports whether t is strictly after End.
func (this TimeRange) After(t Time) bool {
	return t.Time().After(this.End.Time())
}

// Overlaps reports whether other shares at least one instant with the range.
// Bounds are inclusive, so ranges that only touch at an endpoint overlap.
func (this TimeRange) Overlaps(other TimeRange) bool {
	return !this.Before(other.End) && !this.After(other.Start)
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidTimeRange and the field's sentinel (ErrStart,
// ErrEnd, or ErrOrder). A Time is always valid, so in practice only ErrOrder
// can occur.
func (this TimeRange) Validate() error {
	var errs []error

	if err := this.Start.Validate(); err != nil {
		errs = append(errs, errTimeRangeStartFn("%w", err))
	}

	if err := this.End.Validate(); err != nil {
		errs = append(errs, errTimeRangeEndFn("%w", err))
	}

	if this.Start.Time().After(this.End.Time()) {
		errs = append(errs, errTimeRangeOrderFn("%s > %s", this.Start, this.End))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this TimeRange) MarshalJSON() ([]byte, error) {
	type alias TimeRange
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// TimeRange. A JSON null leaves the TimeRange unchanged.
func (this *TimeRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias TimeRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := TimeRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewTimeRange returns a validated TimeRange from start to end.
func NewTimeRange(start, end Time) (TimeRange, error) {
	r := TimeRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return TimeRange{}, err
	}

	return r, nil
}

// MustNewTimeRange is like NewTimeRange but panics if the TimeRange is invalid.
func MustNewTimeRange(start, end Time) TimeRange {
	r, err := NewTimeRange(start, end)

	if err != nil {
		panic(err)
	}

	return r
}
