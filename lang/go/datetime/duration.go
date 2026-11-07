package datetime

// Source of truth: domain-model/datetime/duration.md

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	durationYearsValidation   = validation.Integer[int]("years").Min(0)
	durationMonthsValidation  = validation.Integer[int]("months").Min(0)
	durationWeeksValidation   = validation.Integer[int]("weeks").Min(0)
	durationDaysValidation    = validation.Integer[int]("days").Min(0)
	durationHoursValidation   = validation.Integer[int]("hours").Min(0)
	durationMinutesValidation = validation.Integer[int]("minutes").Min(0)
	durationSecondsValidation = validation.Integer[int]("seconds").Min(0)

	// ErrInvalidDuration is wrapped by every Duration parse and validation error.
	ErrInvalidDuration = errorer.New("invalid duration")
	ErrYears           = errorer.New("years")
	ErrMonths          = errorer.New("months")
	ErrWeeks           = errorer.New("weeks")
	ErrDays            = errorer.New("days")
	ErrHours           = errorer.New("hours")
	ErrMinutes         = errorer.New("minutes")
	ErrSeconds         = errorer.New("seconds")
	ErrWeeksCombined   = errorer.New("weeks cannot be combined with other components")

	// ErrNotFixedLength is returned by Duration.Time when Years, Months, or
	// Days is non-zero.
	ErrNotFixedLength = errorer.New("years, months, and days have no fixed length")
	// ErrOverflow is returned by Duration.Time when the result does not fit
	// in a time.Duration.
	ErrOverflow = errorer.New("duration overflows time.Duration")

	errInvalidDurationFn = errorer.Func(ErrInvalidDuration)
	errYearsFn           = errorer.Func(ErrInvalidDuration, ErrYears)
	errMonthsFn          = errorer.Func(ErrInvalidDuration, ErrMonths)
	errWeeksFn           = errorer.Func(ErrInvalidDuration, ErrWeeks)
	errDaysFn            = errorer.Func(ErrInvalidDuration, ErrDays)
	errHoursFn           = errorer.Func(ErrInvalidDuration, ErrHours)
	errMinutesFn         = errorer.Func(ErrInvalidDuration, ErrMinutes)
	errSecondsFn         = errorer.Func(ErrInvalidDuration, ErrSeconds)
	errWeeksCombinedFn   = errorer.Func(ErrInvalidDuration, ErrWeeksCombined)
	errNotFixedLengthFn  = errorer.Func(ErrNotFixedLength)
	errOverflowFn        = errorer.Func(ErrOverflow)
)

// ZeroDuration is the zero value, a duration of no length ("PT0S"). It is a
// valid Duration.
var ZeroDuration = Duration{}

// Duration is a length of time as calendar and clock components, following
// the RFC 3339 duration format (Appendix A).
//
// Components are kept as written and never converted to one another, since
// years, months, and days have no fixed length: "P1M" and "P30D" differ.
// It serializes as a string, e.g. "P1Y2M3DT4H5M6S".
type Duration struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// String returns the RFC 3339 duration representation, using the fewest
// components the grammar allows and filling skipped units with zeros
// (Years 1 and Days 1 give "P1Y0M1D"). The zero Duration gives "PT0S".
func (this Duration) String() string {
	var b strings.Builder
	b.WriteString("P")

	if this.Weeks != 0 {
		fmt.Fprintf(&b, "%dW", this.Weeks)
	}

	date := [3]int{this.Years, this.Months, this.Days}
	clock := [3]int{this.Hours, this.Minutes, this.Seconds}

	if first, last, ok := durationSpan(date); ok {
		for i := first; i <= last; i++ {
			fmt.Fprintf(&b, "%d%c", date[i], "YMD"[i])
		}
	}

	if first, last, ok := durationSpan(clock); ok {
		b.WriteString("T")

		for i := first; i <= last; i++ {
			fmt.Fprintf(&b, "%d%c", clock[i], "HMS"[i])
		}
	}

	if b.Len() == 1 {
		return "PT0S"
	}

	return b.String()
}

// Time returns the Duration as a time.Duration. Weeks convert as 7 days of
// 24 hours; Hours, Minutes, and Seconds convert exactly.
//
// It returns the validation error if the Duration is invalid, an error
// wrapping ErrNotFixedLength if any of Years, Months, or Days is non-zero,
// and an error wrapping ErrOverflow if the result exceeds time.Duration.
func (this Duration) Time() (time.Duration, error) {
	if err := this.Validate(); err != nil {
		return 0, err
	}

	if this.Years != 0 || this.Months != 0 || this.Days != 0 {
		return 0, errNotFixedLengthFn("%s", this)
	}

	const maxSeconds = int64(math.MaxInt64 / int64(time.Second))

	components := []struct {
		value   int
		seconds int64
	}{
		{this.Weeks, 7 * 24 * 60 * 60},
		{this.Hours, 60 * 60},
		{this.Minutes, 60},
		{this.Seconds, 1},
	}

	var total int64

	for _, c := range components {
		if int64(c.value) > (maxSeconds-total)/c.seconds {
			return 0, errOverflowFn("%s", this)
		}

		total += int64(c.value) * c.seconds
	}

	return time.Duration(total) * time.Second, nil
}

// IsZero reports whether the Duration is the zero value.
func (this Duration) IsZero() bool {
	return this == ZeroDuration
}

// Validate checks every component and returns all failures joined together.
// Each failure wraps ErrInvalidDuration and the component's sentinel (ErrYears,
// etc.), or ErrWeeksCombined when Weeks is used with another component.
func (this Duration) Validate() error {
	var errs []error

	if err := durationYearsValidation.Validate(this.Years); err != nil {
		errs = append(errs, errYearsFn("%w", err))
	}

	if err := durationMonthsValidation.Validate(this.Months); err != nil {
		errs = append(errs, errMonthsFn("%w", err))
	}

	if err := durationWeeksValidation.Validate(this.Weeks); err != nil {
		errs = append(errs, errWeeksFn("%w", err))
	}

	if err := durationDaysValidation.Validate(this.Days); err != nil {
		errs = append(errs, errDaysFn("%w", err))
	}

	if err := durationHoursValidation.Validate(this.Hours); err != nil {
		errs = append(errs, errHoursFn("%w", err))
	}

	if err := durationMinutesValidation.Validate(this.Minutes); err != nil {
		errs = append(errs, errMinutesFn("%w", err))
	}

	if err := durationSecondsValidation.Validate(this.Seconds); err != nil {
		errs = append(errs, errSecondsFn("%w", err))
	}

	if this.Weeks != 0 && (Duration{Weeks: this.Weeks}) != this {
		errs = append(errs, errWeeksCombinedFn("%s", this))
	}

	return errorer.Join(errs...)
}

// MarshalText implements encoding.TextMarshaler, so the Duration serializes
// as a string. It returns the validation error for an invalid Duration.
func (this Duration) MarshalText() ([]byte, error) {
	if err := this.Validate(); err != nil {
		return nil, err
	}

	return []byte(this.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler and validates the input.
func (this *Duration) UnmarshalText(b []byte) error {
	d, err := ParseDuration(string(b))

	if err != nil {
		return err
	}

	*this = d
	return nil
}

// MarshalJSON implements json.Marshaler, encoding the Duration as a JSON
// string. It returns the validation error for an invalid Duration.
func (this Duration) MarshalJSON() ([]byte, error) {
	text, err := this.MarshalText()

	if err != nil {
		return nil, err
	}

	return json.Marshal(string(text))
}

// UnmarshalJSON implements json.Unmarshaler. It requires a JSON string
// containing a valid RFC 3339 duration. A JSON null leaves the Duration unchanged.
func (this *Duration) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return errInvalidDurationFn("%w", err)
	}

	return this.UnmarshalText([]byte(s))
}

// NewDuration returns a validated Duration from date and time components.
func NewDuration(years, months, days, hours, minutes, seconds int) (Duration, error) {
	d := Duration{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}

	if err := d.Validate(); err != nil {
		return ZeroDuration, err
	}

	return d, nil
}

// MustNewDuration is like NewDuration but panics if the Duration is invalid.
func MustNewDuration(years, months, days, hours, minutes, seconds int) Duration {
	d, err := NewDuration(years, months, days, hours, minutes, seconds)

	if err != nil {
		panic(err)
	}

	return d
}

// NewWeeksDuration returns a validated Duration of the given number of weeks.
func NewWeeksDuration(weeks int) (Duration, error) {
	d := Duration{Weeks: weeks}

	if err := d.Validate(); err != nil {
		return ZeroDuration, err
	}

	return d, nil
}

// MustNewWeeksDuration is like NewWeeksDuration but panics if weeks is negative.
func MustNewWeeksDuration(weeks int) Duration {
	d, err := NewWeeksDuration(weeks)

	if err != nil {
		panic(err)
	}

	return d
}

// ParseDuration returns the Duration represented by s.
// It returns an error wrapping ErrInvalidDuration if s is not a valid
// RFC 3339 duration.
func ParseDuration(s string) (Duration, error) {
	rest, ok := strings.CutPrefix(s, "P")

	if !ok {
		return ZeroDuration, errInvalidDurationFn("%q: must start with P", s)
	}

	date, clock, hasT := strings.Cut(rest, "T")

	if hasT && clock == "" {
		return ZeroDuration, errInvalidDurationFn("%q: T must be followed by a time component", s)
	}

	if date == "" && !hasT {
		return ZeroDuration, errInvalidDurationFn("%q: at least one component is required", s)
	}

	var d Duration

	if strings.HasSuffix(date, "W") {
		if hasT {
			return ZeroDuration, errInvalidDurationFn("%q: weeks cannot be combined with other components", s)
		}

		weeks, err := parseDurationPart(date, "W")

		if err != nil {
			return ZeroDuration, errInvalidDurationFn("%q: %w", s, err)
		}

		d.Weeks = weeks[0]
		return d, nil
	}

	values, err := parseDurationPart(date, "YMD")

	if err != nil {
		return ZeroDuration, errInvalidDurationFn("%q: %w", s, err)
	}

	d.Years, d.Months, d.Days = values[0], values[1], values[2]

	values, err = parseDurationPart(clock, "HMS")

	if err != nil {
		return ZeroDuration, errInvalidDurationFn("%q: %w", s, err)
	}

	d.Hours, d.Minutes, d.Seconds = values[0], values[1], values[2]

	return d, nil
}

// MustParseDuration is like ParseDuration but panics if s is not a valid
// RFC 3339 duration.
func MustParseDuration(s string) Duration {
	d, err := ParseDuration(s)

	if err != nil {
		panic(err)
	}

	return d
}

// parseDurationPart parses a run of <digits><designator> components, where
// each designator is from units, in order and without skipping a unit.
// The returned values are indexed by the designator's position in units.
func parseDurationPart(s, units string) ([3]int, error) {
	var values [3]int
	prev := -1

	for s != "" {
		end := strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' })

		if end == 0 {
			return values, fmt.Errorf("expected digits at %q", s)
		}

		if end < 0 {
			return values, fmt.Errorf("missing designator after %q", s)
		}

		unit := strings.IndexByte(units, s[end])

		if unit < 0 {
			return values, fmt.Errorf("unexpected designator %q", s[end])
		}

		if prev >= 0 && unit != prev+1 {
			return values, fmt.Errorf("%c cannot follow %c", units[unit], units[prev])
		}

		v, err := strconv.Atoi(s[:end])

		if err != nil {
			return values, err
		}

		values[unit] = v
		prev = unit
		s = s[end+1:]
	}

	return values, nil
}

// durationSpan returns the indexes of the first and last non-zero values, and
// false if every value is zero.
func durationSpan(values [3]int) (int, int, bool) {
	first, last := -1, -1

	for i, v := range values {
		if v != 0 {
			if first < 0 {
				first = i
			}

			last = i
		}
	}

	return first, last, first >= 0
}
