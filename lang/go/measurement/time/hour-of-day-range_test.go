package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_HourOfDayRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.HourOfDayRange
		errs  []error
	}{
		{name: "valid", input: time.HourOfDayRange{Start: 9, End: 17}},
		{name: "start equals end", input: time.HourOfDayRange{Start: 9, End: 9}},
		{name: "whole day", input: time.HourOfDayRange{Start: 0, End: 23}},
		{name: "zero value is midnight to midnight", input: time.HourOfDayRange{}},
		{
			name:  "invalid start",
			input: time.HourOfDayRange{Start: -1, End: 17},
			errs:  []error{time.ErrInvalidHourOfDayRange, time.ErrStart, time.ErrInvalidHourOfDay},
		},
		{
			name:  "invalid end",
			input: time.HourOfDayRange{Start: 9, End: 24},
			errs:  []error{time.ErrInvalidHourOfDayRange, time.ErrEnd, time.ErrInvalidHourOfDay},
		},
		{
			name:  "missing end reads as midnight",
			input: time.HourOfDayRange{Start: 9},
			errs:  []error{time.ErrInvalidHourOfDayRange, time.ErrOrder},
		},
		{
			name:  "start after end",
			input: time.HourOfDayRange{Start: 17, End: 9},
			errs:  []error{time.ErrInvalidHourOfDayRange, time.ErrOrder},
		},
		{
			name:  "does not wrap past midnight",
			input: time.HourOfDayRange{Start: 22, End: 2},
			errs:  []error{time.ErrInvalidHourOfDayRange, time.ErrOrder},
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

func Test_HourOfDayRange_Operations(t *testing.T) {
	r := time.MustNewHourOfDayRange(9, 17)

	tcs := []struct {
		name    string
		input   time.HourOfDay
		before  bool
		between bool
		after   bool
	}{
		{name: "midnight", input: 0, before: true},
		{name: "before start", input: 8, before: true},
		{name: "at start", input: 9, between: true},
		{name: "inside", input: 12, between: true},
		{name: "at end", input: 17, between: true},
		{name: "after end", input: 18, after: true},
		{name: "last hour", input: 23, after: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.before, r.Before(tc.input), "Before")
			assert.Equal(tt, tc.between, r.Between(tc.input), "Between")
			assert.Equal(tt, tc.after, r.After(tc.input), "After")
		})
	}
}

func Test_HourOfDayRange_Overlaps(t *testing.T) {
	r := time.MustNewHourOfDayRange(9, 17)

	tcs := []struct {
		name     string
		input    time.HourOfDayRange
		expected bool
	}{
		{name: "entirely before", input: time.MustNewHourOfDayRange(0, 8)},
		{name: "touches start", input: time.MustNewHourOfDayRange(6, 9), expected: true},
		{name: "overlaps start", input: time.MustNewHourOfDayRange(6, 10), expected: true},
		{name: "inside", input: time.MustNewHourOfDayRange(12, 13), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: time.MustNewHourOfDayRange(0, 23), expected: true},
		{name: "overlaps end", input: time.MustNewHourOfDayRange(16, 20), expected: true},
		{name: "touches end", input: time.MustNewHourOfDayRange(17, 20), expected: true},
		{name: "entirely after", input: time.MustNewHourOfDayRange(18, 23)},
		{name: "single hour inside", input: time.MustNewHourOfDayRange(12, 12), expected: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_HourOfDayRange_SingleHour(t *testing.T) {
	r := time.MustNewHourOfDayRange(9, 9)

	assert.True(t, r.Between(9))
	assert.True(t, r.Before(8))
	assert.True(t, r.After(10))
}

func Test_NewHourOfDayRange(t *testing.T) {
	actual, err := time.NewHourOfDayRange(9, 17)
	assert.NoError(t, err)
	assert.Equal(t, time.HourOfDay(9), actual.Start)
	assert.Equal(t, time.HourOfDay(17), actual.End)

	_, err = time.NewHourOfDayRange(17, 9)
	assert.ErrorIs(t, err, time.ErrOrder)

	assert.Panics(t, func() { time.MustNewHourOfDayRange(22, 2) })
}

func Test_HourOfDayRange_JSON(t *testing.T) {
	expected := time.MustNewHourOfDayRange(9, 17)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":9,"end":17}`, string(data))

	var actual time.HourOfDayRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func Test_HourOfDayRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewHourOfDayRange(9, 17)

	tcs := []struct {
		name     string
		input    string
		start    time.HourOfDayRange
		expected time.HourOfDayRange
		err      error
	}{
		{name: "valid", input: `{"start":9,"end":17}`, expected: valid},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid},
		{name: "start after end leaves value unchanged", input: `{"start":17,"end":9}`, start: valid, expected: valid, err: time.ErrOrder},
		{name: "missing end reads as midnight", input: `{"start":9}`, err: time.ErrOrder},
		{name: "invalid hour", input: `{"start":9,"end":24}`, err: time.ErrInvalidHourOfDay},
		{name: "string hour", input: `{"start":"9","end":17}`, err: time.ErrInvalidHourOfDay},
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
