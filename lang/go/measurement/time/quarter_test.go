package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseQuarter(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.Quarter
		err      error
	}{
		{name: "q1", input: "Q1", expected: time.Q1},
		{name: "q2", input: "Q2", expected: time.Q2},
		{name: "q3", input: "Q3", expected: time.Q3},
		{name: "q4", input: "Q4", expected: time.Q4},
		{name: "lowercase", input: "q3", expected: time.Q3},
		{name: "empty", input: "", err: time.ErrInvalidQuarter},
		{name: "number only", input: "1", err: time.ErrInvalidQuarter},
		{name: "out of range", input: "Q5", err: time.ErrInvalidQuarter},
		{name: "zero", input: "Q0", err: time.ErrInvalidQuarter},
		{name: "word", input: "first", err: time.ErrInvalidQuarter},
		{name: "surrounding space", input: " Q1", err: time.ErrInvalidQuarter},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseQuarter(tc.input)

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

func Test_Quarter_Values(t *testing.T) {
	tcs := []struct {
		quarter time.Quarter
		name    string
		number  int
		months  []time.Month
	}{
		{quarter: time.Q1, name: "Q1", number: 1, months: []time.Month{time.January, time.February, time.March}},
		{quarter: time.Q2, name: "Q2", number: 2, months: []time.Month{time.April, time.May, time.June}},
		{quarter: time.Q3, name: "Q3", number: 3, months: []time.Month{time.July, time.August, time.September}},
		{quarter: time.Q4, name: "Q4", number: 4, months: []time.Month{time.October, time.November, time.December}},
	}

	assert.Len(t, time.QuartersOfYear, len(tcs))

	for i, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.quarter, time.QuartersOfYear[i], "QuartersOfYear order")
			assert.Equal(tt, tc.name, tc.quarter.String())
			assert.Equal(tt, tc.number, tc.quarter.Number())
			assert.Equal(tt, tc.months, tc.quarter.Months())

			for _, m := range tc.months {
				actual, err := time.QuarterOfMonth(m)
				assert.NoError(tt, err)
				assert.Equal(tt, tc.quarter, actual, "QuarterOfMonth(%s)", m)
			}
		})
	}
}

func Test_NewQuarter(t *testing.T) {
	assert.Equal(t, time.Q4, time.NewQuarter(gotime.Date(2026, 11, 7, 12, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.Q1, time.NewQuarter(gotime.Date(2027, 1, 1, 0, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.Q2, time.NewQuarter(gotime.Date(2026, 6, 30, 12, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.Q3, time.NewQuarter(gotime.Date(2026, 7, 1, 12, 0, 0, 0, gotime.UTC)))

	// 23:30 on Mar 31 in UTC-5 is already Apr 1 in UTC; the quarter in t's own
	// location is kept.
	loc := gotime.FixedZone("", -5*60*60)
	assert.Equal(t, time.Q1, time.NewQuarter(gotime.Date(2026, 3, 31, 23, 30, 0, 0, loc)))
}

func Test_QuarterOfMonth_Invalid(t *testing.T) {
	_, err := time.QuarterOfMonth(time.ZeroMonth)
	assert.ErrorIs(t, err, time.ErrInvalidMonth)

	_, err = time.QuarterOfMonth(time.Month(13))
	assert.ErrorIs(t, err, time.ErrInvalidMonth)
}

func Test_Quarter_Validate(t *testing.T) {
	for _, q := range time.QuartersOfYear {
		assert.NoError(t, q.Validate(), q.String())
	}

	assert.ErrorIs(t, time.ZeroQuarter.Validate(), time.ErrInvalidQuarter)
	assert.ErrorIs(t, time.Quarter(5).Validate(), time.ErrInvalidQuarter)
	assert.ErrorIs(t, time.Quarter(-1).Validate(), time.ErrInvalidQuarter)
	assert.Equal(t, "Quarter(0)", time.ZeroQuarter.String())
	assert.Equal(t, "Quarter(5)", time.Quarter(5).String())
}

func Test_Quarter_JSON(t *testing.T) {
	type wrapper struct {
		In time.Quarter `json:"in"`
	}

	expected := wrapper{In: time.Q3}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"in":"Q3"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"in":"Q5"}`), &actual), time.ErrInvalidQuarter)

	_, err = json.Marshal(wrapper{})
	assert.ErrorIs(t, err, time.ErrInvalidQuarter, "an unset Quarter is not encoded")
}

func Test_Quarter_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		start    time.Quarter
		expected time.Quarter
		err      error
	}{
		{name: "name", input: `"Q2"`, expected: time.Q2},
		{name: "lowercase", input: `"q2"`, expected: time.Q2},
		{name: "null leaves value unchanged", input: `null`, start: time.Q2, expected: time.Q2},
		{name: "invalid name leaves value unchanged", input: `"Q5"`, start: time.Q2, expected: time.Q2, err: time.ErrInvalidQuarter},
		{name: "empty string", input: `""`, err: time.ErrInvalidQuarter},
		{name: "number", input: `2`, err: time.ErrInvalidQuarter},
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

func Test_Quarter_Text(t *testing.T) {
	m := map[time.Quarter]string{
		time.Q1: "a",
		time.Q4: "b",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"Q1":"a","Q4":"b"}`, string(data))

	var actual map[time.Quarter]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
