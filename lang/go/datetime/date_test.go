package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/boundedinfinity/canonical-go/lang/go/datetime"
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
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "not a leap year",
			input: "2026-02-29",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "out of range month",
			input: "2026-13-07",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "out of range day",
			input: "2026-11-31",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "missing zero padding",
			input: "2026-1-7",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "with time",
			input: "2026-11-07T14:30:00Z",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "slash separator",
			input: "2026/11/07",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "zero value",
			input: "0001-01-01",
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "not a date",
			input: "not-a-date",
			err:   datetime.ErrInvalidDate,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := datetime.ParseDate(tc.input)

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
	loc := time.FixedZone("", -5*60*60)

	// 23:30 on Nov 7 in UTC-5 is already Nov 8 in UTC; the date in t's own
	// location is kept.
	actual, err := datetime.NewDate(time.Date(2026, 11, 7, 23, 30, 0, 0, loc))
	assert.NoError(t, err)
	assert.Equal(t, "2026-11-07", actual.String())
	assert.Equal(t, datetime.MustParseDate("2026-11-07"), actual, "comparable with ==")
	assert.Equal(t, time.Date(2026, 11, 7, 0, 0, 0, 0, time.UTC), actual.Time())

	_, err = datetime.NewDate(time.Time{})
	assert.ErrorIs(t, err, datetime.ErrInvalidDate)

	assert.Panics(t, func() { datetime.MustNewDate(time.Time{}) })
}

func Test_Today(t *testing.T) {
	today := datetime.Today()

	assert.NoError(t, today.Validate())
	assert.Equal(t, time.Now().Format(time.DateOnly), today.String())
}

func Test_Date_Validate(t *testing.T) {
	assert.NoError(t, datetime.MustParseDate("2026-11-07").Validate())
	assert.ErrorIs(t, datetime.ZeroDate.Validate(), datetime.ErrInvalidDate)
}

func Test_Date_JSON(t *testing.T) {
	type wrapper struct {
		On datetime.Date `json:"on"`
	}

	expected := wrapper{On: datetime.MustParseDate("2026-11-07")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"on":"2026-11-07"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"on":"nope"}`), &actual), datetime.ErrInvalidDate)
}

func Test_Date_UnmarshalJSON(t *testing.T) {
	valid := datetime.MustParseDate("2026-11-07")

	tcs := []struct {
		name     string
		input    string
		start    datetime.Date
		expected datetime.Date
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
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   datetime.ErrInvalidDate,
		},
		{
			name:  "number",
			input: `123`,
			err:   datetime.ErrInvalidDate,
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
	m := map[datetime.Date]string{
		datetime.MustParseDate("2026-11-07"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"2026-11-07":"a"}`, string(data))

	var actual map[datetime.Date]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
