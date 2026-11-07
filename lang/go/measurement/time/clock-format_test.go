package time_test

import (
	"encoding/json"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseClockFormat(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.ClockFormat
		err      error
	}{
		{name: "12-hour", input: "12-hour", expected: time.ClockFormat12Hour},
		{name: "24-hour", input: "24-hour", expected: time.ClockFormat24Hour},
		{name: "uppercase", input: "24-HOUR", expected: time.ClockFormat24Hour},
		{name: "empty", input: "", err: time.ErrInvalidClockFormat},
		{name: "number only", input: "24", err: time.ErrInvalidClockFormat},
		{name: "short form", input: "24h", err: time.ErrInvalidClockFormat},
		{name: "other name", input: "military", err: time.ErrInvalidClockFormat},
		{name: "legacy name", input: "measurement.time.clock-format.24-hour", err: time.ErrInvalidClockFormat},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseClockFormat(tc.input)

			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
				assert.True(tt, actual.IsZero())
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
			assert.NoError(tt, actual.Validate())
		})
	}
}

func Test_ClockFormat_Values(t *testing.T) {
	assert.Equal(t, []time.ClockFormat{time.ClockFormat12Hour, time.ClockFormat24Hour}, time.ClockFormats)
	assert.Equal(t, "12-hour", time.ClockFormat12Hour.String())
	assert.Equal(t, "24-hour", time.ClockFormat24Hour.String())
}

func Test_ClockFormat_Validate(t *testing.T) {
	for _, c := range time.ClockFormats {
		assert.NoError(t, c.Validate(), c.String())
	}

	assert.ErrorIs(t, time.ZeroClockFormat.Validate(), time.ErrInvalidClockFormat)
	assert.ErrorIs(t, time.ClockFormat(3).Validate(), time.ErrInvalidClockFormat)
	assert.Equal(t, "ClockFormat(0)", time.ZeroClockFormat.String())
}

func Test_ClockFormat_JSON(t *testing.T) {
	type wrapper struct {
		Clock time.ClockFormat `json:"clock"`
	}

	expected := wrapper{Clock: time.ClockFormat12Hour}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"clock":"12-hour"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"clock":"24h"}`), &actual), time.ErrInvalidClockFormat)
	assert.Equal(t, expected, actual, "invalid JSON leaves the value unchanged")

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"clock":24}`), &actual), time.ErrInvalidClockFormat)

	_, err = json.Marshal(wrapper{})
	assert.ErrorIs(t, err, time.ErrInvalidClockFormat, "an unset ClockFormat is not encoded")
}

func Test_ClockFormat_Text(t *testing.T) {
	m := map[time.ClockFormat]string{time.ClockFormat24Hour: "a"}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"24-hour":"a"}`, string(data))

	var actual map[time.ClockFormat]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
