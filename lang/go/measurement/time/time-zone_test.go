package time_test

import (
	"encoding/json"
	"strings"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/stretchr/testify/assert"
)

func validTimeZone() time.TimeZone {
	return time.TimeZone{
		Identifier:   "America/New_York",
		Name:         "Eastern Standard Time",
		Abbreviation: "EST",
		UtcOffset:    "-05:00",
	}
}

func Test_TimeZone_Validate(t *testing.T) {
	tcs := []struct {
		name   string
		modify func(*time.TimeZone)
		errs   []error
	}{
		{name: "valid", modify: func(*time.TimeZone) {}},
		{name: "single part identifier", modify: func(tz *time.TimeZone) { tz.Identifier = "UTC" }},
		{name: "three part identifier", modify: func(tz *time.TimeZone) { tz.Identifier = "America/Argentina/Buenos_Aires" }},
		{name: "identifier with sign and digits", modify: func(tz *time.TimeZone) { tz.Identifier = "Etc/GMT+5" }},
		{name: "positive offset with minutes", modify: func(tz *time.TimeZone) { tz.UtcOffset = "+05:45" }},
		{name: "largest offset", modify: func(tz *time.TimeZone) { tz.UtcOffset = "+14:00" }},
		{name: "six letter abbreviation", modify: func(tz *time.TimeZone) { tz.Abbreviation = "ABCDEF" }},
		{
			name:   "empty identifier",
			modify: func(tz *time.TimeZone) { tz.Identifier = "" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrIdentifier},
		},
		{
			name:   "identifier with space",
			modify: func(tz *time.TimeZone) { tz.Identifier = "America/New York" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrIdentifier},
		},
		{
			name:   "identifier with empty part",
			modify: func(tz *time.TimeZone) { tz.Identifier = "America//New_York" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrIdentifier},
		},
		{
			name:   "identifier with trailing slash",
			modify: func(tz *time.TimeZone) { tz.Identifier = "America/" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrIdentifier},
		},
		{
			name:   "identifier too long",
			modify: func(tz *time.TimeZone) { tz.Identifier = strings.Repeat("a", 101) },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrIdentifier},
		},
		{
			name:   "empty name",
			modify: func(tz *time.TimeZone) { tz.Name = "" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrName},
		},
		{
			name:   "name too long",
			modify: func(tz *time.TimeZone) { tz.Name = strings.Repeat("a", 101) },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrName},
		},
		{
			name:   "lowercase abbreviation",
			modify: func(tz *time.TimeZone) { tz.Abbreviation = "est" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrAbbreviation},
		},
		{
			name:   "abbreviation too short",
			modify: func(tz *time.TimeZone) { tz.Abbreviation = "E" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrAbbreviation},
		},
		{
			name:   "abbreviation too long",
			modify: func(tz *time.TimeZone) { tz.Abbreviation = "ABCDEFG" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrAbbreviation},
		},
		{
			name:   "numeric abbreviation",
			modify: func(tz *time.TimeZone) { tz.Abbreviation = "+05" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrAbbreviation},
		},
		{
			name:   "offset without sign",
			modify: func(tz *time.TimeZone) { tz.UtcOffset = "05:00" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrUtcOffset},
		},
		{
			name:   "offset hours out of range",
			modify: func(tz *time.TimeZone) { tz.UtcOffset = "+15:00" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrUtcOffset},
		},
		{
			name:   "offset minutes not a quarter hour",
			modify: func(tz *time.TimeZone) { tz.UtcOffset = "+05:20" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrUtcOffset},
		},
		{
			name:   "offset without colon",
			modify: func(tz *time.TimeZone) { tz.UtcOffset = "-0500" },
			errs:   []error{time.ErrInvalidTimeZone, time.ErrUtcOffset},
		},
		{
			name:   "every failure is reported",
			modify: func(tz *time.TimeZone) { *tz = time.TimeZone{} },
			errs:   []error{time.ErrIdentifier, time.ErrName, time.ErrAbbreviation, time.ErrUtcOffset},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			tz := validTimeZone()
			tc.modify(&tz)
			err := tz.Validate()

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

func Test_TimeZone_Location(t *testing.T) {
	location, err := validTimeZone().Location()
	assert.NoError(t, err)
	assert.Equal(t, "America/New_York", location.String())

	winter := gotime.Date(2026, 1, 15, 12, 0, 0, 0, location)
	name, offset := winter.Zone()
	assert.Equal(t, "EST", name)
	assert.Equal(t, -5*60*60, offset)

	summer := gotime.Date(2026, 7, 15, 12, 0, 0, 0, location)
	name, offset = summer.Zone()
	assert.Equal(t, "EDT", name, "Location applies daylight saving time")
	assert.Equal(t, -4*60*60, offset)

	unknown := validTimeZone()
	unknown.Identifier = "Mars/Olympus_Mons"
	assert.NoError(t, unknown.Validate(), "well formed")

	_, err = unknown.Location()
	assert.ErrorIs(t, err, time.ErrInvalidTimeZone)
	assert.ErrorIs(t, err, time.ErrIdentifier)
}

func Test_NewTimeZone(t *testing.T) {
	actual, err := time.NewTimeZone("Asia/Kolkata", "India Standard Time", "IST", "+05:30")
	assert.NoError(t, err)
	assert.Equal(t, "Asia/Kolkata", actual.Identifier)

	_, err = time.NewTimeZone("Asia/Kolkata", "India Standard Time", "IST", "5:30")
	assert.ErrorIs(t, err, time.ErrUtcOffset)

	assert.Panics(t, func() { time.MustNewTimeZone("", "", "", "") })
}

func Test_TimeZone_JSON(t *testing.T) {
	expected := validTimeZone()

	data, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"identifier":"America/New_York","name":"Eastern Standard Time","abbreviation":"EST","utc-offset":"-05:00"}`, string(data))

	var actual time.TimeZone
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func Test_TimeZone_UnmarshalJSON(t *testing.T) {
	valid := validTimeZone()

	tcs := []struct {
		name     string
		input    string
		start    time.TimeZone
		expected time.TimeZone
		err      error
	}{
		{
			name:     "valid",
			input:    `{"identifier":"America/New_York","name":"Eastern Standard Time","abbreviation":"EST","utc-offset":"-05:00"}`,
			expected: valid,
		},
		{name: "null leaves value unchanged", input: `null`, start: valid, expected: valid},
		{
			name:     "invalid leaves value unchanged",
			input:    `{"identifier":"America/New_York","name":"Eastern Standard Time","abbreviation":"EST","utc-offset":"-5"}`,
			start:    valid,
			expected: valid,
			err:      time.ErrUtcOffset,
		},
		{
			name:  "missing fields",
			input: `{"identifier":"America/New_York"}`,
			err:   time.ErrName,
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

			assert.Equal(tt, tc.expected, actual)
		})
	}
}
