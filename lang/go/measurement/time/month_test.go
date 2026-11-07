package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseMonth(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.Month
		err      error
	}{
		{name: "january", input: "January", expected: time.January},
		{name: "may", input: "May", expected: time.May},
		{name: "september", input: "September", expected: time.September},
		{name: "december", input: "December", expected: time.December},
		{name: "lowercase", input: "january", expected: time.January},
		{name: "uppercase", input: "DECEMBER", expected: time.December},
		{name: "abbreviation", input: "Sep", expected: time.September},
		{name: "lowercase abbreviation", input: "oct", expected: time.October},
		{name: "empty", input: "", err: time.ErrInvalidMonth},
		{name: "other abbreviation", input: "Sept", err: time.ErrInvalidMonth},
		{name: "number", input: "1", err: time.ErrInvalidMonth},
		{name: "surrounding space", input: " January", err: time.ErrInvalidMonth},
		{name: "prefix only", input: "Ja", err: time.ErrInvalidMonth},
		{name: "not a month", input: "Smarch", err: time.ErrInvalidMonth},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseMonth(tc.input)

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

func Test_Month_Values(t *testing.T) {
	tcs := []struct {
		month  time.Month
		name   string
		number int
		native gotime.Month
	}{
		{month: time.January, name: "January", number: 1, native: gotime.January},
		{month: time.February, name: "February", number: 2, native: gotime.February},
		{month: time.March, name: "March", number: 3, native: gotime.March},
		{month: time.April, name: "April", number: 4, native: gotime.April},
		{month: time.May, name: "May", number: 5, native: gotime.May},
		{month: time.June, name: "June", number: 6, native: gotime.June},
		{month: time.July, name: "July", number: 7, native: gotime.July},
		{month: time.August, name: "August", number: 8, native: gotime.August},
		{month: time.September, name: "September", number: 9, native: gotime.September},
		{month: time.October, name: "October", number: 10, native: gotime.October},
		{month: time.November, name: "November", number: 11, native: gotime.November},
		{month: time.December, name: "December", number: 12, native: gotime.December},
	}

	assert.Len(t, time.MonthsOfYear, len(tcs))

	for i, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.month, time.MonthsOfYear[i], "MonthsOfYear order")
			assert.Equal(tt, tc.name, tc.month.String())
			assert.Equal(tt, tc.number, tc.month.Number())
			assert.Equal(tt, tc.native, tc.month.Month())

			parsed, err := time.ParseMonth(tc.name[:3])
			assert.NoError(tt, err, "abbreviation")
			assert.Equal(tt, tc.month, parsed)
		})
	}
}

func Test_NewMonth(t *testing.T) {
	assert.Equal(t, time.November, time.NewMonth(gotime.Date(2026, 11, 7, 12, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.January, time.NewMonth(gotime.Date(2027, 1, 1, 0, 0, 0, 0, gotime.UTC)))

	// 23:30 on Nov 30 in UTC-5 is already Dec 1 in UTC; the month in t's own
	// location is kept.
	loc := gotime.FixedZone("", -5*60*60)
	assert.Equal(t, time.November, time.NewMonth(gotime.Date(2026, 11, 30, 23, 30, 0, 0, loc)))

	assert.Equal(t, time.November, time.NewMonth(time.MustParseDate("2026-11-07").Time()), "from a Date")
}

func Test_Month_Validate(t *testing.T) {
	for _, m := range time.MonthsOfYear {
		assert.NoError(t, m.Validate(), m.String())
	}

	assert.ErrorIs(t, time.ZeroMonth.Validate(), time.ErrInvalidMonth)
	assert.ErrorIs(t, time.Month(13).Validate(), time.ErrInvalidMonth)
	assert.ErrorIs(t, time.Month(-1).Validate(), time.ErrInvalidMonth)
	assert.Equal(t, "Month(0)", time.ZeroMonth.String())
	assert.Equal(t, "Month(13)", time.Month(13).String())
}

func Test_Month_JSON(t *testing.T) {
	type wrapper struct {
		In time.Month `json:"in"`
	}

	expected := wrapper{In: time.March}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"in":"March"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"in":"Smarch"}`), &actual), time.ErrInvalidMonth)

	_, err = json.Marshal(wrapper{})
	assert.ErrorIs(t, err, time.ErrInvalidMonth, "an unset Month is not encoded")
}

func Test_Month_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		start    time.Month
		expected time.Month
		err      error
	}{
		{name: "name", input: `"April"`, expected: time.April},
		{name: "abbreviation", input: `"apr"`, expected: time.April},
		{name: "null leaves value unchanged", input: `null`, start: time.April, expected: time.April},
		{name: "invalid name leaves value unchanged", input: `"Smarch"`, start: time.April, expected: time.April, err: time.ErrInvalidMonth},
		{name: "empty string", input: `""`, err: time.ErrInvalidMonth},
		{name: "number", input: `4`, err: time.ErrInvalidMonth},
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

func Test_Month_Text(t *testing.T) {
	m := map[time.Month]string{
		time.January:  "a",
		time.December: "b",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"January":"a","December":"b"}`, string(data))

	var actual map[time.Month]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
