package time

// Source of truth: domain-model/measurement/time/month-range.md

import (
	"encoding/json"
	"slices"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidMonthRange is wrapped by every MonthRange validation error.
	// The field sentinel ErrMonths is shared with Duration.
	ErrInvalidMonthRange = errorer.New("invalid month range")
	ErrDuplicateMonth    = errorer.New("duplicate month")

	errMonthRangeMonthsFn    = errorer.Func(ErrInvalidMonthRange, ErrMonths)
	errMonthRangeDuplicateFn = errorer.Func(ErrInvalidMonthRange, ErrMonths, ErrDuplicateMonth)
)

// MonthRange is a set of months of the year, e.g. a business's season.
// Months are kept in ISO 8601 order (January first), each at most once.
type MonthRange struct {
	Months []Month `json:"months"`
}

// AddMonth adds m to Months, keeping ISO 8601 order. Adding a month that is
// already present does nothing. It returns an error wrapping
// ErrInvalidMonth if m is invalid.
func (this *MonthRange) AddMonth(m Month) error {
	if err := m.Validate(); err != nil {
		return err
	}

	if slices.Contains(this.Months, m) {
		return nil
	}

	i := slices.IndexFunc(this.Months, func(other Month) bool { return other > m })

	if i < 0 {
		i = len(this.Months)
	}

	this.Months = slices.Insert(this.Months, i, m)
	return nil
}

// AddRange adds every month from start to end, inclusive. If start comes
// after end, the range wraps past December: November to February adds
// November, December, January, and February. It returns an error wrapping
// ErrInvalidMonth, and adds nothing, if either month is invalid.
func (this *MonthRange) AddRange(start, end Month) error {
	if err := start.Validate(); err != nil {
		return err
	}

	if err := end.Validate(); err != nil {
		return err
	}

	for m := start; ; m = m%12 + 1 {
		this.AddMonth(m)

		if m == end {
			return nil
		}
	}
}

// Contains reports whether m is in Months. It returns an error wrapping
// ErrInvalidMonth if m is invalid.
func (this MonthRange) Contains(m Month) (bool, error) {
	if err := m.Validate(); err != nil {
		return false, err
	}

	return slices.Contains(this.Months, m), nil
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidMonthRange and ErrMonths. An invalid month
// also wraps ErrInvalidMonth, and a repeated month ErrDuplicateMonth.
func (this MonthRange) Validate() error {
	var errs []error
	var seen [13]bool

	for _, m := range this.Months {
		if err := m.Validate(); err != nil {
			errs = append(errs, errMonthRangeMonthsFn("%w", err))
			continue
		}

		if seen[m] {
			errs = append(errs, errMonthRangeDuplicateFn("%s", m))
		}

		seen[m] = true
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler. An empty MonthRange encodes its
// Months as [] rather than null.
func (this MonthRange) MarshalJSON() ([]byte, error) {
	type alias MonthRange

	if this.Months == nil {
		this.Months = []Month{}
	}

	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// MonthRange. A JSON null leaves the MonthRange unchanged.
func (this *MonthRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias MonthRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := MonthRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewMonthRange returns a MonthRange holding months, in ISO 8601 order with
// repeats removed. It returns an error wrapping ErrInvalidMonth if any month
// is invalid.
func NewMonthRange(months ...Month) (MonthRange, error) {
	var r MonthRange

	for _, m := range months {
		if err := r.AddMonth(m); err != nil {
			return MonthRange{}, err
		}
	}

	return r, nil
}

// MustNewMonthRange is like NewMonthRange but panics if any month is invalid.
func MustNewMonthRange(months ...Month) MonthRange {
	r, err := NewMonthRange(months...)

	if err != nil {
		panic(err)
	}

	return r
}
