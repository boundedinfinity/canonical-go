package datetime_test

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/boundedinfinity/canonical-go/lang/go/datetime"
	"github.com/stretchr/testify/assert"
)

func Test_ParseDuration(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected datetime.Duration
		output   string
		err      error
	}{
		{
			name:     "every date and time component",
			input:    "P1Y2M3DT4H5M6S",
			expected: datetime.Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6},
			output:   "P1Y2M3DT4H5M6S",
		},
		{
			name:     "years only",
			input:    "P1Y",
			expected: datetime.Duration{Years: 1},
			output:   "P1Y",
		},
		{
			name:     "months and days",
			input:    "P2M3D",
			expected: datetime.Duration{Months: 2, Days: 3},
			output:   "P2M3D",
		},
		{
			name:     "days only",
			input:    "P3D",
			expected: datetime.Duration{Days: 3},
			output:   "P3D",
		},
		{
			name:     "time only",
			input:    "PT30M",
			expected: datetime.Duration{Minutes: 30},
			output:   "PT30M",
		},
		{
			name:     "minutes and seconds",
			input:    "PT5M6S",
			expected: datetime.Duration{Minutes: 5, Seconds: 6},
			output:   "PT5M6S",
		},
		{
			name:     "skipping across T",
			input:    "P1YT1S",
			expected: datetime.Duration{Years: 1, Seconds: 1},
			output:   "P1YT1S",
		},
		{
			name:     "explicit zeros fill a skipped unit",
			input:    "P1Y0M1D",
			expected: datetime.Duration{Years: 1, Days: 1},
			output:   "P1Y0M1D",
		},
		{
			name:     "leading zero components are dropped on output",
			input:    "P0Y0M1D",
			expected: datetime.Duration{Days: 1},
			output:   "P1D",
		},
		{
			name:     "weeks",
			input:    "P2W",
			expected: datetime.Duration{Weeks: 2},
			output:   "P2W",
		},
		{
			name:     "zero days",
			input:    "P0D",
			expected: datetime.ZeroDuration,
			output:   "PT0S",
		},
		{
			name:     "zero seconds",
			input:    "PT0S",
			expected: datetime.ZeroDuration,
			output:   "PT0S",
		},
		{
			name:     "large values are not normalized",
			input:    "PT90M",
			expected: datetime.Duration{Minutes: 90},
			output:   "PT90M",
		},
		{name: "empty", input: "", err: datetime.ErrInvalidDuration},
		{name: "missing P", input: "1Y", err: datetime.ErrInvalidDuration},
		{name: "P alone", input: "P", err: datetime.ErrInvalidDuration},
		{name: "PT alone", input: "PT", err: datetime.ErrInvalidDuration},
		{name: "T with nothing after", input: "P1DT", err: datetime.ErrInvalidDuration},
		{name: "skipped date unit", input: "P1Y1D", err: datetime.ErrInvalidDuration},
		{name: "skipped time unit", input: "PT1H1S", err: datetime.ErrInvalidDuration},
		{name: "out of order", input: "P1D1Y", err: datetime.ErrInvalidDuration},
		{name: "repeated unit", input: "P1D1D", err: datetime.ErrInvalidDuration},
		{name: "time unit before T", input: "P1H", err: datetime.ErrInvalidDuration},
		{name: "date unit after T", input: "PT1D", err: datetime.ErrInvalidDuration},
		{name: "weeks with years", input: "P1Y2W", err: datetime.ErrInvalidDuration},
		{name: "weeks with time", input: "P2WT1H", err: datetime.ErrInvalidDuration},
		{name: "repeated weeks", input: "P1W2W", err: datetime.ErrInvalidDuration},
		{name: "fraction", input: "PT1.5S", err: datetime.ErrInvalidDuration},
		{name: "negative", input: "P-1D", err: datetime.ErrInvalidDuration},
		{name: "missing number", input: "PD", err: datetime.ErrInvalidDuration},
		{name: "missing designator", input: "P1", err: datetime.ErrInvalidDuration},
		{name: "lowercase", input: "p1d", err: datetime.ErrInvalidDuration},
		{name: "second T", input: "PT1HT1M", err: datetime.ErrInvalidDuration},
		{name: "overflow", input: "P99999999999999999999D", err: datetime.ErrInvalidDuration},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := datetime.ParseDuration(tc.input)

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
		input    datetime.Duration
		expected string
	}{
		{name: "zero", input: datetime.ZeroDuration, expected: "PT0S"},
		{name: "fills skipped date unit", input: datetime.Duration{Years: 1, Days: 1}, expected: "P1Y0M1D"},
		{name: "fills skipped time unit", input: datetime.Duration{Hours: 1, Seconds: 1}, expected: "PT1H0M1S"},
		{name: "fills both parts", input: datetime.Duration{Years: 1, Days: 2, Hours: 3, Seconds: 4}, expected: "P1Y0M2DT3H0M4S"},
		{name: "date and time without skipping", input: datetime.Duration{Days: 1, Hours: 12}, expected: "P1DT12H"},
		{name: "weeks", input: datetime.Duration{Weeks: 3}, expected: "P3W"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, tc.input.String())

			parsed, err := datetime.ParseDuration(tc.input.String())
			assert.NoError(tt, err, "output round-trips")
			assert.Equal(tt, tc.input, parsed)
		})
	}
}

func Test_Duration_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input datetime.Duration
		errs  []error
	}{
		{name: "zero", input: datetime.ZeroDuration},
		{name: "every date and time component", input: datetime.Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6}},
		{name: "weeks alone", input: datetime.Duration{Weeks: 2}},
		{name: "negative years", input: datetime.Duration{Years: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrYears}},
		{name: "negative months", input: datetime.Duration{Months: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrMonths}},
		{name: "negative weeks", input: datetime.Duration{Weeks: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrWeeks}},
		{name: "negative days", input: datetime.Duration{Days: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrDays}},
		{name: "negative hours", input: datetime.Duration{Hours: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrHours}},
		{name: "negative minutes", input: datetime.Duration{Minutes: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrMinutes}},
		{name: "negative seconds", input: datetime.Duration{Seconds: -1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrSeconds}},
		{name: "weeks with days", input: datetime.Duration{Weeks: 1, Days: 1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrWeeksCombined}},
		{name: "weeks with seconds", input: datetime.Duration{Weeks: 1, Seconds: 1}, errs: []error{datetime.ErrInvalidDuration, datetime.ErrWeeksCombined}},
		{
			name:  "every failure is reported",
			input: datetime.Duration{Years: -1, Weeks: 1, Seconds: -1},
			errs:  []error{datetime.ErrYears, datetime.ErrSeconds, datetime.ErrWeeksCombined},
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
		input    datetime.Duration
		expected time.Duration
		err      error
	}{
		{name: "zero", input: datetime.ZeroDuration, expected: 0},
		{name: "seconds", input: datetime.MustParseDuration("PT45S"), expected: 45 * time.Second},
		{name: "minutes", input: datetime.MustParseDuration("PT90M"), expected: 90 * time.Minute},
		{name: "hours, minutes, and seconds", input: datetime.MustParseDuration("PT4H5M6S"), expected: 4*time.Hour + 5*time.Minute + 6*time.Second},
		{name: "exact 24 hours", input: datetime.MustParseDuration("PT24H"), expected: 24 * time.Hour},
		{name: "weeks", input: datetime.MustParseDuration("P2W"), expected: 2 * 7 * 24 * time.Hour},
		{name: "largest seconds", input: datetime.Duration{Seconds: int(math.MaxInt64 / int64(time.Second))}, expected: time.Duration(math.MaxInt64/int64(time.Second)) * time.Second},
		{name: "years", input: datetime.MustParseDuration("P1Y"), err: datetime.ErrNotFixedLength},
		{name: "months", input: datetime.MustParseDuration("P1M"), err: datetime.ErrNotFixedLength},
		{name: "days", input: datetime.MustParseDuration("P1D"), err: datetime.ErrNotFixedLength},
		{name: "days with time", input: datetime.MustParseDuration("P1DT1H"), err: datetime.ErrNotFixedLength},
		{name: "invalid", input: datetime.Duration{Hours: -1}, err: datetime.ErrHours},
		{name: "weeks combined", input: datetime.Duration{Weeks: 1, Hours: 1}, err: datetime.ErrWeeksCombined},
		{name: "overflow in seconds", input: datetime.Duration{Seconds: int(math.MaxInt64/int64(time.Second)) + 1}, err: datetime.ErrOverflow},
		{name: "overflow in hours", input: datetime.Duration{Hours: 3_000_000}, err: datetime.ErrOverflow},
		{name: "overflow in weeks", input: datetime.Duration{Weeks: 20_000}, err: datetime.ErrOverflow},
		{name: "overflow from the sum", input: datetime.Duration{Hours: 2_562_047, Minutes: 60}, err: datetime.ErrOverflow},
		{name: "overflow without wrapping int64", input: datetime.Duration{Weeks: math.MaxInt64}, err: datetime.ErrOverflow},
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
	actual, err := datetime.NewDuration(1, 2, 3, 4, 5, 6)
	assert.NoError(t, err)
	assert.Equal(t, "P1Y2M3DT4H5M6S", actual.String())

	_, err = datetime.NewDuration(0, 0, -1, 0, 0, 0)
	assert.ErrorIs(t, err, datetime.ErrDays)

	assert.Panics(t, func() { datetime.MustNewDuration(-1, 0, 0, 0, 0, 0) })
}

func Test_NewWeeksDuration(t *testing.T) {
	actual, err := datetime.NewWeeksDuration(2)
	assert.NoError(t, err)
	assert.Equal(t, "P2W", actual.String())

	_, err = datetime.NewWeeksDuration(-1)
	assert.ErrorIs(t, err, datetime.ErrWeeks)

	assert.Panics(t, func() { datetime.MustNewWeeksDuration(-1) })
}

func Test_Duration_JSON(t *testing.T) {
	type wrapper struct {
		For datetime.Duration `json:"for"`
	}

	expected := wrapper{For: datetime.MustParseDuration("P1Y0M1DT2H")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"for":"P1Y0M1DT2H"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"for":"P1Y1D"}`), &actual), datetime.ErrInvalidDuration)

	_, err = json.Marshal(wrapper{For: datetime.Duration{Weeks: 1, Days: 1}})
	assert.ErrorIs(t, err, datetime.ErrWeeksCombined, "invalid durations are not encoded")
}

func Test_Duration_UnmarshalJSON(t *testing.T) {
	valid := datetime.MustParseDuration("PT30M")

	tcs := []struct {
		name     string
		input    string
		start    datetime.Duration
		expected datetime.Duration
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
			err:      datetime.ErrInvalidDuration,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   datetime.ErrInvalidDuration,
		},
		{
			name:  "number",
			input: `123`,
			err:   datetime.ErrInvalidDuration,
		},
		{
			name:  "object",
			input: `{"years":1}`,
			err:   datetime.ErrInvalidDuration,
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
	m := map[datetime.Duration]string{
		datetime.MustParseDuration("P2W"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"P2W":"a"}`, string(data))

	var actual map[datetime.Duration]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
