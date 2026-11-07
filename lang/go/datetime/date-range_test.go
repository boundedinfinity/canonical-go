package datetime_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/datetime"
	"github.com/stretchr/testify/assert"
)

var (
	dateRangeStart = datetime.MustParseDate("2026-11-01")
	dateRangeEnd   = datetime.MustParseDate("2026-11-30")
)

func Test_DateRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input datetime.DateRange
		errs  []error
	}{
		{
			name:  "valid",
			input: datetime.DateRange{Start: dateRangeStart, End: dateRangeEnd},
		},
		{
			name:  "start equals end",
			input: datetime.DateRange{Start: dateRangeStart, End: dateRangeStart},
		},
		{
			name:  "missing start",
			input: datetime.DateRange{End: dateRangeEnd},
			errs:  []error{datetime.ErrInvalidDateRange, datetime.ErrStart, datetime.ErrInvalidDate},
		},
		{
			name:  "missing end",
			input: datetime.DateRange{Start: dateRangeStart},
			errs:  []error{datetime.ErrInvalidDateRange, datetime.ErrEnd, datetime.ErrInvalidDate, datetime.ErrOrder},
		},
		{
			name:  "missing both",
			input: datetime.DateRange{},
			errs:  []error{datetime.ErrInvalidDateRange, datetime.ErrStart, datetime.ErrEnd},
		},
		{
			name:  "start after end",
			input: datetime.DateRange{Start: dateRangeEnd, End: dateRangeStart},
			errs:  []error{datetime.ErrInvalidDateRange, datetime.ErrOrder},
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
	r := datetime.MustNewDateRange(dateRangeStart, dateRangeEnd)

	tcs := []struct {
		name    string
		input   datetime.Date
		before  bool
		between bool
		after   bool
	}{
		{name: "before start", input: datetime.MustParseDate("2026-10-31"), before: true},
		{name: "at start", input: dateRangeStart, between: true},
		{name: "inside", input: datetime.MustParseDate("2026-11-15"), between: true},
		{name: "at end", input: dateRangeEnd, between: true},
		{name: "after end", input: datetime.MustParseDate("2026-12-01"), after: true},
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
	r := datetime.MustNewDateRange(dateRangeStart, dateRangeEnd)
	on := datetime.MustParseDate

	tcs := []struct {
		name     string
		input    datetime.DateRange
		expected bool
	}{
		{name: "entirely before", input: datetime.MustNewDateRange(on("2026-10-01"), on("2026-10-31"))},
		{name: "touches start", input: datetime.MustNewDateRange(on("2026-10-01"), dateRangeStart), expected: true},
		{name: "overlaps start", input: datetime.MustNewDateRange(on("2026-10-01"), on("2026-11-10")), expected: true},
		{name: "inside", input: datetime.MustNewDateRange(on("2026-11-10"), on("2026-11-20")), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: datetime.MustNewDateRange(on("2026-10-01"), on("2026-12-31")), expected: true},
		{name: "overlaps end", input: datetime.MustNewDateRange(on("2026-11-20"), on("2026-12-31")), expected: true},
		{name: "touches end", input: datetime.MustNewDateRange(dateRangeEnd, on("2026-12-31")), expected: true},
		{name: "entirely after", input: datetime.MustNewDateRange(on("2026-12-01"), on("2026-12-31"))},
		{name: "single day inside", input: datetime.MustNewDateRange(on("2026-11-15"), on("2026-11-15")), expected: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_DateRange_SingleDay(t *testing.T) {
	r := datetime.MustNewDateRange(dateRangeStart, dateRangeStart)

	assert.True(t, r.Between(dateRangeStart))
	assert.True(t, r.Before(datetime.MustParseDate("2026-10-31")))
	assert.True(t, r.After(datetime.MustParseDate("2026-11-02")))
}

func Test_NewDateRange(t *testing.T) {
	actual, err := datetime.NewDateRange(dateRangeStart, dateRangeEnd)
	assert.NoError(t, err)
	assert.Equal(t, dateRangeStart, actual.Start)
	assert.Equal(t, dateRangeEnd, actual.End)

	_, err = datetime.NewDateRange(dateRangeEnd, dateRangeStart)
	assert.ErrorIs(t, err, datetime.ErrOrder)

	assert.Panics(t, func() { datetime.MustNewDateRange(dateRangeEnd, dateRangeStart) })
}

func Test_DateRange_JSON(t *testing.T) {
	expected := datetime.MustNewDateRange(dateRangeStart, dateRangeEnd)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"2026-11-01","end":"2026-11-30"}`, string(data))

	var actual datetime.DateRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func Test_DateRange_UnmarshalJSON(t *testing.T) {
	valid := datetime.MustNewDateRange(dateRangeStart, dateRangeEnd)

	tcs := []struct {
		name     string
		input    string
		start    datetime.DateRange
		expected datetime.DateRange
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
			err:      datetime.ErrOrder,
		},
		{
			name:  "missing start",
			input: `{"end":"2026-11-30"}`,
			err:   datetime.ErrStart,
		},
		{
			name:  "missing end",
			input: `{"start":"2026-11-01"}`,
			err:   datetime.ErrEnd,
		},
		{
			name:  "invalid date",
			input: `{"start":"nope","end":"2026-11-30"}`,
			err:   datetime.ErrInvalidDate,
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
