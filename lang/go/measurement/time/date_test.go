package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseDate(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "valid",
			input:    "2026-11-07",
			expected: "2026-11-07",
		},
		{
			name:     "leap day",
			input:    "2028-02-29",
			expected: "2028-02-29",
		},
		{
			name:  "empty",
			input: "",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "not a leap year",
			input: "2026-02-29",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "out of range month",
			input: "2026-13-07",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "out of range day",
			input: "2026-11-31",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "missing zero padding",
			input: "2026-1-7",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "with time",
			input: "2026-11-07T14:30:00Z",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "slash separator",
			input: "2026/11/07",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "zero value",
			input: "0001-01-01",
			err:   time.ErrInvalidDate,
		},
		{
			name:  "not a date",
			input: "not-a-date",
			err:   time.ErrInvalidDate,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseDate(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.True(tt, actual.IsZero())
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual.String())
		})
	}
}

func Test_NewDate(t *testing.T) {
	loc := gotime.FixedZone("", -5*60*60)

	// 23:30 on Nov 7 in UTC-5 is already Nov 8 in UTC; the date in t's own
	// location is kept.
	actual, err := time.NewDate(gotime.Date(2026, 11, 7, 23, 30, 0, 0, loc))
	assert.NoError(t, err)
	assert.Equal(t, "2026-11-07", actual.String())
	assert.Equal(t, time.MustParseDate("2026-11-07"), actual, "comparable with ==")
	assert.Equal(t, gotime.Date(2026, 11, 7, 0, 0, 0, 0, gotime.UTC), actual.Time())

	_, err = time.NewDate(gotime.Time{})
	assert.ErrorIs(t, err, time.ErrInvalidDate)

	assert.Panics(t, func() { time.MustNewDate(gotime.Time{}) })
}

func Test_Today(t *testing.T) {
	today := time.Today()

	assert.NoError(t, today.Validate())
	assert.Equal(t, gotime.Now().Format(gotime.DateOnly), today.String())
}

func Test_Date_Validate(t *testing.T) {
	assert.NoError(t, time.MustParseDate("2026-11-07").Validate())
	assert.ErrorIs(t, time.ZeroDate.Validate(), time.ErrInvalidDate)
}

func Test_Date_JSON(t *testing.T) {
	type wrapper struct {
		On time.Date `json:"on"`
	}

	expected := wrapper{On: time.MustParseDate("2026-11-07")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"on":"2026-11-07"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"on":"nope"}`), &actual), time.ErrInvalidDate)
}

func Test_Date_UnmarshalJSON(t *testing.T) {
	valid := time.MustParseDate("2026-11-07")

	tcs := []struct {
		name     string
		input    string
		start    time.Date
		expected time.Date
		err      error
	}{
		{
			name:     "string",
			input:    `"2026-11-07"`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:  "invalid date string",
			input: `"nope"`,
			err:   time.ErrInvalidDate,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   time.ErrInvalidDate,
		},
		{
			name:  "number",
			input: `123`,
			err:   time.ErrInvalidDate,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual := tc.start
			err := actual.UnmarshalJSON([]byte(tc.input))

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_Date_Text(t *testing.T) {
	m := map[time.Date]string{
		time.MustParseDate("2026-11-07"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"2026-11-07":"a"}`, string(data))

	var actual map[time.Date]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
