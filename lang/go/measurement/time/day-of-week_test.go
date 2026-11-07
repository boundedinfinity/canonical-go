package time_test

import (
	"encoding/json"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func Test_ParseDayOfWeek(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected time.DayOfWeek
		err      error
	}{
		{name: "monday", input: "Monday", expected: time.Monday},
		{name: "tuesday", input: "Tuesday", expected: time.Tuesday},
		{name: "wednesday", input: "Wednesday", expected: time.Wednesday},
		{name: "thursday", input: "Thursday", expected: time.Thursday},
		{name: "friday", input: "Friday", expected: time.Friday},
		{name: "saturday", input: "Saturday", expected: time.Saturday},
		{name: "sunday", input: "Sunday", expected: time.Sunday},
		{name: "lowercase", input: "monday", expected: time.Monday},
		{name: "uppercase", input: "SUNDAY", expected: time.Sunday},
		{name: "abbreviation", input: "Wed", expected: time.Wednesday},
		{name: "lowercase abbreviation", input: "thu", expected: time.Thursday},
		{name: "empty", input: "", err: time.ErrInvalidDayOfWeek},
		{name: "other abbreviation", input: "Tues", err: time.ErrInvalidDayOfWeek},
		{name: "number", input: "1", err: time.ErrInvalidDayOfWeek},
		{name: "surrounding space", input: " Monday", err: time.ErrInvalidDayOfWeek},
		{name: "prefix only", input: "Mo", err: time.ErrInvalidDayOfWeek},
		{name: "not a day", input: "Funday", err: time.ErrInvalidDayOfWeek},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := time.ParseDayOfWeek(tc.input)

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

func Test_DayOfWeek_Values(t *testing.T) {
	tcs := []struct {
		day     time.DayOfWeek
		name    string
		number  int
		weekday gotime.Weekday
	}{
		{day: time.Monday, name: "Monday", number: 1, weekday: gotime.Monday},
		{day: time.Tuesday, name: "Tuesday", number: 2, weekday: gotime.Tuesday},
		{day: time.Wednesday, name: "Wednesday", number: 3, weekday: gotime.Wednesday},
		{day: time.Thursday, name: "Thursday", number: 4, weekday: gotime.Thursday},
		{day: time.Friday, name: "Friday", number: 5, weekday: gotime.Friday},
		{day: time.Saturday, name: "Saturday", number: 6, weekday: gotime.Saturday},
		{day: time.Sunday, name: "Sunday", number: 7, weekday: gotime.Sunday},
	}

	assert.Len(t, time.DaysOfWeek, len(tcs))

	for i, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			assert.Equal(tt, tc.day, time.DaysOfWeek[i], "DaysOfWeek order")
			assert.Equal(tt, tc.name, tc.day.String())
			assert.Equal(tt, tc.number, tc.day.Number())
			assert.Equal(tt, tc.weekday, tc.day.Weekday())
		})
	}
}

func Test_NewDayOfWeek(t *testing.T) {
	// 2026-11-07 is a Saturday.
	assert.Equal(t, time.Saturday, time.NewDayOfWeek(gotime.Date(2026, 11, 7, 12, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.Sunday, time.NewDayOfWeek(gotime.Date(2026, 11, 8, 12, 0, 0, 0, gotime.UTC)))
	assert.Equal(t, time.Monday, time.NewDayOfWeek(gotime.Date(2026, 11, 9, 12, 0, 0, 0, gotime.UTC)))

	// 23:30 on Saturday in UTC-5 is already Sunday in UTC; the day in t's own
	// location is kept.
	loc := gotime.FixedZone("", -5*60*60)
	assert.Equal(t, time.Saturday, time.NewDayOfWeek(gotime.Date(2026, 11, 7, 23, 30, 0, 0, loc)))

	assert.Equal(t, time.Saturday, time.NewDayOfWeek(time.MustParseDate("2026-11-07").Time()), "from a Date")
}

func Test_DayOfWeek_Validate(t *testing.T) {
	for _, d := range time.DaysOfWeek {
		assert.NoError(t, d.Validate(), d.String())
	}

	assert.ErrorIs(t, time.ZeroDayOfWeek.Validate(), time.ErrInvalidDayOfWeek)
	assert.ErrorIs(t, time.DayOfWeek(8).Validate(), time.ErrInvalidDayOfWeek)
	assert.ErrorIs(t, time.DayOfWeek(-1).Validate(), time.ErrInvalidDayOfWeek)
	assert.Equal(t, "DayOfWeek(0)", time.ZeroDayOfWeek.String())
	assert.Equal(t, "DayOfWeek(8)", time.DayOfWeek(8).String())
}

func Test_DayOfWeek_JSON(t *testing.T) {
	type wrapper struct {
		On time.DayOfWeek `json:"on"`
	}

	expected := wrapper{On: time.Wednesday}

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"on":"Wednesday"}`, string(data))

	var actual wrapper
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)

	assert.ErrorIs(t, json.Unmarshal([]byte(`{"on":"Funday"}`), &actual), time.ErrInvalidDayOfWeek)

	_, err = json.Marshal(wrapper{})
	assert.ErrorIs(t, err, time.ErrInvalidDayOfWeek, "an unset DayOfWeek is not encoded")
}

func Test_DayOfWeek_UnmarshalJSON(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		start    time.DayOfWeek
		expected time.DayOfWeek
		err      error
	}{
		{name: "name", input: `"Friday"`, expected: time.Friday},
		{name: "abbreviation", input: `"fri"`, expected: time.Friday},
		{name: "null leaves value unchanged", input: `null`, start: time.Friday, expected: time.Friday},
		{name: "invalid name leaves value unchanged", input: `"Funday"`, start: time.Friday, expected: time.Friday, err: time.ErrInvalidDayOfWeek},
		{name: "empty string", input: `""`, err: time.ErrInvalidDayOfWeek},
		{name: "number", input: `1`, err: time.ErrInvalidDayOfWeek},
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

func Test_DayOfWeek_Text(t *testing.T) {
	m := map[time.DayOfWeek]string{
		time.Monday: "a",
		time.Sunday: "b",
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"Monday":"a","Sunday":"b"}`, string(data))

	var actual map[time.DayOfWeek]string
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, m, actual)
}
