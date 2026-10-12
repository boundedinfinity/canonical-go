package id_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
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
			name:     "valid",
			input:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			expected: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		},
		{
			name:     "uppercase is normalized",
			input:    "F47AC10B-58CC-4372-A567-0E02B2C3D479",
			expected: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
		},
		{
			name:  "empty",
			input: "",
			err:   id.ErrInvalid,
		},
		{
			name:  "not a uuid",
			input: "not-a-uuid",
			err:   id.ErrInvalid,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := id.Parse(tc.input)

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
	a, b := id.New(), id.New()

	assert.False(t, a.IsZero())
	assert.NotEqual(t, a, b)
	assert.Equal(t, byte('4'), a.String()[14], "version nibble")
}

func Test_Validate(t *testing.T) {
	assert.NoError(t, id.New().Validate())
	assert.ErrorIs(t, id.Zero.Validate(), id.ErrInvalid)
}

func Test_JSON(t *testing.T) {
	type wrapper struct {
		Id id.ID `json:"id"`
	}

	expected := wrapper{Id: id.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":"f47ac10b-58cc-4372-a567-0e02b2c3d479"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"id":"nope"}`), &actual), id.ErrInvalid)
}

func Test_UnmarshalJSON(t *testing.T) {
	valid := id.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")

	tcs := []struct {
		name     string
		input    string
		start    id.ID
		expected id.ID
		err      error
	}{
		{
			name:     "string",
			input:    `"f47ac10b-58cc-4372-a567-0e02b2c3d479"`,
			expected: valid,
		},
		{
			name:     "null leaves value unchanged",
			input:    `null`,
			start:    valid,
			expected: valid,
		},
		{
			name:  "invalid uuid string",
			input: `"nope"`,
			err:   id.ErrInvalid,
		},
		{
			name:  "empty string",
			input: `""`,
			err:   id.ErrInvalid,
		},
		{
			name:  "number",
			input: `123`,
			err:   id.ErrInvalid,
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

func Test_MarshalJSON(t *testing.T) {
	data, err := json.Marshal(id.MustParse("F47AC10B-58CC-4372-A567-0E02B2C3D479"))

	assert.NoError(t, err)
	assert.Equal(t, `"f47ac10b-58cc-4372-a567-0e02b2c3d479"`, string(data))
}
