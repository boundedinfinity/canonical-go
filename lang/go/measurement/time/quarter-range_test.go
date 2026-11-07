package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_QuarterRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.QuarterRange
		errs  []error
	}{
		{name: "valid", input: time.QuarterRange{Start: time.Q2, End: time.Q3}},
		{name: "start equals end", input: time.QuarterRange{Start: time.Q2, End: time.Q2}},
		{name: "whole year", input: time.QuarterRange{Start: time.Q1, End: time.Q4}},
		{
			name:  "missing start",
			input: time.QuarterRange{End: time.Q3},
			errs:  []error{time.ErrInvalidQuarterRange, time.ErrStart, time.ErrInvalidQuarter},
		},
		{
			name:  "missing end",
			input: time.QuarterRange{Start: time.Q2},
			errs:  []error{time.ErrInvalidQuarterRange, time.ErrEnd, time.ErrInvalidQuarter, time.ErrOrder},
		},
		{
			name:  "missing both",
			input: time.QuarterRange{},
			errs:  []error{time.ErrInvalidQuarterRange, time.ErrStart, time.ErrEnd},
		},
		{
			name:  "start after end",
			input: time.QuarterRange{Start: time.Q3, End: time.Q2},
			errs:  []error{time.ErrInvalidQuarterRange, time.ErrOrder},
		},
		{
			name:  "does not wrap past Q4",
			input: time.QuarterRange{Start: time.Q4, End: time.Q1},
			errs:  []error{time.ErrInvalidQuarterRange, time.ErrOrder},
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

func Test_QuarterRange_Operations(t *testing.T) {
	r := time.MustNewQuarterRange(time.Q2, time.Q3)

	tcs := []struct {
		name    string
		input   time.Quarter
		before  bool
		between bool
		after   bool
	}{
		{name: "before start", input: time.Q1, before: true},
		{name: "at start", input: time.Q2, between: true},
		{name: "at end", input: time.Q3, between: true},
		{name: "after end", input: time.Q4, after: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.before, r.Before(tc.input), "Before")
			assert.Equal(tt, tc.between, r.Between(tc.input), "Between")
			assert.Equal(tt, tc.after, r.After(tc.input), "After")
		})
	}
}

func Test_QuarterRange_Overlaps(t *testing.T) {
	r := time.MustNewQuarterRange(time.Q2, time.Q3)

	tcs := []struct {
		name     string
		input    time.QuarterRange
		expected bool
	}{
		{name: "entirely before", input: time.MustNewQuarterRange(time.Q1, time.Q1)},
		{name: "touches start", input: time.MustNewQuarterRange(time.Q1, time.Q2), expected: true},
		{name: "inside", input: time.MustNewQuarterRange(time.Q3, time.Q3), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: time.MustNewQuarterRange(time.Q1, time.Q4), expected: true},
		{name: "touches end", input: time.MustNewQuarterRange(time.Q3, time.Q4), expected: true},
		{name: "entirely after", input: time.MustNewQuarterRange(time.Q4, time.Q4)},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_QuarterRange_SingleQuarter(t *testing.T) {
	r := time.MustNewQuarterRange(time.Q2, time.Q2)

	assert.True(t, r.Between(time.Q2))
	assert.True(t, r.Before(time.Q1))
	assert.True(t, r.After(time.Q3))
}

func Test_NewQuarterRange(t *testing.T) {
	actual, err := time.NewQuarterRange(time.Q1, time.Q3)
	assert.NoError(t, err)
	assert.Equal(t, time.Q1, actual.Start)
	assert.Equal(t, time.Q3, actual.End)

	_, err = time.NewQuarterRange(time.Q3, time.Q1)
	assert.ErrorIs(t, err, time.ErrOrder)

	assert.Panics(t, func() { time.MustNewQuarterRange(time.Q3, time.Q1) })
}

func Test_QuarterRange_JSON(t *testing.T) {
	expected := time.MustNewQuarterRange(time.Q2, time.Q4)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"Q2","end":"Q4"}`, string(data))

	var actual time.QuarterRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func Test_QuarterRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewQuarterRange(time.Q2, time.Q3)

	tcs := []struct {
		name     string
		input    string
		start    time.QuarterRange
		expected time.QuarterRange
		err      error
	}{
		{name: "valid", input: `{"start":"Q2","end":"Q3"}`, expected: valid},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid},
		{name: "start after end leaves value unchanged", input: `{"start":"Q3","end":"Q2"}`, start: valid, expected: valid, err: time.ErrOrder},
		{name: "missing start", input: `{"end":"Q3"}`, err: time.ErrStart},
		{name: "missing end", input: `{"start":"Q2"}`, err: time.ErrEnd},
		{name: "invalid quarter", input: `{"start":"Q5","end":"Q3"}`, err: time.ErrInvalidQuarter},
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
