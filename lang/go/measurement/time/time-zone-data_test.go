package time_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	gotime "time"

	"github.com/boundedinfinity/canonical-go/lang/go/measurement/time"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_embedded_time_zones_match_domain_model fails when the embedded copy
// drifts from the source of truth. Run `go generate ./lang/go/measurement/...`
// to refresh it.
func Test_embedded_time_zones_match_domain_model(t *testing.T) {
	embedded, err := os.ReadFile("time-zone.json")
	require.NoError(t, err)

	source, err := os.ReadFile("../../../../domain-model/measurement/time/time-zone.json")
	require.NoError(t, err)

	assert.Equal(t, string(source), string(embedded), "run go generate to refresh time-zone.json")
}

// Test_time_zone_data_round_trips checks that converting the loaded models
// back to DTOs reproduces the data file exactly.
func Test_time_zone_data_round_trips(t *testing.T) {
	dtos := make([]time.TimeZoneDto, 0)

	for _, timeZone := range time.TimeZones.All() {
		dtos = append(dtos, timeZone.ToDto())
	}

	expected, err := os.ReadFile("time-zone.json")
	require.NoError(t, err)

	actual, err := json.Marshal(dtos)
	require.NoError(t, err)

	assert.JSONEq(t, string(expected), string(actual))
}

// Test_time_zone_data_matches_iana checks each zone against the IANA database:
// its Location loads, and in mid-January (standard time for every zone in the
// data, which are all in the northern hemisphere) its abbreviation and offset
// match Abbreviation and UtcOffset.
func Test_time_zone_data_matches_iana(t *testing.T) {
	for _, timeZone := range time.TimeZones.All() {
		t.Run(timeZone.Identifier, func(tt *testing.T) {
			location, err := timeZone.Location()
			require.NoError(tt, err)

			name, offset := gotime.Date(2026, 1, 15, 12, 0, 0, 0, location).Zone()

			sign := '+'
			if offset < 0 {
				sign, offset = '-', -offset
			}

			assert.Equal(tt, timeZone.Abbreviation, name)
			assert.Equal(tt, timeZone.UtcOffset, fmt.Sprintf("%c%02d:%02d", sign, offset/3600, offset%3600/60))
		})
	}
}

func Test_TimeZones_All(t *testing.T) {
	all := time.TimeZones.All()

	assert.Len(t, all, 8)
	assert.Equal(t, "America/Chicago", all[0].Identifier)

	for i := 1; i < len(all); i++ {
		assert.Less(t, all[i-1].Identifier, all[i].Identifier, "sorted by Identifier")
	}

	all[0].Name = "changed"
	assert.Equal(t, "Central Standard Time", time.TimeZones.All()[0].Name, "All returns a copy")
}

func Test_TimeZones_Get(t *testing.T) {
	tcs := []struct {
		term     string
		expected string
		found    bool
	}{
		{term: "America/New_York", expected: "America/New_York", found: true},
		{term: "america/new_york", expected: "America/New_York", found: true},
		{term: "Eastern Standard Time", expected: "America/New_York", found: true},
		{term: "coordinated universal time", expected: "Etc/UTC", found: true},
		{term: "EST"},
		{term: "-05:00"},
		{term: "America/New"},
		{term: ""},
	}

	for _, tc := range tcs {
		t.Run(tc.term, func(tt *testing.T) {
			actual, ok := time.TimeZones.Get(tc.term)

			assert.Equal(tt, tc.found, ok)
			assert.Equal(tt, tc.expected, actual.Identifier)
		})
	}
}

func Test_TimeZones_Find(t *testing.T) {
	tcs := []struct {
		term     string
		expected []string
	}{
		{term: "Asia/Kolkata", expected: []string{"Asia/Kolkata"}},
		{term: "Pacific Standard Time", expected: []string{"America/Los_Angeles"}},
		{term: "EST", expected: []string{"America/New_York"}},
		{term: "ist", expected: []string{"Asia/Kolkata"}},
		{term: "+05:30", expected: []string{"Asia/Kolkata"}},
		{term: "+00:00", expected: []string{"Etc/UTC", "Europe/London"}},
		{term: "XYZ"},
		{term: "Standard"},
	}

	for _, tc := range tcs {
		t.Run(tc.term, func(tt *testing.T) {
			var actual []string

			for _, timeZone := range time.TimeZones.Find(tc.term) {
				actual = append(actual, timeZone.Identifier)
			}

			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_TimeZones_Resolve(t *testing.T) {
	var _ resolver.Resolver[string, time.TimeZone] = time.TimeZones

	actual, err := time.TimeZones.Resolve("europe/london")
	assert.NoError(t, err)
	assert.Equal(t, "Europe/London", actual.Identifier)

	_, err = time.TimeZones.Resolve("Greenwich Mean Time")
	assert.ErrorIs(t, err, time.ErrNotFound, "Resolve uses the natural key only")

	_, err = resolver.One[string, time.TimeZone](time.TimeZones, "Mars/Olympus_Mons")
	assert.ErrorIs(t, err, resolver.ErrResolve)
	assert.ErrorIs(t, err, time.ErrNotFound)
}
