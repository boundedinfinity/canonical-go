package time

// Source of truth: domain-model/measurement/time/day-of-week-range.md

import (
	"encoding/json"
	"slices"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	// ErrInvalidDayOfWeekRange is wrapped by every DayOfWeekRange validation
	// error. The field sentinel ErrDays is shared with Duration.
	ErrInvalidDayOfWeekRange = errorer.New("invalid day of week range")
	ErrDuplicateDay          = errorer.New("duplicate day")

	errDayOfWeekRangeDaysFn      = errorer.Func(ErrInvalidDayOfWeekRange, ErrDays)
	errDayOfWeekRangeDuplicateFn = errorer.Func(ErrInvalidDayOfWeekRange, ErrDays, ErrDuplicateDay)
)

// DayOfWeekRange is a set of days of the week, e.g. a business's working days.
// Days are kept in ISO 8601 order (Monday first), each at most once.
type DayOfWeekRange struct {
	Days []DayOfWeek `json:"days"`
}

// AddDay adds d to Days, keeping ISO 8601 order. Adding a day that is
// already present does nothing. It returns an error wrapping
// ErrInvalidDayOfWeek if d is invalid.
func (this *DayOfWeekRange) AddDay(d DayOfWeek) error {
	if err := d.Validate(); err != nil {
		return err
	}

	if slices.Contains(this.Days, d) {
		return nil
	}

	i := slices.IndexFunc(this.Days, func(other DayOfWeek) bool { return other > d })

	if i < 0 {
		i = len(this.Days)
	}

	this.Days = slices.Insert(this.Days, i, d)
	return nil
}

// AddRange adds every day from start to end, inclusive. If start comes after
// end, the range wraps past Sunday: Friday to Monday adds Friday, Saturday,
// Sunday, and Monday. It returns an error wrapping ErrInvalidDayOfWeek, and
// adds nothing, if either day is invalid.
func (this *DayOfWeekRange) AddRange(start, end DayOfWeek) error {
	if err := start.Validate(); err != nil {
		return err
	}

	if err := end.Validate(); err != nil {
		return err
	}

	for d := start; ; d = d%7 + 1 {
		this.AddDay(d)

		if d == end {
			return nil
		}
	}
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidDayOfWeekRange and ErrDays. An invalid day
// also wraps ErrInvalidDayOfWeek, and a repeated day ErrDuplicateDay.
func (this DayOfWeekRange) Validate() error {
	var errs []error
	var seen [8]bool

	for _, d := range this.Days {
		if err := d.Validate(); err != nil {
			errs = append(errs, errDayOfWeekRangeDaysFn("%w", err))
			continue
		}

		if seen[d] {
			errs = append(errs, errDayOfWeekRangeDuplicateFn("%s", d))
		}

		seen[d] = true
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler. An empty DayOfWeekRange encodes its
// Days as [] rather than null.
func (this DayOfWeekRange) MarshalJSON() ([]byte, error) {
	type alias DayOfWeekRange

	if this.Days == nil {
		this.Days = []DayOfWeek{}
	}

	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// DayOfWeekRange. A JSON null leaves the DayOfWeekRange unchanged.
func (this *DayOfWeekRange) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias DayOfWeekRange
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r := DayOfWeekRange(decoded)

	if err := r.Validate(); err != nil {
		return err
	}

	*this = r
	return nil
}

// NewDayOfWeekRange returns a DayOfWeekRange holding days, in ISO 8601 order
// with repeats removed. It returns an error wrapping ErrInvalidDayOfWeek if
// any day is invalid.
func NewDayOfWeekRange(days ...DayOfWeek) (DayOfWeekRange, error) {
	var r DayOfWeekRange

	for _, d := range days {
		if err := r.AddDay(d); err != nil {
			return DayOfWeekRange{}, err
		}
	}

	return r, nil
}

// MustNewDayOfWeekRange is like NewDayOfWeekRange but panics if any day is
// invalid.
func MustNewDayOfWeekRange(days ...DayOfWeek) DayOfWeekRange {
	r, err := NewDayOfWeekRange(days...)

	if err != nil {
		panic(err)
	}

	return r
}
