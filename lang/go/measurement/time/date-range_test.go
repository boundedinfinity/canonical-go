package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

var (
	dateRangeStart = time.MustParseDate("2026-11-01")
	dateRangeEnd   = time.MustParseDate("2026-11-30")
)

func Test_DateRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.DateRange
		errs  []error
	}{
		{
			name:  "valid",
			input: time.DateRange{Start: dateRangeStart, End: dateRangeEnd},
		},
		{
			name:  "start equals end",
			input: time.DateRange{Start: dateRangeStart, End: dateRangeStart},
		},
		{
			name:  "missing start",
			input: time.DateRange{End: dateRangeEnd},
			errs:  []error{time.ErrInvalidDateRange, time.ErrStart, time.ErrInvalidDate},
		},
		{
			name:  "missing end",
			input: time.DateRange{Start: dateRangeStart},
			errs:  []error{time.ErrInvalidDateRange, time.ErrEnd, time.ErrInvalidDate, time.ErrOrder},
		},
		{
			name:  "missing both",
			input: time.DateRange{},
			errs:  []error{time.ErrInvalidDateRange, time.ErrStart, time.ErrEnd},
		},
		{
			name:  "start after end",
			input: time.DateRange{Start: dateRangeEnd, End: dateRangeStart},
			errs:  []error{time.ErrInvalidDateRange, time.ErrOrder},
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

func Test_DateRange_Operations(t *testing.T) {
	r := time.MustNewDateRange(dateRangeStart, dateRangeEnd)

	tcs := []struct {
		name    string
		input   time.Date
		before  bool
		between bool
		after   bool
	}{
		{name: "before start", input: time.MustParseDate("2026-10-31"), before: true},
		{name: "at start", input: dateRangeStart, between: true},
		{name: "inside", input: time.MustParseDate("2026-11-15"), between: true},
		{name: "at end", input: dateRangeEnd, between: true},
		{name: "after end", input: time.MustParseDate("2026-12-01"), after: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.before, r.Before(tc.input), "Before")
			assert.Equal(tt, tc.between, r.Between(tc.input), "Between")
			assert.Equal(tt, tc.after, r.After(tc.input), "After")
		})
	}
}

func Test_DateRange_Overlaps(t *testing.T) {
	r := time.MustNewDateRange(dateRangeStart, dateRangeEnd)
	on := time.MustParseDate

	tcs := []struct {
		name     string
		input    time.DateRange
		expected bool
	}{
		{name: "entirely before", input: time.MustNewDateRange(on("2026-10-01"), on("2026-10-31"))},
		{name: "touches start", input: time.MustNewDateRange(on("2026-10-01"), dateRangeStart), expected: true},
		{name: "overlaps start", input: time.MustNewDateRange(on("2026-10-01"), on("2026-11-10")), expected: true},
		{name: "inside", input: time.MustNewDateRange(on("2026-11-10"), on("2026-11-20")), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: time.MustNewDateRange(on("2026-10-01"), on("2026-12-31")), expected: true},
		{name: "overlaps end", input: time.MustNewDateRange(on("2026-11-20"), on("2026-12-31")), expected: true},
		{name: "touches end", input: time.MustNewDateRange(dateRangeEnd, on("2026-12-31")), expected: true},
		{name: "entirely after", input: time.MustNewDateRange(on("2026-12-01"), on("2026-12-31"))},
		{name: "single day inside", input: time.MustNewDateRange(on("2026-11-15"), on("2026-11-15")), expected: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_DateRange_SingleDay(t *testing.T) {
	r := time.MustNewDateRange(dateRangeStart, dateRangeStart)

	assert.True(t, r.Between(dateRangeStart))
	assert.True(t, r.Before(time.MustParseDate("2026-10-31")))
	assert.True(t, r.After(time.MustParseDate("2026-11-02")))
}

func Test_NewDateRange(t *testing.T) {
	actual, err := time.NewDateRange(dateRangeStart, dateRangeEnd)
	assert.NoError(t, err)
	assert.Equal(t, dateRangeStart, actual.Start)
	assert.Equal(t, dateRangeEnd, actual.End)

	_, err = time.NewDateRange(dateRangeEnd, dateRangeStart)
	assert.ErrorIs(t, err, time.ErrOrder)

	assert.Panics(t, func() { time.MustNewDateRange(dateRangeEnd, dateRangeStart) })
}

func Test_DateRange_JSON(t *testing.T) {
	expected := time.MustNewDateRange(dateRangeStart, dateRangeEnd)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"2026-11-01","end":"2026-11-30"}`, string(data))

	var actual time.DateRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func Test_DateRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewDateRange(dateRangeStart, dateRangeEnd)

	tcs := []struct {
		name     string
		input    string
		start    time.DateRange
		expected time.DateRange
		err      error
	}{
		{
			name:     "valid",
			input:    `{"start":"2026-11-01","end":"2026-11-30"}`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:     "start after end leaves value unchanged",
			input:    `{"start":"2026-11-30","end":"2026-11-01"}`,
			start:    valid,
			expected: valid,
			err:      time.ErrOrder,
		},
		{
			name:  "missing start",
			input: `{"end":"2026-11-30"}`,
			err:   time.ErrStart,
		},
		{
			name:  "missing end",
			input: `{"start":"2026-11-01"}`,
			err:   time.ErrEnd,
		},
		{
			name:  "invalid date",
			input: `{"start":"nope","end":"2026-11-30"}`,
			err:   time.ErrInvalidDate,
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
