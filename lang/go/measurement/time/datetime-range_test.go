package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

var (
	rangeStart = time.MustParse("2026-11-07T09:00:00Z")
	rangeEnd   = time.MustParse("2026-11-07T17:00:00Z")
)

func Test_DateTimeRange_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input time.DateTimeRange
		errs  []error
	}{
		{
			name:  "valid",
			input: time.DateTimeRange{Start: rangeStart, End: rangeEnd},
		},
		{
			name:  "start equals end",
			input: time.DateTimeRange{Start: rangeStart, End: rangeStart},
		},
		{
			name: "same instant with different offsets",
			input: time.DateTimeRange{
				Start: time.MustParse("2026-11-07T09:00:00-05:00"),
				End:   time.MustParse("2026-11-07T14:00:00Z"),
			},
		},
		{
			name:  "missing start",
			input: time.DateTimeRange{End: rangeEnd},
			errs:  []error{time.ErrInvalidDateTimeRange, time.ErrStart, time.ErrInvalid},
		},
		{
			name:  "missing end",
			input: time.DateTimeRange{Start: rangeStart},
			errs:  []error{time.ErrInvalidDateTimeRange, time.ErrEnd, time.ErrInvalid, time.ErrOrder},
		},
		{
			name:  "missing both",
			input: time.DateTimeRange{},
			errs:  []error{time.ErrInvalidDateTimeRange, time.ErrStart, time.ErrEnd},
		},
		{
			name:  "start after end",
			input: time.DateTimeRange{Start: rangeEnd, End: rangeStart},
			errs:  []error{time.ErrInvalidDateTimeRange, time.ErrOrder},
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

func Test_DateTimeRange_Operations(t *testing.T) {
	r := time.MustNewRange(rangeStart, rangeEnd)

	tcs := []struct {
		name    string
		input   time.DateTime
		before  bool
		between bool
		after   bool
	}{
		{name: "before start", input: time.MustParse("2026-11-07T08:59:59Z"), before: true},
		{name: "at start", input: rangeStart, between: true},
		{name: "inside", input: time.MustParse("2026-11-07T12:00:00Z"), between: true},
		{name: "at end", input: rangeEnd, between: true},
		{name: "after end", input: time.MustParse("2026-11-07T17:00:00.001Z"), after: true},
		{name: "start in another offset", input: time.MustParse("2026-11-07T04:00:00-05:00"), between: true},
		{name: "before start in another offset", input: time.MustParse("2026-11-07T03:59:59-05:00"), before: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.before, r.Before(tc.input), "Before")
			assert.Equal(tt, tc.between, r.Between(tc.input), "Between")
			assert.Equal(tt, tc.after, r.After(tc.input), "After")
		})
	}
}

func Test_DateTimeRange_Overlaps(t *testing.T) {
	r := time.MustNewRange(rangeStart, rangeEnd)
	at := time.MustParse

	tcs := []struct {
		name     string
		input    time.DateTimeRange
		expected bool
	}{
		{name: "entirely before", input: time.MustNewRange(at("2026-11-07T06:00:00Z"), at("2026-11-07T08:59:59Z"))},
		{name: "touches start", input: time.MustNewRange(at("2026-11-07T06:00:00Z"), rangeStart), expected: true},
		{name: "overlaps start", input: time.MustNewRange(at("2026-11-07T06:00:00Z"), at("2026-11-07T10:00:00Z")), expected: true},
		{name: "inside", input: time.MustNewRange(at("2026-11-07T10:00:00Z"), at("2026-11-07T11:00:00Z")), expected: true},
		{name: "same range", input: r, expected: true},
		{name: "contains", input: time.MustNewRange(at("2026-11-07T06:00:00Z"), at("2026-11-07T20:00:00Z")), expected: true},
		{name: "overlaps end", input: time.MustNewRange(at("2026-11-07T16:00:00Z"), at("2026-11-07T20:00:00Z")), expected: true},
		{name: "touches end", input: time.MustNewRange(rangeEnd, at("2026-11-07T20:00:00Z")), expected: true},
		{name: "entirely after", input: time.MustNewRange(at("2026-11-07T17:00:01Z"), at("2026-11-07T20:00:00Z"))},
		{name: "single instant inside", input: time.MustNewRange(at("2026-11-07T12:00:00Z"), at("2026-11-07T12:00:00Z")), expected: true},
		{name: "touches end in another offset", input: time.MustNewRange(at("2026-11-07T12:00:00-05:00"), at("2026-11-07T15:00:00-05:00")), expected: true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.expected, r.Overlaps(tc.input))
			assert.Equal(tt, tc.expected, tc.input.Overlaps(r), "symmetric")
		})
	}
}

func Test_DateTimeRange_SingleInstant(t *testing.T) {
	r := time.MustNewRange(rangeStart, rangeStart)

	assert.True(t, r.Between(rangeStart))
	assert.True(t, r.After(rangeEnd))
	assert.True(t, r.Before(time.MustParse("2026-11-07T08:00:00Z")))
}

func Test_NewRange(t *testing.T) {
	actual, err := time.NewRange(rangeStart, rangeEnd)
	assert.NoError(t, err)
	assert.Equal(t, rangeStart, actual.Start)
	assert.Equal(t, rangeEnd, actual.End)

	_, err = time.NewRange(rangeEnd, rangeStart)
	assert.ErrorIs(t, err, time.ErrOrder)

	assert.Panics(t, func() { time.MustNewRange(rangeEnd, rangeStart) })
}

func Test_DateTimeRange_JSON(t *testing.T) {
	expected := time.MustNewRange(rangeStart, rangeEnd)

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"2026-11-07T09:00:00Z","end":"2026-11-07T17:00:00Z"}`, string(data))

	var actual time.DateTimeRange
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.True(t, expected.Start.Equal(actual.Start))
	assert.True(t, expected.End.Equal(actual.End))
}

func Test_DateTimeRange_UnmarshalJSON(t *testing.T) {
	valid := time.MustNewRange(rangeStart, rangeEnd)

	tcs := []struct {
		name     string
		input    string
		start    time.DateTimeRange
		expected time.DateTimeRange
		err      error
	}{
		{
			name:     "valid",
			input:    `{"start":"2026-11-07T09:00:00Z","end":"2026-11-07T17:00:00Z"}`,
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
			input:    `{"start":"2026-11-07T17:00:00Z","end":"2026-11-07T09:00:00Z"}`,
			start:    valid,
			expected: valid,
			err:      time.ErrOrder,
		},
		{
			name:  "missing start",
			input: `{"end":"2026-11-07T17:00:00Z"}`,
			err:   time.ErrStart,
		},
		{
			name:  "missing end",
			input: `{"start":"2026-11-07T09:00:00Z"}`,
			err:   time.ErrEnd,
		},
		{
			name:  "invalid datetime",
			input: `{"start":"nope","end":"2026-11-07T17:00:00Z"}`,
			err:   time.ErrInvalid,
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
