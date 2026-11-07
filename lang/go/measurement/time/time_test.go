package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseTime(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "utc",
			input:    "14:30:00Z",
			expected: "14:30:00Z",
		},
		{
			name:     "offset is preserved",
			input:    "09:30:00-05:00",
			expected: "09:30:00-05:00",
		},
		{
			name:     "fractional seconds",
			input:    "09:30:00.123-05:00",
			expected: "09:30:00.123-05:00",
		},
		{
			name:     "trailing zeros in fractional seconds are dropped",
			input:    "14:30:00.500000Z",
			expected: "14:30:00.5Z",
		},
		{
			name:     "midnight is valid",
			input:    "00:00:00Z",
			expected: "00:00:00Z",
		},
		{
			name:     "end of day",
			input:    "23:59:59.999999999+14:00",
			expected: "23:59:59.999999999+14:00",
		},
		{
			name:  "empty",
			input: "",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "missing offset",
			input: "14:30:00",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "missing seconds",
			input: "14:30Z",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "out of range hour",
			input: "24:00:00Z",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "out of range minute",
			input: "14:60:00Z",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "with date",
			input: "2026-11-07T14:30:00Z",
			err:   time.ErrInvalidTime,
		},
		{
			name:  "not a time",
			input: "not-a-time",
			err:   time.ErrInvalidTime,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseTime(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.True(tt, actual.IsZero())
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual.String())
			assert.NoError(tt, actual.Validate())
		})
	}
}

func Test_NewTime(t *testing.T) {
	loc := gotime.FixedZone("", -5*60*60)

	actual := time.NewTime(gotime.Date(2026, 11, 7, 9, 30, 0, 0, loc))
	assert.Equal(t, "09:30:00-05:00", actual.String())

	year, month, day := actual.Time().Date()
	assert.Equal(t, 1, year)
	assert.Equal(t, gotime.January, month)
	assert.Equal(t, 1, day)
}

func Test_CurrentTime(t *testing.T) {
	assert.NoError(t, time.CurrentTime().Validate())
}

func Test_ZeroTime(t *testing.T) {
	assert.Equal(t, "00:00:00Z", time.ZeroTime.String())
	assert.NoError(t, time.ZeroTime.Validate())
	assert.Equal(t, time.ZeroTime, time.MustParseTime("00:00:00Z"))
	assert.True(t, time.MustParseTime("00:00:00Z").IsZero())
	assert.False(t, time.MustParseTime("00:00:01Z").IsZero())
}

func Test_Time_Equal(t *testing.T) {
	tcs := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{name: "identical", a: "14:30:00Z", b: "14:30:00Z", expected: true},
		{name: "different offsets", a: "14:30:00Z", b: "09:30:00-05:00", expected: true},
		{name: "wraps past midnight", a: "23:00:00-05:00", b: "04:00:00Z", expected: true},
		{name: "wraps before midnight", a: "01:00:00+05:00", b: "20:00:00Z", expected: true},
		{name: "different second", a: "14:30:00Z", b: "14:30:01Z"},
		{name: "different nanosecond", a: "14:30:00Z", b: "14:30:00.000000001Z"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			a, b := time.MustParseTime(tc.a), time.MustParseTime(tc.b)

			assert.Equal(tt, tc.expected, a.Equal(b))
			assert.Equal(tt, tc.expected, b.Equal(a), "symmetric")
		})
	}
}

func Test_Time_JSON(t *testing.T) {
	type wrapper struct {
		At time.Time `json:"at"`
	}

	expected := wrapper{At: time.MustParseTime("09:30:00.123-05:00")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"at":"09:30:00.123-05:00"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.True(t, expected.At.Equal(actual.At))
	assert.Equal(t, expected.At.String(), actual.At.String())

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"at":"nope"}`), &actual), time.ErrInvalidTime)
}

func Test_Time_UnmarshalJSON(t *testing.T) {
	valid := time.MustParseTime("14:30:00Z")

	tcs := []struct {
		name     string
		input    string
		start    time.Time
		expected time.Time
		err      error
	}{
		{
			name:     "string",
			input:    `"14:30:00Z"`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:  "invalid time string",
			input: `"nope"`,
			err:   time.ErrInvalidTime,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   time.ErrInvalidTime,
		},
		{
			name:  "number",
			input: `123`,
			err:   time.ErrInvalidTime,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual := tc.start
			err := actual.UnmarshalJSON([]byte(tc.input))

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				return
			}

			assert.NoError(tt, err)
			assert.True(tt, tc.expected.Equal(actual))
		})
	}
}

func Test_Time_Text(t *testing.T) {
	m := map[time.Time]string{
		time.MustParseTime("14:30:00Z"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"14:30:00Z":"a"}`, string(data))

	var actual map[time.Time]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Len(t, actual, 1)

	for k, v := range actual {
		assert.Equal(t, "14:30:00Z", k.String())
		assert.Equal(t, "a", v)
	}
}
