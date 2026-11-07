package time

import "github.com/boundedinfinity/go-commoner/errorer"

var (
	ErrHourOfDay             = errorer.New("hour of day")
	ErrHourOfDayOutOfRange   = errorer.New("out of range")
	errHourOfDayOutOfRangeFn = errorer.Func(ErrHourOfDay, ErrHourOfDayOutOfRange)
)

type HourOfDay int

func (this HourOfDay) Validate() error {
	if this < 0 || this > 23 {
		return errHourOfDayOutOfRangeFn("%d is outside of the range 0 - 23", this)
	}

	return nil
}

type HourOfDayRange struct {
	Min HourOfDay `json:"min"`
	Max HourOfDay `json:"max"`
}
