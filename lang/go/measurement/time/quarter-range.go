package time

// Source of truth: domain-model/measurement/time/quarter-range.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidQuarterRange is wrapped by every QuarterRange validation error.
	// The field sentinels (ErrStart, ErrEnd, ErrOrder) are shared with
	// DateTimeRange.
	ErrInvalidQuarterRange = errorer.New("invalid quarter range")

	errQuarterRangeStartFn = errorer.Func(ErrInvalidQuarterRange, ErrStart)
	errQuarterRangeEndFn   = errorer.Func(ErrInvalidQuarterRange, ErrEnd)
	errQuarterRangeOrderFn = errorer.Func(ErrInvalidQuarterRange, ErrOrder)
)

// QuarterRange is the range of Quarters from Start to End, inclusive. It does
// not wrap past Q4.
type QuarterRange struct {
	Start Quarter `json:"start"`
	End   Quarter `json:"end"`
}

// Between reports whether q is within the range, inclusive of Start and End.
func (this QuarterRange) Between(q Quarter) bool {
	return !this.Before(q) && !this.After(q)
}

// Before reports whether q is strictly before Start.
func (this QuarterRange) Before(q Quarter) bool {
	return q < this.Start
}

// After reports whether q is strictly after End.
func (this QuarterRange) After(q Quarter) bool {
	return q > this.End
}

// Overlaps reports whether other shares at least one Quarter with the range.
// Bounds are inclusive, so ranges that only touch at an endpoint overlap.
func (this QuarterRange) Overlaps(other QuarterRange) bool {
	return !this.Before(other.End) && !this.After(other.Start)
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidQuarterRange and the field's sentinel
// (ErrStart, ErrEnd, or ErrOrder). A Start or End failure also wraps
// ErrInvalidQuarter.
func (this QuarterRange) Validate() error {
	var errs []error

	if err := this.Start.Validate(); err != nil {
		errs = append(errs, errQuarterRangeStartFn("%w", err))
	}

	if err := this.End.Validate(); err != nil {
		errs = append(errs, errQuarterRangeEndFn("%w", err))
	}

	if this.Start > this.End {
		errs = append(errs, errQuarterRangeOrderFn("%s > %s", this.Start, this.End))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this QuarterRange) MarshalJSON() ([]byte, error) {
	type alias QuarterRange
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// QuarterRange. A JSON null leaves the QuarterRange unchanged.
func (this *QuarterRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias QuarterRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := QuarterRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewQuarterRange returns a validated QuarterRange from start to end.
func NewQuarterRange(start, end Quarter) (QuarterRange, error) {
	r := QuarterRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return QuarterRange{}, err
	}

	return r, nil
}

// MustNewQuarterRange is like NewQuarterRange but panics if the QuarterRange
// is invalid.
func MustNewQuarterRange(start, end Quarter) QuarterRange {
	r, err := NewQuarterRange(start, end)

	if err != nil {
		panic(err)
	}

	return r
}
