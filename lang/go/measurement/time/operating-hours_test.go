package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func weekdays() time.DayOfWeekRange {
	var days time.DayOfWeekRange
	days.AddRange(time.Monday, time.Friday)
	return days
}

func validOperatingHours() time.OperatingHours {
	return time.OperatingHours{
		Days:     weekdays(),
		Hours:    time.MustNewHourOfDayRange(9, 16),
		TimeZone: validTimeZone(),
	}
}

func Test_OperatingHours_Validate(t *testing.T) {
	tcs := []struct {
		name   string
		modify func(*time.OperatingHours)
		errs   []error
	}{
		{name: "valid", modify: func(*time.OperatingHours) {}},
		{name: "with clock", modify: func(h *time.OperatingHours) { h.Clock = optioner.Some(time.ClockFormat12Hour) }},
		{name: "single day", modify: func(h *time.OperatingHours) { h.Days = time.MustNewDayOfWeekRange(time.Saturday) }},
		{name: "single hour", modify: func(h *time.OperatingHours) { h.Hours = time.MustNewHourOfDayRange(12, 12) }},
		{name: "until midnight", modify: func(h *time.OperatingHours) { h.Hours = time.MustNewHourOfDayRange(9, 23) }},
		{name: "24 hours", modify: func(h *time.OperatingHours) { h.Hours = time.MustNewHourOfDayRange(0, 23) }},
		{
			name:   "no days",
			modify: func(h *time.OperatingHours) { h.Days = time.DayOfWeekRange{} },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrDays},
		},
		{
			name:   "invalid days",
			modify: func(h *time.OperatingHours) { h.Days.Days = append(h.Days.Days, time.Monday) },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrDays, time.ErrDuplicateDay},
		},
		{
			name:   "overnight does not wrap",
			modify: func(h *time.OperatingHours) { h.Hours = time.HourOfDayRange{Start: 18, End: 2} },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrHours, time.ErrInvalidHourOfDayRange, time.ErrOrder},
		},
		{
			name:   "invalid hour",
			modify: func(h *time.OperatingHours) { h.Hours = time.HourOfDayRange{Start: 9, End: 24} },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrHours, time.ErrInvalidHourOfDay},
		},
		{
			name:   "missing time zone",
			modify: func(h *time.OperatingHours) { h.TimeZone = time.TimeZone{} },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrTimeZone, time.ErrInvalidTimeZone},
		},
		{
			name:   "invalid clock",
			modify: func(h *time.OperatingHours) { h.Clock = optioner.Some(time.ZeroClockFormat) },
			errs:   []error{time.ErrInvalidOperatingHours, time.ErrClock, time.ErrInvalidClockFormat},
		},
		{
			name: "every failure is reported",
			modify: func(h *time.OperatingHours) {
				*h = time.OperatingHours{Hours: time.HourOfDayRange{Start: 9, End: 24}}
			},
			errs: []error{time.ErrDays, time.ErrHours, time.ErrTimeZone},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			hours := validOperatingHours()
			tc.modify(&hours)
			err := hours.Validate()

			if len(tc.errs) == 0 {
				assert.NoError(tt, err)
				return
			}

			for _, expected := range tc.errs {
				assert.ErrorIs(tt, err, expected)
			}
		})
	}
}

func Test_NewOperatingHours(t *testing.T) {
	tz, ok := time.TimeZones.Get("America/Chicago")
	assert.True(t, ok)

	nineToFive := time.MustNewHourOfDayRange(9, 16)

	actual, err := time.NewOperatingHours(weekdays(), nineToFive, tz, time.ZeroClockFormat)
	assert.NoError(t, err)
	assert.Equal(t, nineToFive, actual.Hours)
	assert.False(t, actual.Clock.Defined(), "zero clock is treated as not present")

	actual, err = time.NewOperatingHours(weekdays(), nineToFive, tz, time.ClockFormat24Hour)
	assert.NoError(t, err)
	assert.Equal(t, time.ClockFormat24Hour, actual.Clock.Get())

	_, err = time.NewOperatingHours(time.DayOfWeekRange{}, nineToFive, tz, time.ZeroClockFormat)
	assert.ErrorIs(t, err, time.ErrDays)

	assert.Panics(t, func() {
		time.MustNewOperatingHours(weekdays(), time.HourOfDayRange{Start: 18, End: 2}, tz, time.ZeroClockFormat)
	})
}

func Test_OperatingHours_JSON(t *testing.T) {
	expected := validOperatingHours()
	expected.Clock = optioner.Some(time.ClockFormat12Hour)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"days": {"days": ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]},
		"hours": {"start": 9, "end": 16},
		"time-zone": {
			"identifier": "America/New_York",
			"name": "Eastern Standard Time",
			"abbreviation": "EST",
			"utc-offset": "-05:00"
		},
		"clock": "12-hour"
	}`, string(data))

	var actual time.OperatingHours
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	data, err = json.Marshal(validOperatingHours())
	assert.NoError(t, err)
	assert.NotContains(t, string(data), "clock", "absent clock is omitted")
}

func Test_OperatingHours_UnmarshalJSON(t *testing.T) {
	valid := validOperatingHours()
	tz := `{"identifier":"America/New_York","name":"Eastern Standard Time","abbreviation":"EST","utc-offset":"-05:00"}`

	tcs := []struct {
		name     string
		input    string
		start    time.OperatingHours
		expected time.OperatingHours
		err      error
	}{
		{
			name:     "valid",
			input:    `{"days":{"days":["mon","tue","wed","thu","fri"]},"hours":{"start":9,"end":16},"time-zone":` + tz + `}`,
			expected: valid,
		},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid},
		{
			name:     "no days leaves value unchanged",
			input:    `{"days":{"days":[]},"hours":{"start":9,"end":16},"time-zone":` + tz + `}`,
			start:    valid,
			expected: valid,
			err:      time.ErrDays,
		},
		{
			name:  "missing time zone",
			input: `{"days":{"days":["Mon"]},"hours":{"start":9,"end":16}}`,
			err:   time.ErrTimeZone,
		},
		{
			name:  "overnight hours",
			input: `{"days":{"days":["Fri"]},"hours":{"start":18,"end":2},"time-zone":` + tz + `}`,
			err:   time.ErrOrder,
		},
		{
			name:  "invalid hour",
			input: `{"days":{"days":["Mon"]},"hours":{"start":9,"end":24},"time-zone":` + tz + `}`,
			err:   time.ErrInvalidHourOfDay,
		},
		{
			name:  "invalid clock",
			input: `{"days":{"days":["Mon"]},"hours":{"start":9,"end":16},"time-zone":` + tz + `,"clock":"24h"}`,
			err:   time.ErrInvalidClockFormat,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual := tc.start
			err := actual.UnmarshalJSON([]byte(tc.input))

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
			} else {
				assert.NoError(tt, err)
			}

			assert.Equal(tt, tc.expected, actual)
		})
	}
}
