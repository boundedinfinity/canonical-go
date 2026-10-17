package iso3166_test

import (
	"encoding/json"
	"strings"
	"testing"

	iso3166 "github.com/boundedinfinity/canonical-go/lang/go/specification/iso-3166"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_NewCountryCode(t *testing.T) {
	tcs := []struct {
		name      string
		alpha2    string
		alpha3    string
		numeric   string
		cname     string
		reference string
		errs      []error
	}{
		{name: "valid", alpha2: "US", alpha3: "USA", numeric: "840", cname: "United States of America", reference: "https://en.wikipedia.org/wiki/United_States"},
		{name: "no reference", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: "Afghanistan"},
		{name: "http reference", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: "Afghanistan", reference: "http://example.com"},
		{name: "name max", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: strings.Repeat("a", 100)},
		{name: "name max counts runes", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: strings.Repeat("é", 100)},
		{name: "alpha2 lowercase", alpha2: "us", alpha3: "USA", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha2}},
		{name: "alpha2 too short", alpha2: "U", alpha3: "USA", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha2}},
		{name: "alpha2 too long", alpha2: "USA", alpha3: "USA", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha2}},
		{name: "alpha2 digit", alpha2: "U1", alpha3: "USA", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha2}},
		{name: "alpha3 lowercase", alpha2: "US", alpha3: "usa", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha3}},
		{name: "alpha3 too short", alpha2: "US", alpha3: "US", numeric: "840", cname: "n", errs: []error{iso3166.ErrAlpha3}},
		{name: "numeric without leading zeros", alpha2: "AF", alpha3: "AFG", numeric: "4", cname: "n", errs: []error{iso3166.ErrNumeric}},
		{name: "numeric letters", alpha2: "AF", alpha3: "AFG", numeric: "0a4", cname: "n", errs: []error{iso3166.ErrNumeric}},
		{name: "numeric too long", alpha2: "AF", alpha3: "AFG", numeric: "0004", cname: "n", errs: []error{iso3166.ErrNumeric}},
		{name: "name empty", alpha2: "AF", alpha3: "AFG", numeric: "004", errs: []error{iso3166.ErrName}},
		{name: "name too long", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: strings.Repeat("a", 101), errs: []error{iso3166.ErrName}},
		{name: "relative reference", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: "n", reference: "/wiki/Afghanistan", errs: []error{iso3166.ErrReference}},
		{name: "ftp reference", alpha2: "AF", alpha3: "AFG", numeric: "004", cname: "n", reference: "ftp://example.com", errs: []error{iso3166.ErrReference}},
		{
			name: "reports every failure",
			errs: []error{iso3166.ErrAlpha2, iso3166.ErrAlpha3, iso3166.ErrNumeric, iso3166.ErrName},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := iso3166.NewCountryCode(tc.alpha2, tc.alpha3, tc.numeric, tc.cname, true, tc.reference)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, iso3166.ErrInvalidCountryCode)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, iso3166.CountryCode{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.alpha2, actual.Alpha2)
			assert.Equal(tt, tc.alpha3, actual.Alpha3)
			assert.Equal(tt, tc.numeric, actual.Numeric)
			assert.Equal(tt, tc.cname, actual.Name)
			assert.True(tt, actual.Independent)
			assert.Equal(tt, optioner.OfZero(tc.reference), actual.Reference)
		})
	}
}

func Test_CountryCode_JSON(t *testing.T) {
	tcs := []struct {
		name     string
		input    iso3166.CountryCode
		expected string
	}{
		{
			name:     "absent reference is omitted and false is written",
			input:    iso3166.CountryCode{Alpha2: "GI", Alpha3: "GIB", Numeric: "292", Name: "Gibraltar"},
			expected: `{"alpha-2":"GI","alpha-3":"GIB","numeric":"292","name":"Gibraltar","independent":false}`,
		},
		{
			name: "all fields",
			input: iso3166.CountryCode{
				Alpha2:      "US",
				Alpha3:      "USA",
				Numeric:     "840",
				Name:        "United States of America",
				Independent: true,
				Reference:   optioner.Some("https://en.wikipedia.org/wiki/United_States"),
			},
			expected: `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"United States of America","independent":true,` +
				`"reference":"https://en.wikipedia.org/wiki/United_States"}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual iso3166.CountryCode
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_CountryCode_UnmarshalJSON_invalid(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		errs  []error
	}{
		{name: "missing independent", input: `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"n"}`, errs: []error{iso3166.ErrIndependent}},
		{name: "missing alpha-2", input: `{"alpha-3":"USA","numeric":"840","name":"n","independent":true}`, errs: []error{iso3166.ErrAlpha2}},
		{name: "numeric as number", input: `{"alpha-2":"US","alpha-3":"USA","numeric":840,"name":"n","independent":true}`},
		{name: "bad reference", input: `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"n","independent":true,"reference":"nope"}`, errs: []error{iso3166.ErrReference}},
		{
			name:  "reports every failure",
			input: `{"alpha-2":"us","numeric":"840","name":"n"}`,
			errs:  []error{iso3166.ErrIndependent, iso3166.ErrAlpha2, iso3166.ErrAlpha3},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual iso3166.CountryCode
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.Error(tt, err)
			for _, e := range tc.errs {
				assert.ErrorIs(tt, err, iso3166.ErrInvalidCountryCode)
				assert.ErrorIs(tt, err, e)
			}
			assert.Equal(tt, iso3166.CountryCode{}, actual)
		})
	}
}

func Test_CountryCode_UnmarshalJSON_null(t *testing.T) {
	expected := iso3166.MustNewCountryCode("US", "USA", "840", "United States of America", true, "")
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
