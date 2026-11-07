package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseHourOfDay(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.HourOfDay
		err      error
	}{
		{name: "midnight", input: "0", expected: 0},
		{name: "morning", input: "7", expected: 7},
		{name: "leading zero", input: "07", expected: 7},
		{name: "afternoon", input: "14", expected: 14},
		{name: "last hour", input: "23", expected: 23},
		{name: "padded midnight", input: "00", expected: 0},
		{name: "empty", input: "", err: time.ErrInvalidHourOfDay},
		{name: "out of range", input: "24", err: time.ErrInvalidHourOfDay},
		{name: "negative", input: "-1", err: time.ErrInvalidHourOfDay},
		{name: "two leading zeros", input: "007", err: time.ErrInvalidHourOfDay},
		{name: "with minutes", input: "7:00", err: time.ErrInvalidHourOfDay},
		{name: "12-hour form", input: "7 AM", err: time.ErrInvalidHourOfDay},
		{name: "surrounding space", input: " 7", err: time.ErrInvalidHourOfDay},
		{name: "plus sign", input: "+7", err: time.ErrInvalidHourOfDay},
		{name: "fraction", input: "7.5", err: time.ErrInvalidHourOfDay},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseHourOfDay(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.True(tt, actual.IsZero())
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
			assert.NoError(tt, actual.Validate())
		})
	}
}

func Test_HourOfDay_Values(t *testing.T) {
	for n := range 24 {
		h := time.HourOfDay(n)

		assert.NoError(t, h.Validate())
		assert.Equal(t, n, h.Number())

		parsed, err := time.ParseHourOfDay(h.String())
		assert.NoError(t, err, "String round-trips")
		assert.Equal(t, h, parsed)
	}

	assert.Equal(t, "14", time.HourOfDay(14).String())
}

func Test_NewHourOfDay(t *testing.T) {
	assert.Equal(t, time.HourOfDay(14), time.NewHourOfDay(gotime.Date(2026, 11, 7, 14, 30, 0, 0, gotime.UTC)))
	assert.Equal(t, time.HourOfDay(0), time.NewHourOfDay(gotime.Date(2026, 11, 7, 0, 0, 0, 0, gotime.UTC)))

	// 23:30 in UTC-5 is already 04:30 in UTC; the hour in t's own location is
	// kept.
	loc := gotime.FixedZone("", -5*60*60)
	assert.Equal(t, time.HourOfDay(23), time.NewHourOfDay(gotime.Date(2026, 11, 7, 23, 30, 0, 0, loc)))
}

func Test_HourOfDay_Validate(t *testing.T) {
	assert.NoError(t, time.ZeroHourOfDay.Validate(), "midnight is valid")
	assert.True(t, time.ZeroHourOfDay.IsZero())
	assert.False(t, time.HourOfDay(1).IsZero())

	assert.ErrorIs(t, time.HourOfDay(24).Validate(), time.ErrInvalidHourOfDay)
	assert.ErrorIs(t, time.HourOfDay(-1).Validate(), time.ErrInvalidHourOfDay)
}

func Test_HourOfDay_JSON(t *testing.T) {
	type wrapper struct {
		At time.HourOfDay `json:"at"`
	}

	expected := wrapper{At: 14}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"at":14}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	data, err = json.Marshal(wrapper{})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"at":0}`, string(data), "midnight is encoded")

	_, err = json.Marshal(wrapper{At: 24})
	assert.ErrorIs(t, err, time.ErrInvalidHourOfDay, "an invalid hour is not encoded")
}

func Test_HourOfDay_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		start    time.HourOfDay
		expected time.HourOfDay
		err      error
	}{
		{name: "number", input: `9`, expected: 9},
		{name: "midnight", input: `0`, start: 5, expected: 0},
		{name: "null leaves value unchanged", input: `null`, start: 9, expected: 9},
		{name: "out of range leaves value unchanged", input: `24`, start: 9, expected: 9, err: time.ErrInvalidHourOfDay},
		{name: "negative", input: `-1`, err: time.ErrInvalidHourOfDay},
		{name: "fraction", input: `9.5`, err: time.ErrInvalidHourOfDay},
		{name: "string", input: `"9"`, err: time.ErrInvalidHourOfDay},
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

func Test_HourOfDay_Text(t *testing.T) {
	m := map[time.HourOfDay]string{
		0:  "midnight",
		14: "afternoon",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"0":"midnight","14":"afternoon"}`, string(data))

	var actual map[time.HourOfDay]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
