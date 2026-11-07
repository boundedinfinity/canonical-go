package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_MonthRange_AddMonth(t *testing.T) {
	var r time.MonthRange

	assert.NoError(t, r.AddMonth(time.December))
	assert.NoError(t, r.AddMonth(time.January))
	assert.NoError(t, r.AddMonth(time.June))
	assert.NoError(t, r.AddMonth(time.January), "adding a present month does nothing")
	assert.Equal(t, []time.Month{time.January, time.June, time.December}, r.Months, "ISO order, no repeats")

	assert.ErrorIs(t, r.AddMonth(time.ZeroMonth), time.ErrInvalidMonth)
	assert.ErrorIs(t, r.AddMonth(time.Month(13)), time.ErrInvalidMonth)
	assert.Equal(t, []time.Month{time.January, time.June, time.December}, r.Months, "invalid months are not added")
}

func Test_MonthRange_AddRange(t *testing.T) {
	tcs := []struct {
		name       string
		start, end time.Month
		expected   []time.Month
	}{
		{
			name:     "spring",
			start:    time.March,
			end:      time.May,
			expected: []time.Month{time.March, time.April, time.May},
		},
		{
			name:     "single month",
			start:    time.December,
			end:      time.December,
			expected: []time.Month{time.December},
		},
		{
			name:     "whole year",
			start:    time.January,
			end:      time.December,
			expected: time.MonthsOfYear,
		},
		{
			name:     "wraps past december",
			start:    time.November,
			end:      time.February,
			expected: []time.Month{time.January, time.February, time.November, time.December},
		},
		{
			name:     "wraps to the whole year",
			start:    time.February,
			end:      time.January,
			expected: time.MonthsOfYear,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var r time.MonthRange

			assert.NoError(tt, r.AddRange(tc.start, tc.end))
			assert.Equal(tt, tc.expected, r.Months)
			assert.NoError(tt, r.Validate())
		})
	}
}

func Test_MonthRange_AddRange_Invalid(t *testing.T) {
	r := time.MustNewMonthRange(time.July)

	assert.ErrorIs(t, r.AddRange(time.ZeroMonth, time.March), time.ErrInvalidMonth)
	assert.ErrorIs(t, r.AddRange(time.January, time.Month(13)), time.ErrInvalidMonth)
	assert.Equal(t, []time.Month{time.July}, r.Months, "nothing is added")
}

func Test_MonthRange_Season(t *testing.T) {
	var r time.MonthRange

	assert.NoError(t, r.AddRange(time.April, time.September))
	assert.NoError(t, r.AddMonth(time.December))
	assert.Equal(t, []time.Month{time.April, time.May, time.June, time.July, time.August, time.September, time.December}, r.Months)

	assert.NoError(t, r.AddRange(time.August, time.December), "overlapping ranges store each month once")
	assert.Len(t, r.Months, 9)
}

func Test_MonthRange_Contains(t *testing.T) {
	var season time.MonthRange
	assert.NoError(t, season.AddRange(time.November, time.February))

	tcs := []struct {
		name     string
		input    time.Month
		expected bool
		err      error
	}{
		{name: "range start", input: time.November, expected: true},
		{name: "after the wrap", input: time.January, expected: true},
		{name: "range end", input: time.February, expected: true},
		{name: "just after", input: time.March},
		{name: "just before", input: time.October},
		{name: "zero month", input: time.ZeroMonth, err: time.ErrInvalidMonth},
		{name: "out of range", input: time.Month(13), err: time.ErrInvalidMonth},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := season.Contains(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.False(tt, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}

	actual, err := time.MonthRange{}.Contains(time.January)
	assert.NoError(t, err)
	assert.False(t, actual, "an empty range contains nothing")
}

func Test_MonthRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.MonthRange
		errs  []error
	}{
		{name: "empty", input: time.MonthRange{}},
		{name: "valid", input: time.MonthRange{Months: []time.Month{time.January, time.December}}},
		{
			name:  "invalid month",
			input: time.MonthRange{Months: []time.Month{time.January, time.ZeroMonth}},
			errs:  []error{time.ErrInvalidMonthRange, time.ErrMonths, time.ErrInvalidMonth},
		},
		{
			name:  "duplicate month",
			input: time.MonthRange{Months: []time.Month{time.January, time.February, time.January}},
			errs:  []error{time.ErrInvalidMonthRange, time.ErrMonths, time.ErrDuplicateMonth},
		},
		{
			name:  "every failure is reported",
			input: time.MonthRange{Months: []time.Month{time.May, time.May, time.Month(13)}},
			errs:  []error{time.ErrDuplicateMonth, time.ErrInvalidMonth},
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

func Test_NewMonthRange(t *testing.T) {
	actual, err := time.NewMonthRange(time.December, time.January, time.December)
	assert.NoError(t, err)
	assert.Equal(t, []time.Month{time.January, time.December}, actual.Months)

	actual, err = time.NewMonthRange()
	assert.NoError(t, err)
	assert.Empty(t, actual.Months)

	_, err = time.NewMonthRange(time.January, time.ZeroMonth)
	assert.ErrorIs(t, err, time.ErrInvalidMonth)

	assert.Panics(t, func() { time.MustNewMonthRange(time.Month(13)) })
}

func Test_MonthRange_JSON(t *testing.T) {
	expected := time.MustNewMonthRange(time.January, time.February, time.December)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"months":["January","February","December"]}`, string(data))

	var actual time.MonthRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	data, err = json.Marshal(time.MonthRange{})
	assert.NoError(t, err)
	assert.JSONEq(t, `{"months":[]}`, string(data), "empty encodes as []")
}

func Test_MonthRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewMonthRange(time.July)

	tcs := []struct {
		name     string
		input    string
		start    time.MonthRange
		expected []time.Month
		err      error
	}{
		{name: "months", input: `{"months":["jan","Dec"]}`, expected: []time.Month{time.January, time.December}},
		{name: "empty", input: `{"months":[]}`, expected: []time.Month{}},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid.Months},
		{name: "duplicate leaves value unchanged", input: `{"months":["Jan","January"]}`, start: valid, expected: valid.Months, err: time.ErrDuplicateMonth},
		{name: "invalid month leaves value unchanged", input: `{"months":["Smarch"]}`, start: valid, expected: valid.Months, err: time.ErrInvalidMonth},
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

			assert.Equal(tt, tc.expected, actual.Months)
		})
	}
}
