package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_Week_Days(t *testing.T) {
	tcs := []struct {
		name     string
		input    time.Week
		first    time.DayOfWeek
		last     time.DayOfWeek
		expected []time.DayOfWeek
	}{
		{
			name:     "zero value starts on Monday",
			input:    time.Week{},
			first:    time.Monday,
			last:     time.Sunday,
			expected: []time.DayOfWeek{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
		},
		{
			name:     "default week",
			input:    time.DefaultWeek,
			first:    time.Monday,
			last:     time.Sunday,
			expected: time.DaysOfWeek,
		},
		{
			name:     "sunday start",
			input:    time.MustNewWeek(time.Sunday),
			first:    time.Sunday,
			last:     time.Saturday,
			expected: []time.DayOfWeek{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
		},
		{
			name:     "saturday start",
			input:    time.MustNewWeek(time.Saturday),
			first:    time.Saturday,
			last:     time.Friday,
			expected: []time.DayOfWeek{time.Saturday, time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		{
			name:     "wednesday start",
			input:    time.MustNewWeek(time.Wednesday),
			first:    time.Wednesday,
			last:     time.Tuesday,
			expected: []time.DayOfWeek{time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday, time.Monday, time.Tuesday},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.first, tc.input.First())
			assert.Equal(tt, tc.last, tc.input.Last())
			assert.Equal(tt, tc.expected, tc.input.Days())

			for i, d := range tc.expected {
				assert.Equal(tt, i+1, tc.input.Position(d), "Position(%s)", d)
			}
		})
	}
}

func Test_Week_Position_Invalid(t *testing.T) {
	assert.Equal(t, 0, time.DefaultWeek.Position(time.ZeroDayOfWeek))
	assert.Equal(t, 0, time.DefaultWeek.Position(time.DayOfWeek(8)))
}

func Test_Week_Compare(t *testing.T) {
	monday := time.DefaultWeek
	sunday := time.MustNewWeek(time.Sunday)

	assert.Equal(t, -1, monday.Compare(time.Monday, time.Sunday))
	assert.Equal(t, 1, sunday.Compare(time.Monday, time.Sunday), "Sunday comes first in a Sunday week")
	assert.Equal(t, 0, sunday.Compare(time.Friday, time.Friday))
	assert.Equal(t, -1, sunday.Compare(time.Sunday, time.Saturday))
	assert.Equal(t, 1, monday.Compare(time.Sunday, time.Saturday))
}

func Test_Week_Number_Unchanged(t *testing.T) {
	sunday := time.MustNewWeek(time.Sunday)

	assert.Equal(t, 7, sunday.First().Number(), "Number stays ISO 8601")
	assert.Equal(t, 1, sunday.Position(time.Sunday))
}

func Test_Week_Validate(t *testing.T) {
	assert.NoError(t, time.Week{}.Validate())
	assert.NoError(t, time.DefaultWeek.Validate())
	assert.NoError(t, time.MustNewWeek(time.Sunday).Validate())

	err := time.Week{Start: optioner.Some(time.ZeroDayOfWeek)}.Validate()
	assert.ErrorIs(t, err, time.ErrInvalidWeek)
	assert.ErrorIs(t, err, time.ErrStart)
	assert.ErrorIs(t, err, time.ErrInvalidDayOfWeek)
}

func Test_NewWeek(t *testing.T) {
	actual, err := time.NewWeek(time.Sunday)
	assert.NoError(t, err)
	assert.Equal(t, time.Sunday, actual.First())

	actual, err = time.NewWeek(time.ZeroDayOfWeek)
	assert.NoError(t, err, "zero start is treated as not present")
	assert.False(t, actual.Start.Defined())
	assert.Equal(t, time.Monday, actual.First())

	_, err = time.NewWeek(time.DayOfWeek(8))
	assert.ErrorIs(t, err, time.ErrInvalidWeek)

	assert.Panics(t, func() { time.MustNewWeek(time.DayOfWeek(8)) })
}

func Test_Week_JSON(t *testing.T) {
	data, err := json.Marshal(time.MustNewWeek(time.Sunday))
	assert.NoError(t, err)
	assert.JSONEq(t, `{"start":"Sunday"}`, string(data))

	data, err = json.Marshal(time.Week{})
	assert.NoError(t, err)
	assert.JSONEq(t, `{}`, string(data), "absent start is omitted")

	var actual time.Week
	assert.NoError(t, json.Unmarshal([]byte(`{"start":"sun"}`), &actual))
	assert.Equal(t, time.Sunday, actual.First())

	actual = time.Week{}
	assert.NoError(t, json.Unmarshal([]byte(`{}`), &actual))
	assert.Equal(t, time.Monday, actual.First())
}

func Test_Week_UnmarshalJSON(t *testing.T) {
	sunday := time.MustNewWeek(time.Sunday)

	tcs := []struct {
		name     string
		input    string
		start    time.Week
		expected time.DayOfWeek
		err      error
	}{
		{name: "start", input: `{"start":"Saturday"}`, expected: time.Saturday},
		{name: "no start", input: `{}`, expected: time.Monday},
		{name: "null leaves value unchanged", input: `null`, start: sunday, expected: time.Sunday},
		{name: "invalid start leaves value unchanged", input: `{"start":"Funday"}`, start: sunday, expected: time.Sunday, err: time.ErrInvalidDayOfWeek},
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

			assert.Equal(tt, tc.expected, actual.First())
		})
	}
}
