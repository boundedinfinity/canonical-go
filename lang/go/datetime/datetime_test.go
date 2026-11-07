package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/boundedinfinity/canonical-go/lang/go/datetime"
	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "utc",
			input:    "2026-11-07T14:30:00Z",
			expected: "2026-11-07T14:30:00Z",
		},
		{
			name:     "offset is preserved",
			input:    "2026-11-07T09:30:00-05:00",
			expected: "2026-11-07T09:30:00-05:00",
		},
		{
			name:     "fractional seconds",
			input:    "2026-11-07T09:30:00.123-05:00",
			expected: "2026-11-07T09:30:00.123-05:00",
		},
		{
			name:     "trailing zeros in fractional seconds are dropped",
			input:    "2026-11-07T14:30:00.500000Z",
			expected: "2026-11-07T14:30:00.5Z",
		},
		{
			name:  "empty",
			input: "",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "missing offset",
			input: "2026-11-07T14:30:00",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "date only",
			input: "2026-11-07",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "space separator",
			input: "2026-11-07 14:30:00Z",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "out of range month",
			input: "2026-13-07T14:30:00Z",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "zero value",
			input: "0001-01-01T00:00:00Z",
			err:   datetime.ErrInvalid,
		},
		{
			name:  "not a datetime",
			input: "not-a-datetime",
			err:   datetime.ErrInvalid,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := datetime.Parse(tc.input)

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

func Test_New(t *testing.T) {
	loc := time.FixedZone("", -5*60*60)
	input := time.Date(2026, 11, 7, 9, 30, 0, 0, loc)

	actual, err := datetime.New(input)
	assert.NoError(t, err)
	assert.Equal(t, "2026-11-07T09:30:00-05:00", actual.String())
	assert.True(t, input.Equal(actual.Time()))

	_, err = datetime.New(time.Time{})
	assert.ErrorIs(t, err, datetime.ErrInvalid)

	assert.Panics(t, func() { datetime.MustNew(time.Time{}) })
}

func Test_Now(t *testing.T) {
	assert.NoError(t, datetime.Now().Validate())
}

func Test_Equal(t *testing.T) {
	a := datetime.MustParse("2026-11-07T14:30:00Z")
	b := datetime.MustParse("2026-11-07T09:30:00-05:00")
	c := datetime.MustParse("2026-11-07T14:30:01Z")

	assert.True(t, a.Equal(b))
	assert.False(t, a.Equal(c))
}

func Test_Validate(t *testing.T) {
	assert.NoError(t, datetime.MustParse("2026-11-07T14:30:00Z").Validate())
	assert.ErrorIs(t, datetime.Zero.Validate(), datetime.ErrInvalid)
}

func Test_JSON(t *testing.T) {
	type wrapper struct {
		At datetime.DateTime `json:"at"`
	}

	expected := wrapper{At: datetime.MustParse("2026-11-07T09:30:00.123-05:00")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"at":"2026-11-07T09:30:00.123-05:00"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.True(t, expected.At.Equal(actual.At))
	assert.Equal(t, expected.At.String(), actual.At.String())

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"at":"nope"}`), &actual), datetime.ErrInvalid)
}

func Test_UnmarshalJSON(t *testing.T) {
	valid := datetime.MustParse("2026-11-07T14:30:00Z")

	tcs := []struct {
		name     string
		input    string
		start    datetime.DateTime
		expected datetime.DateTime
		err      error
	}{
		{
			name:     "string",
			input:    `"2026-11-07T14:30:00Z"`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:  "invalid datetime string",
			input: `"nope"`,
			err:   datetime.ErrInvalid,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   datetime.ErrInvalid,
		},
		{
			name:  "number",
			input: `123`,
			err:   datetime.ErrInvalid,
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
			assert.True(tt, tc.expected.Equal(actual))
		})
	}
}

func Test_Text(t *testing.T) {
	m := map[datetime.DateTime]string{
		datetime.MustParse("2026-11-07T14:30:00Z"): "a",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"2026-11-07T14:30:00Z":"a"}`, string(data))

	var actual map[datetime.DateTime]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Len(t, actual, 1)

	for k, v := range actual {
		assert.Equal(t, "2026-11-07T14:30:00Z", k.String())
		assert.Equal(t, "a", v)
	}
}
