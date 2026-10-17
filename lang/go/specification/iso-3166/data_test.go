package iso3166_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	iso3166 "github.com/boundedinfinity/canonical-go/lang/go/specification/iso-3166"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_embedded_data_matches_domain_model fails when the embedded copies drift
// from the source of truth. Run `go generate ./lang/go/specification/...` to
// refresh them.
func Test_embedded_data_matches_domain_model(t *testing.T) {
	for _, name := range []string{"country-code.json", "subdivision.json"} {
		t.Run(name, func(tt *testing.T) {
			embedded, err := os.ReadFile(name)
			require.NoError(tt, err)

			source, err := os.ReadFile("../../../../domain-model/specification/iso-3166/" + name)
			require.NoError(tt, err)

			assert.Equal(tt, string(source), string(embedded), "run go generate to refresh %s", name)
		})
	}
}

// Test_data_round_trips checks that converting the loaded models back to DTOs
// reproduces the data files exactly.
func Test_data_round_trips(t *testing.T) {
	countries := make([]iso3166.CountryCodeDto, 0)
	for _, country := range iso3166.CountryCodes.All() {
		countries = append(countries, country.ToDto())
	}

	subdivisions := make([]iso3166.SubdivisionDto, 0)
	for _, subdivision := range iso3166.Subdivisions.All() {
		subdivisions = append(subdivisions, subdivision.ToDto())
	}

	for name, dtos := range map[string]any{"country-code.json": countries, "subdivision.json": subdivisions} {
		t.Run(name, func(tt *testing.T) {
			expected, err := os.ReadFile(name)
			require.NoError(tt, err)

			actual, err := json.Marshal(dtos)
			require.NoError(tt, err)

			assert.JSONEq(tt, string(expected), string(actual))
		})
	}
}

func Test_CountryCodes_All(t *testing.T) {
	all := iso3166.CountryCodes.All()

	assert.Len(t, all, 249)
	assert.Equal(t, "AD", all[0].Alpha2)

	all[0].Name = "changed"
	assert.Equal(t, "Andorra", iso3166.CountryCodes.All()[0].Name, "All returns a copy")
}

func Test_CountryCodes_Get(t *testing.T) {
	tcs := []struct {
		term     string
		expected string
	}{
		{term: "United States of America", expected: "US"},
		{term: "united states of america", expected: "US"},
		{term: "US", expected: "US"},
		{term: "us", expected: "US"},
		{term: "USA", expected: "US"},
		{term: "840", expected: "US"},
		{term: "004", expected: "AF"},
		{term: "4", expected: "AF"},
		{term: "Afghanistan", expected: "AF"},
		{term: "Narnia"},
		{term: ""},
	}

	for _, tc := range tcs {
		t.Run(tc.term, func(tt *testing.T) {
			actual, ok := iso3166.CountryCodes.Get(tc.term)

			if tc.expected == "" {
				assert.False(tt, ok)
				assert.Equal(tt, iso3166.CountryCode{}, actual)
				return
			}

			assert.True(tt, ok)
			assert.Equal(tt, tc.expected, actual.Alpha2)
		})
	}
}

func Test_CountryCodes_Get_every_country(t *testing.T) {
	for _, country := range iso3166.CountryCodes.All() {
		for _, term := range []string{country.Name, country.Alpha2, country.Alpha3, country.Numeric} {
			actual, ok := iso3166.CountryCodes.Get(term)

			assert.True(t, ok, term)
			assert.Equal(t, country, actual, term)
		}
	}
}

func Test_CountryCodes_Find(t *testing.T) {
	alpha2s := func(countries []iso3166.CountryCode) []string {
		var codes []string
		for _, country := range countries {
			codes = append(codes, country.Alpha2)
		}
		return codes
	}

	assert.Equal(t, []string{"US"}, alpha2s(iso3166.CountryCodes.Find("usa")), "exact match")
	assert.Equal(t, []string{"GN", "GQ", "GW", "PG"}, alpha2s(iso3166.CountryCodes.Find("guin")), "substring of names")
	assert.Equal(t, []string{"GN"}, alpha2s(iso3166.CountryCodes.Find("guinea")), "an exact name wins over substrings")
	assert.Empty(t, iso3166.CountryCodes.Find("narnia"))
}

func Test_CountryCodes_Resolve(t *testing.T) {
	var countries resolver.Resolver[string, iso3166.CountryCode] = iso3166.CountryCodes

	actual, err := countries.Resolve("US")
	assert.NoError(t, err)
	assert.Equal(t, "United States of America", actual.Name)

	_, err = countries.Resolve("ZZ")
	assert.ErrorIs(t, err, iso3166.ErrNotFound)
}

func Test_Subdivisions_All(t *testing.T) {
	all := iso3166.Subdivisions.All()

	assert.Len(t, all, 3379)
	assert.Equal(t, "AD-02", all[0].Code)

	for _, subdivision := range all {
		country, err := iso3166.CountryCodes.Resolve(subdivision.Country.Alpha2)
		assert.NoError(t, err)
		assert.Equal(t, country, subdivision.Country)
	}
}

func Test_Subdivisions_Get(t *testing.T) {
	tcs := []struct {
		term     string
		expected string
	}{
		{term: "US-AL", expected: "US-AL"},
		{term: "us-al", expected: "US-AL"},
		{term: "USAL", expected: "US-AL"},
		{term: "usal", expected: "US-AL"},
		{term: "Alabama"},
		{term: "US-ZZ"},
	}

	for _, tc := range tcs {
		t.Run(tc.term, func(tt *testing.T) {
			actual, ok := iso3166.Subdivisions.Get(tc.term)

			if tc.expected == "" {
				assert.False(tt, ok)
				return
			}

			assert.True(tt, ok)
			assert.Equal(tt, tc.expected, actual.Code)
			assert.Equal(tt, "Alabama", actual.Name)
		})
	}
}

func Test_Subdivisions_FindByName(t *testing.T) {
	codes := func(subdivisions []iso3166.Subdivision) []string {
		var codes []string
		for _, subdivision := range subdivisions {
			codes = append(codes, subdivision.Code)
		}
		return codes
	}

	assert.Equal(t, []string{"US-AL"}, codes(iso3166.Subdivisions.FindByName("alabama")))
	assert.Equal(t,
		[]string{"FJ-W", "GH-WP", "GM-W", "PG-WPD", "RW-04", "SB-WE", "UG-W", "ZM-01"},
		codes(iso3166.Subdivisions.FindByName("Western")),
		"names are not unique",
	)
	assert.Empty(t, iso3166.Subdivisions.FindByName("Narnia"))
}

func Test_Subdivisions_ByCountry(t *testing.T) {
	us, _ := iso3166.CountryCodes.Get("US")
	gi, _ := iso3166.CountryCodes.Get("GI")

	subdivisions := iso3166.Subdivisions.ByCountry(us)
	assert.Len(t, subdivisions, 57)
	for _, subdivision := range subdivisions {
		assert.Equal(t, "US", subdivision.Country.Alpha2)
	}

	assert.Empty(t, iso3166.Subdivisions.ByCountry(gi))
}

func Test_Subdivisions_Resolve(t *testing.T) {
	var subdivisions resolver.Resolver[string, iso3166.Subdivision] = iso3166.Subdivisions

	actual, err := subdivisions.Resolve("US-AL")
	assert.NoError(t, err)
	assert.Equal(t, "Alabama", actual.Name)

	_, err = subdivisions.Resolve("US-ZZ")
	assert.ErrorIs(t, err, iso3166.ErrNotFound)
}

func Test_SubdivisionDto_ToModel_with_data(t *testing.T) {
	dto := iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "County"}

	actual, err := dto.ToModel(iso3166.CountryCodes, iso3166.Subdivisions)
	assert.NoError(t, err)
	assert.Equal(t, "United States of America", actual.Country.Name)
}
