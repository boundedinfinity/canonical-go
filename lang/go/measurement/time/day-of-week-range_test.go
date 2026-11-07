package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_DayOfWeekRange_AddDay(t *testing.T) {
	var r time.DayOfWeekRange

	assert.NoError(t, r.AddDay(time.Saturday))
	assert.NoError(t, r.AddDay(time.Monday))
	assert.NoError(t, r.AddDay(time.Wednesday))
	assert.NoError(t, r.AddDay(time.Monday), "adding a present day does nothing")
	assert.Equal(t, []time.DayOfWeek{time.Monday, time.Wednesday, time.Saturday}, r.Days, "ISO order, no repeats")

	assert.ErrorIs(t, r.AddDay(time.ZeroDayOfWeek), time.ErrInvalidDayOfWeek)
	assert.ErrorIs(t, r.AddDay(time.DayOfWeek(8)), time.ErrInvalidDayOfWeek)
	assert.Equal(t, []time.DayOfWeek{time.Monday, time.Wednesday, time.Saturday}, r.Days, "invalid days are not added")
}

func Test_DayOfWeekRange_AddRange(t *testing.T) {
	tcs := []struct {
		name       string
		start, end time.DayOfWeek
		expected   []time.DayOfWeek
	}{
		{
			name:     "weekdays",
			start:    time.Monday,
			end:      time.Friday,
			expected: []time.DayOfWeek{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		{
			name:     "single day",
			start:    time.Saturday,
			end:      time.Saturday,
			expected: []time.DayOfWeek{time.Saturday},
		},
		{
			name:     "whole week",
			start:    time.Monday,
			end:      time.Sunday,
			expected: time.DaysOfWeek,
		},
		{
			name:     "wraps past sunday",
			start:    time.Friday,
			end:      time.Monday,
			expected: []time.DayOfWeek{time.Monday, time.Friday, time.Saturday, time.Sunday},
		},
		{
			name:     "wraps to the whole week",
			start:    time.Tuesday,
			end:      time.Monday,
			expected: time.DaysOfWeek,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var r time.DayOfWeekRange

			assert.NoError(tt, r.AddRange(tc.start, tc.end))
			assert.Equal(tt, tc.expected, r.Days)
			assert.NoError(tt, r.Validate())
		})
	}
}

func Test_DayOfWeekRange_AddRange_Invalid(t *testing.T) {
	r := time.MustNewDayOfWeekRange(time.Sunday)

	assert.ErrorIs(t, r.AddRange(time.ZeroDayOfWeek, time.Friday), time.ErrInvalidDayOfWeek)
	assert.ErrorIs(t, r.AddRange(time.Monday, time.DayOfWeek(8)), time.ErrInvalidDayOfWeek)
	assert.Equal(t, []time.DayOfWeek{time.Sunday}, r.Days, "nothing is added")
}

func Test_DayOfWeekRange_WorkingDays(t *testing.T) {
	var r time.DayOfWeekRange

	assert.NoError(t, r.AddRange(time.Monday, time.Friday))
	assert.NoError(t, r.AddDay(time.Saturday))
	assert.Equal(t, []time.DayOfWeek{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}, r.Days)

	assert.NoError(t, r.AddRange(time.Wednesday, time.Saturday), "overlapping ranges store each day once")
	assert.Len(t, r.Days, 6)
}

func Test_DayOfWeekRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.DayOfWeekRange
		errs  []error
	}{
		{name: "empty", input: time.DayOfWeekRange{}},
		{name: "valid", input: time.DayOfWeekRange{Days: []time.DayOfWeek{time.Monday, time.Saturday}}},
		{
			name:  "invalid day",
			input: time.DayOfWeekRange{Days: []time.DayOfWeek{time.Monday, time.ZeroDayOfWeek}},
			errs:  []error{time.ErrInvalidDayOfWeekRange, time.ErrDays, time.ErrInvalidDayOfWeek},
		},
		{
			name:  "duplicate day",
			input: time.DayOfWeekRange{Days: []time.DayOfWeek{time.Monday, time.Tuesday, time.Monday}},
			errs:  []error{time.ErrInvalidDayOfWeekRange, time.ErrDays, time.ErrDuplicateDay},
		},
		{
			name:  "every failure is reported",
			input: time.DayOfWeekRange{Days: []time.DayOfWeek{time.Friday, time.Friday, time.DayOfWeek(9)}},
			errs:  []error{time.ErrDuplicateDay, time.ErrInvalidDayOfWeek},
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

func Test_NewDayOfWeekRange(t *testing.T) {
	actual, err := time.NewDayOfWeekRange(time.Saturday, time.Monday, time.Saturday)
	assert.NoError(t, err)
	assert.Equal(t, []time.DayOfWeek{time.Monday, time.Saturday}, actual.Days)

	actual, err = time.NewDayOfWeekRange()
	assert.NoError(t, err)
	assert.Empty(t, actual.Days)

	_, err = time.NewDayOfWeekRange(time.Monday, time.ZeroDayOfWeek)
	assert.ErrorIs(t, err, time.ErrInvalidDayOfWeek)

	assert.Panics(t, func() { time.MustNewDayOfWeekRange(time.DayOfWeek(8)) })
}

func Test_DayOfWeekRange_JSON(t *testing.T) {
	expected := time.MustNewDayOfWeekRange(time.Monday, time.Tuesday, time.Saturday)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"days":["Monday","Tuesday","Saturday"]}`, string(data))

	var actual time.DayOfWeekRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	data, err = json.Marshal(time.DayOfWeekRange{})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"days":[]}`, string(data), "empty encodes as []")
}

func Test_DayOfWeekRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewDayOfWeekRange(time.Sunday)

	tcs := []struct {
		name     string
		input    string
		start    time.DayOfWeekRange
		expected []time.DayOfWeek
		err      error
	}{
		{name: "days", input: `{"days":["mon","Fri"]}`, expected: []time.DayOfWeek{time.Monday, time.Friday}},
		{name: "empty", input: `{"days":[]}`, expected: []time.DayOfWeek{}},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid.Days},
		{name: "duplicate leaves value unchanged", input: `{"days":["Mon","Monday"]}`, start: valid, expected: valid.Days, err: time.ErrDuplicateDay},
		{name: "invalid day leaves value unchanged", input: `{"days":["Funday"]}`, start: valid, expected: valid.Days, err: time.ErrInvalidDayOfWeek},
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

			assert.Equal(tt, tc.expected, actual.Days)
		})
	}
}
