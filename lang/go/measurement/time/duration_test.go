package time_test

import (
	"encoding/json"
	"math"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseDuration(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.Duration
		output   string
		err      error
	}{
		{
			name:     "every date and time component",
			input:    "P1Y2M3DT4H5M6S",
			expected: time.Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6},
			output:   "P1Y2M3DT4H5M6S",
		},
		{
			name:     "years only",
			input:    "P1Y",
			expected: time.Duration{Years: 1},
			output:   "P1Y",
		},
		{
			name:     "months and days",
			input:    "P2M3D",
			expected: time.Duration{Months: 2, Days: 3},
			output:   "P2M3D",
		},
		{
			name:     "days only",
			input:    "P3D",
			expected: time.Duration{Days: 3},
			output:   "P3D",
		},
		{
			name:     "time only",
			input:    "PT30M",
			expected: time.Duration{Minutes: 30},
			output:   "PT30M",
		},
		{
			name:     "minutes and seconds",
			input:    "PT5M6S",
			expected: time.Duration{Minutes: 5, Seconds: 6},
			output:   "PT5M6S",
		},
		{
			name:     "skipping across T",
			input:    "P1YT1S",
			expected: time.Duration{Years: 1, Seconds: 1},
			output:   "P1YT1S",
		},
		{
			name:     "explicit zeros fill a skipped unit",
			input:    "P1Y0M1D",
			expected: time.Duration{Years: 1, Days: 1},
			output:   "P1Y0M1D",
		},
		{
			name:     "leading zero components are dropped on output",
			input:    "P0Y0M1D",
			expected: time.Duration{Days: 1},
			output:   "P1D",
		},
		{
			name:     "weeks",
			input:    "P2W",
			expected: time.Duration{Weeks: 2},
			output:   "P2W",
		},
		{
			name:     "zero days",
			input:    "P0D",
			expected: time.ZeroDuration,
			output:   "PT0S",
		},
		{
			name:     "zero seconds",
			input:    "PT0S",
			expected: time.ZeroDuration,
			output:   "PT0S",
		},
		{
			name:     "large values are not normalized",
			input:    "PT90M",
			expected: time.Duration{Minutes: 90},
			output:   "PT90M",
		},
		{name: "empty", input: "", err: time.ErrInvalidDuration},
		{name: "missing P", input: "1Y", err: time.ErrInvalidDuration},
		{name: "P alone", input: "P", err: time.ErrInvalidDuration},
		{name: "PT alone", input: "PT", err: time.ErrInvalidDuration},
		{name: "T with nothing after", input: "P1DT", err: time.ErrInvalidDuration},
		{name: "skipped date unit", input: "P1Y1D", err: time.ErrInvalidDuration},
		{name: "skipped time unit", input: "PT1H1S", err: time.ErrInvalidDuration},
		{name: "out of order", input: "P1D1Y", err: time.ErrInvalidDuration},
		{name: "repeated unit", input: "P1D1D", err: time.ErrInvalidDuration},
		{name: "time unit before T", input: "P1H", err: time.ErrInvalidDuration},
		{name: "date unit after T", input: "PT1D", err: time.ErrInvalidDuration},
		{name: "weeks with years", input: "P1Y2W", err: time.ErrInvalidDuration},
		{name: "weeks with time", input: "P2WT1H", err: time.ErrInvalidDuration},
		{name: "repeated weeks", input: "P1W2W", err: time.ErrInvalidDuration},
		{name: "fraction", input: "PT1.5S", err: time.ErrInvalidDuration},
		{name: "negative", input: "P-1D", err: time.ErrInvalidDuration},
		{name: "missing number", input: "PD", err: time.ErrInvalidDuration},
		{name: "missing designator", input: "P1", err: time.ErrInvalidDuration},
		{name: "lowercase", input: "p1d", err: time.ErrInvalidDuration},
		{name: "second T", input: "PT1HT1M", err: time.ErrInvalidDuration},
		{name: "overflow", input: "P99999999999999999999D", err: time.ErrInvalidDuration},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseDuration(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.True(tt, actual.IsZero())
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
			assert.Equal(tt, tc.output, actual.String())
			assert.NoError(tt, actual.Validate())
		})
	}
}

func Test_Duration_String(t *testing.T) {
	tcs := []struct {
		name     string
		input    time.Duration
		expected string
	}{
		{name: "zero", input: time.ZeroDuration, expected: "PT0S"},
		{name: "fills skipped date unit", input: time.Duration{Years: 1, Days: 1}, expected: "P1Y0M1D"},
		{name: "fills skipped time unit", input: time.Duration{Hours: 1, Seconds: 1}, expected: "PT1H0M1S"},
		{name: "fills both parts", input: time.Duration{Years: 1, Days: 2, Hours: 3, Seconds: 4}, expected: "P1Y0M2DT3H0M4S"},
		{name: "date and time without skipping", input: time.Duration{Days: 1, Hours: 12}, expected: "P1DT12H"},
		{name: "weeks", input: time.Duration{Weeks: 3}, expected: "P3W"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, tc.input.String())

			parsed, err := time.ParseDuration(tc.input.String())
			assert.NoError(tt, err, "output round-trips")
			assert.Equal(tt, tc.input, parsed)
		})
	}
}

func Test_Duration_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.Duration
		errs  []error
	}{
		{name: "zero", input: time.ZeroDuration},
		{name: "every date and time component", input: time.Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6}},
		{name: "weeks alone", input: time.Duration{Weeks: 2}},
		{name: "negative years", input: time.Duration{Years: -1}, errs: []error{time.ErrInvalidDuration, time.ErrYears}},
		{name: "negative months", input: time.Duration{Months: -1}, errs: []error{time.ErrInvalidDuration, time.ErrMonths}},
		{name: "negative weeks", input: time.Duration{Weeks: -1}, errs: []error{time.ErrInvalidDuration, time.ErrWeeks}},
		{name: "negative days", input: time.Duration{Days: -1}, errs: []error{time.ErrInvalidDuration, time.ErrDays}},
		{name: "negative hours", input: time.Duration{Hours: -1}, errs: []error{time.ErrInvalidDuration, time.ErrHours}},
		{name: "negative minutes", input: time.Duration{Minutes: -1}, errs: []error{time.ErrInvalidDuration, time.ErrMinutes}},
		{name: "negative seconds", input: time.Duration{Seconds: -1}, errs: []error{time.ErrInvalidDuration, time.ErrSeconds}},
		{name: "weeks with days", input: time.Duration{Weeks: 1, Days: 1}, errs: []error{time.ErrInvalidDuration, time.ErrWeeksCombined}},
		{name: "weeks with seconds", input: time.Duration{Weeks: 1, Seconds: 1}, errs: []error{time.ErrInvalidDuration, time.ErrWeeksCombined}},
		{
			name:  "every failure is reported",
			input: time.Duration{Years: -1, Weeks: 1, Seconds: -1},
			errs:  []error{time.ErrYears, time.ErrSeconds, time.ErrWeeksCombined},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

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

func Test_Duration_Time(t *testing.T) {
	tcs := []struct {
		name     string
		input    time.Duration
		expected gotime.Duration
		err      error
	}{
		{name: "zero", input: time.ZeroDuration, expected: 0},
		{name: "seconds", input: time.MustParseDuration("PT45S"), expected: 45 * gotime.Second},
		{name: "minutes", input: time.MustParseDuration("PT90M"), expected: 90 * gotime.Minute},
		{name: "hours, minutes, and seconds", input: time.MustParseDuration("PT4H5M6S"), expected: 4*gotime.Hour + 5*gotime.Minute + 6*gotime.Second},
		{name: "exact 24 hours", input: time.MustParseDuration("PT24H"), expected: 24 * gotime.Hour},
		{name: "weeks", input: time.MustParseDuration("P2W"), expected: 2 * 7 * 24 * gotime.Hour},
		{name: "largest seconds", input: time.Duration{Seconds: int(math.MaxInt64 / int64(gotime.Second))}, expected: gotime.Duration(math.MaxInt64/int64(gotime.Second)) * gotime.Second},
		{name: "years", input: time.MustParseDuration("P1Y"), err: time.ErrNotFixedLength},
		{name: "months", input: time.MustParseDuration("P1M"), err: time.ErrNotFixedLength},
		{name: "days", input: time.MustParseDuration("P1D"), err: time.ErrNotFixedLength},
		{name: "days with time", input: time.MustParseDuration("P1DT1H"), err: time.ErrNotFixedLength},
		{name: "invalid", input: time.Duration{Hours: -1}, err: time.ErrHours},
		{name: "weeks combined", input: time.Duration{Weeks: 1, Hours: 1}, err: time.ErrWeeksCombined},
		{name: "overflow in seconds", input: time.Duration{Seconds: int(math.MaxInt64/int64(gotime.Second)) + 1}, err: time.ErrOverflow},
		{name: "overflow in hours", input: time.Duration{Hours: 3_000_000}, err: time.ErrOverflow},
		{name: "overflow in weeks", input: time.Duration{Weeks: 20_000}, err: time.ErrOverflow},
		{name: "overflow from the sum", input: time.Duration{Hours: 2_562_047, Minutes: 60}, err: time.ErrOverflow},
		{name: "overflow without wrapping int64", input: time.Duration{Weeks: math.MaxInt64}, err: time.ErrOverflow},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := tc.input.Time()

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.Zero(tt, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_NewDuration(t *testing.T) {
	actual, err := time.NewDuration(1, 2, 3, 4, 5, 6)
	assert.NoError(t, err)
	assert.Equal(t, "P1Y2M3DT4H5M6S", actual.String())

	_, err = time.NewDuration(0, 0, -1, 0, 0, 0)
	assert.ErrorIs(t, err, time.ErrDays)

	assert.Panics(t, func() { time.MustNewDuration(-1, 0, 0, 0, 0, 0) })
}

func Test_NewWeeksDuration(t *testing.T) {
	actual, err := time.NewWeeksDuration(2)
	assert.NoError(t, err)
	assert.Equal(t, "P2W", actual.String())

	_, err = time.NewWeeksDuration(-1)
	assert.ErrorIs(t, err, time.ErrWeeks)

	assert.Panics(t, func() { time.MustNewWeeksDuration(-1) })
}

func Test_Duration_JSON(t *testing.T) {
	type wrapper struct {
		For time.Duration `json:"for"`
	}

	expected := wrapper{For: time.MustParseDuration("P1Y0M1DT2H")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"for":"P1Y0M1DT2H"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"for":"P1Y1D"}`), &actual), time.ErrInvalidDuration)

	_, err = json.Marshal(wrapper{For: time.Duration{Weeks: 1, Days: 1}})
	assert.ErrorIs(t, err, time.ErrWeeksCombined, "invalid durations are not encoded")
}

func Test_Duration_UnmarshalJSON(t *testing.T) {
	valid := time.MustParseDuration("PT30M")

	tcs := []struct {
		name     string
		input    string
		start    time.Duration
		expected time.Duration
		err      error
	}{
		{
			name:     "string",
			input:    `"PT30M"`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:     "invalid duration leaves value unchanged",
			input:    `"PT1H1S"`,
			start:    valid,
			expected: valid,
			err:      time.ErrInvalidDuration,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   time.ErrInvalidDuration,
		},
		{
			name:  "number",
			input: `123`,
			err:   time.ErrInvalidDuration,
		},
		{
			name:  "object",
			input: `{"years":1}`,
			err:   time.ErrInvalidDuration,
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

func Test_Duration_Text(t *testing.T) {
	m := map[time.Duration]string{
		time.MustParseDuration("P2W"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"P2W":"a"}`, string(data))

	var actual map[time.Duration]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
