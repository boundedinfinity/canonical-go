package time

// Source of truth: domain-model/measurement/time/operating-hours.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	operatingHoursDaysValidation = validation.Slice[DayOfWeek]("days").Min(1)

	// ErrInvalidOperatingHours is wrapped by every OperatingHours validation
	// error. The field sentinels ErrDays and ErrHours are shared with Duration.
	ErrInvalidOperatingHours = errorer.New("invalid operating hours")
	ErrTimeZone              = errorer.New("time zone")
	ErrClock                 = errorer.New("clock")

	errOperatingHoursDaysFn     = errorer.Func(ErrInvalidOperatingHours, ErrDays)
	errOperatingHoursHoursFn    = errorer.Func(ErrInvalidOperatingHours, ErrHours)
	errOperatingHoursTimeZoneFn = errorer.Func(ErrInvalidOperatingHours, ErrTimeZone)
	errOperatingHoursClockFn    = errorer.Func(ErrInvalidOperatingHours, ErrClock)
)

// OperatingHours is the hours a business or service is open on a set of days
// of the week, e.g. Monday through Friday, 9 to 16 in America/New_York.
//
// Hours are local clock hours in TimeZone, and both ends are open hours:
// 9 to 16 is open from 9:00 until 17:00. Hours do not wrap past midnight.
type OperatingHours struct {
	Days     DayOfWeekRange               `json:"days"`
	Hours    HourOfDayRange               `json:"hours"`
	TimeZone TimeZone                     `json:"time-zone"`
	Clock    optioner.Option[ClockFormat] `json:"clock,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidOperatingHours and the field's sentinel
// (ErrDays, ErrHours, etc.), plus the nested model's own errors.
func (this OperatingHours) Validate() error {
	var errs []error

	if err := this.Days.Validate(); err != nil {
		errs = append(errs, errOperatingHoursDaysFn("%w", err))
	}

	if err := operatingHoursDaysValidation.Validate(this.Days.Days); err != nil {
		errs = append(errs, errOperatingHoursDaysFn("%w", err))
	}

	if err := this.Hours.Validate(); err != nil {
		errs = append(errs, errOperatingHoursHoursFn("%w", err))
	}

	if err := this.TimeZone.Validate(); err != nil {
		errs = append(errs, errOperatingHoursTimeZoneFn("%w", err))
	}

	if this.Clock.Defined() {
		if err := this.Clock.Get().Validate(); err != nil {
			errs = append(errs, errOperatingHoursClockFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this OperatingHours) MarshalJSON() ([]byte, error) {
	type alias OperatingHours
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// OperatingHours. A JSON null leaves the OperatingHours unchanged.
func (this *OperatingHours) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias OperatingHours
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	hours := OperatingHours(decoded)

	if err := hours.Validate(); err != nil {
		return err
	}

	*this = hours
	return nil
}

// NewOperatingHours returns a validated OperatingHours. A zero clock is
// treated as not present.
func NewOperatingHours(days DayOfWeekRange, hours HourOfDayRange, timeZone TimeZone, clock ClockFormat) (OperatingHours, error) {
	operatingHours := OperatingHours{
		Days:     days,
		Hours:    hours,
		TimeZone: timeZone,
		Clock:    optioner.OfZero(clock),
	}

	if err := operatingHours.Validate(); err != nil {
		return OperatingHours{}, err
	}

	return operatingHours, nil
}

// MustNewOperatingHours is like NewOperatingHours but panics if the
// OperatingHours is invalid.
func MustNewOperatingHours(days DayOfWeekRange, hours HourOfDayRange, timeZone TimeZone, clock ClockFormat) OperatingHours {
	operatingHours, err := NewOperatingHours(days, hours, timeZone, clock)

	if err != nil {
		panic(err)
	}

	return operatingHours
}
