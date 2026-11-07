package time

// Source of truth: domain-model/measurement/time/hour-of-day-range.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidHourOfDayRange is wrapped by every HourOfDayRange validation
	// error. The field sentinels (ErrStart, ErrEnd, ErrOrder) are shared with
	// DateTimeRange.
	ErrInvalidHourOfDayRange = errorer.New("invalid hour of day range")

	errHourOfDayRangeStartFn = errorer.Func(ErrInvalidHourOfDayRange, ErrStart)
	errHourOfDayRangeEndFn   = errorer.Func(ErrInvalidHourOfDayRange, ErrEnd)
	errHourOfDayRangeOrderFn = errorer.Func(ErrInvalidHourOfDayRange, ErrOrder)
)

// HourOfDayRange is the range of HourOfDays from Start to End, inclusive. It
// does not wrap past midnight.
type HourOfDayRange struct {
	Start HourOfDay `json:"start"`
	End   HourOfDay `json:"end"`
}

// Between reports whether h is within the range, inclusive of Start and End.
func (this HourOfDayRange) Between(h HourOfDay) bool {
	return !this.Before(h) && !this.After(h)
}

// Before reports whether h is strictly before Start.
func (this HourOfDayRange) Before(h HourOfDay) bool {
	return h < this.Start
}

// After reports whether h is strictly after End.
func (this HourOfDayRange) After(h HourOfDay) bool {
	return h > this.End
}

// Overlaps reports whether other shares at least one HourOfDay with the
// range. Bounds are inclusive, so ranges that only touch at an endpoint
// overlap.
func (this HourOfDayRange) Overlaps(other HourOfDayRange) bool {
	return !this.Before(other.End) && !this.After(other.Start)
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidHourOfDayRange and the field's sentinel
// (ErrStart, ErrEnd, or ErrOrder). A Start or End failure also wraps
// ErrInvalidHourOfDay.
func (this HourOfDayRange) Validate() error {
	var errs []error

	if err := this.Start.Validate(); err != nil {
		errs = append(errs, errHourOfDayRangeStartFn("%w", err))
	}

	if err := this.End.Validate(); err != nil {
		errs = append(errs, errHourOfDayRangeEndFn("%w", err))
	}

	if this.Start > this.End {
		errs = append(errs, errHourOfDayRangeOrderFn("%s > %s", this.Start, this.End))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this HourOfDayRange) MarshalJSON() ([]byte, error) {
	type alias HourOfDayRange
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// HourOfDayRange. A JSON null leaves the HourOfDayRange unchanged.
func (this *HourOfDayRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias HourOfDayRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := HourOfDayRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewHourOfDayRange returns a validated HourOfDayRange from start to end.
func NewHourOfDayRange(start, end HourOfDay) (HourOfDayRange, error) {
	r := HourOfDayRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return HourOfDayRange{}, err
	}

	return r, nil
}

// MustNewHourOfDayRange is like NewHourOfDayRange but panics if the
// HourOfDayRange is invalid.
func MustNewHourOfDayRange(start, end HourOfDay) HourOfDayRange {
	r, err := NewHourOfDayRange(start, end)

	if err != nil {
		panic(err)
	}

	return r
}
