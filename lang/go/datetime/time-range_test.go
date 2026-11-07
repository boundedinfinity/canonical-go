package datetime_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/datetime"
	"github.com/stretchr/testify/assert"
)

var (
	timeRangeStart = datetime.MustParseTime("09:00:00Z")
	timeRangeEnd   = datetime.MustParseTime("17:00:00Z")
)

func Test_TimeRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input datetime.TimeRange
		errs  []error
	}{
		{
			name:  "valid",
			input: datetime.TimeRange{Start: timeRangeStart, End: timeRangeEnd},
		},
		{
			name:  "start equals end",
			input: datetime.TimeRange{Start: timeRangeStart, End: timeRangeStart},
		},
		{
			name: "evening in a negative offset crosses midnight utc",
			input: datetime.TimeRange{
				Start: datetime.MustParseTime("18:00:00-05:00"),
				End:   datetime.MustParseTime("20:00:00-05:00"),
			},
		},
		{
			name: "mixed offsets compare by instant",
			input: datetime.TimeRange{
				Start: datetime.MustParseTime("09:00:00-05:00"),
				End:   datetime.MustParseTime("14:00:00Z"),
			},
		},
		{
			name:  "zero value is midnight to midnight",
			input: datetime.TimeRange{},
		},
		{
			name:  "missing end reads as midnight",
			input: datetime.TimeRange{Start: timeRangeStart},
			errs:  []error{datetime.ErrInvalidTimeRange, datetime.ErrOrder},
		},
		{
			name:  "start after end",
			input: datetime.TimeRange{Start: timeRangeEnd, End: timeRangeStart},
			errs:  []error{datetime.ErrInvalidTimeRange, datetime.ErrOrder},
		},
		{
			name: "overnight does not wrap",
			input: datetime.TimeRange{
				Start: datetime.MustParseTime("22:00:00Z"),
				End:   datetime.MustParseTime("02:00:00Z"),
			},
			errs: []error{datetime.ErrInvalidTimeRange, datetime.ErrOrder},
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

func Test_TimeRange_Operations(t *testing.T) {
	r := datetime.MustNewTimeRange(timeRangeStart, timeRangeEnd)

	tcs := []struct {
		name    string
		input   datetime.Time
		before  bool
		between bool
		after   bool
	}{
		{name: "before start", input: datetime.MustParseTime("08:59:59Z"), before: true},
		{name: "at start", input: timeRangeStart, between: true},
		{name: "inside", input: datetime.MustParseTime("12:00:00Z"), between: true},
		{name: "at end", input: timeRangeEnd, between: true},
		{name: "after end", input: datetime.MustParseTime("17:00:00.001Z"), after: true},
		{name: "start in another offset", input: datetime.MustParseTime("04:00:00-05:00"), between: true},
		{name: "before start in another offset", input: datetime.MustParseTime("03:59:59-05:00"), before: true},
		{name: "midnight", input: datetime.ZeroTime, before: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.before, r.Before(tc.input), "Before")
			assert.Equal(tt, tc.between, r.Between(tc.input), "Between")
			assert.Equal(tt, tc.after, r.After(tc.input), "After")
		})
	}
}

func Test_TimeRange_Operations_NegativeOffset(t *testing.T) {
	r := datetime.MustNewTimeRange(datetime.MustParseTime("18:00:00-05:00"), datetime.MustParseTime("20:00:00-05:00"))

	assert.True(t, r.Before(datetime.MustParseTime("17:59:59-05:00")))
	assert.True(t, r.Between(datetime.MustParseTime("19:00:00-05:00")))
	assert.True(t, r.After(datetime.MustParseTime("20:00:01-05:00")))
}

func Test_TimeRange_Overlaps(t *testing.T) {
	r := datetime.MustNewTimeRange(timeRangeStart, timeRangeEnd)
	at := datetime.MustParseTime

	tcs := []struct {
		name     string
		input    datetime.TimeRange
		expected bool
	}{
		{name: "entirely before", input: datetime.MustNewTimeRange(at("06:00:00Z"), at("08:59:59Z"))},
		{name: "touches start", input: datetime.MustNewTimeRange(at("06:00:00Z"), timeRangeStart), expected: true},
		{name: "overlaps start", input: datetime.MustNewTimeRange(at("06:00:00Z"), at("10:00:00Z")), expected: true},
		{name: "inside", input: datetime.MustNewTimeRange(at("10:00:00Z"), at("11:00:00Z")), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: datetime.MustNewTimeRange(at("06:00:00Z"), at("20:00:00Z")), expected: true},
		{name: "overlaps end", input: datetime.MustNewTimeRange(at("16:00:00Z"), at("20:00:00Z")), expected: true},
		{name: "touches end", input: datetime.MustNewTimeRange(timeRangeEnd, at("20:00:00Z")), expected: true},
		{name: "entirely after", input: datetime.MustNewTimeRange(at("17:00:01Z"), at("20:00:00Z"))},
		{name: "single instant inside", input: datetime.MustNewTimeRange(at("12:00:00Z"), at("12:00:00Z")), expected: true},
		{name: "touches end in another offset", input: datetime.MustNewTimeRange(at("12:00:00-05:00"), at("15:00:00-05:00")), expected: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_TimeRange_SingleInstant(t *testing.T) {
	r := datetime.MustNewTimeRange(timeRangeStart, timeRangeStart)

	assert.True(t, r.Between(timeRangeStart))
	assert.True(t, r.After(timeRangeEnd))
	assert.True(t, r.Before(datetime.MustParseTime("08:00:00Z")))
}

func Test_NewTimeRange(t *testing.T) {
	actual, err := datetime.NewTimeRange(timeRangeStart, timeRangeEnd)
	assert.NoError(t, err)
	assert.Equal(t, timeRangeStart, actual.Start)
	assert.Equal(t, timeRangeEnd, actual.End)

	_, err = datetime.NewTimeRange(timeRangeEnd, timeRangeStart)
	assert.ErrorIs(t, err, datetime.ErrOrder)

	assert.Panics(t, func() { datetime.MustNewTimeRange(timeRangeEnd, timeRangeStart) })
}

func Test_TimeRange_JSON(t *testing.T) {
	expected := datetime.MustNewTimeRange(timeRangeStart, timeRangeEnd)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"09:00:00Z","end":"17:00:00Z"}`, string(data))

	var actual datetime.TimeRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.True(t, expected.Start.Equal(actual.Start))
	assert.True(t, expected.End.Equal(actual.End))
}

func Test_TimeRange_UnmarshalJSON(t *testing.T) {
	valid := datetime.MustNewTimeRange(timeRangeStart, timeRangeEnd)

	tcs := []struct {
		name     string
		input    string
		start    datetime.TimeRange
		expected datetime.TimeRange
		err      error
	}{
		{
			name:     "valid",
			input:    `{"start":"09:00:00Z","end":"17:00:00Z"}`,
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
			input:    `{"start":"17:00:00Z","end":"09:00:00Z"}`,
			start:    valid,
			expected: valid,
			err:      datetime.ErrOrder,
		},
		{
			name:  "missing end reads as midnight",
			input: `{"start":"09:00:00Z"}`,
			err:   datetime.ErrOrder,
		},
		{
			name:  "invalid time",
			input: `{"start":"nope","end":"17:00:00Z"}`,
			err:   datetime.ErrInvalidTime,
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

			assert.True(tt, tc.expected.Start.Equal(actual.Start))
			assert.True(tt, tc.expected.End.Equal(actual.End))
		})
	}
}
