package iso3166_test

import (
	"encoding/json"
	"strings"
	"testing"

	iso3166 "github.com/boundedinfinity/canonical-go/lang/go/specification/iso-3166"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

var (
	us = iso3166.MustNewCountryCode("US", "USA", "840", "United States of America", true, "")
	ca = iso3166.MustNewCountryCode("CA", "CAN", "124", "Canada", true, "")
)

func Test_NewSubdivision(t *testing.T) {
	tcs := []struct {
		name      string
		code      string
		country   iso3166.CountryCode
		sname     string
		reference string
		errs      []error
	}{
		{name: "valid", code: "US-AL", country: us, sname: "Alabama", reference: "https://en.wikipedia.org/wiki/Alabama"},
		{name: "1 character suffix", code: "US-A", country: us, sname: "n"},
		{name: "digit suffix", code: "US-001", country: us, sname: "n"},
		{name: "name max", code: "US-AL", country: us, sname: strings.Repeat("a", 100)},
		{name: "name max counts runes", code: "US-AL", country: us, sname: strings.Repeat("é", 100)},
		{name: "code lowercase", code: "us-al", country: us, sname: "n", errs: []error{iso3166.ErrCode, iso3166.ErrCountryMismatch}},
		{name: "code without hyphen", code: "USAL", country: us, sname: "n", errs: []error{iso3166.ErrCode}},
		{name: "code suffix too long", code: "US-ALAB", country: us, sname: "n", errs: []error{iso3166.ErrCode}},
		{name: "code empty", code: "", country: us, sname: "n", errs: []error{iso3166.ErrCode, iso3166.ErrCountryMismatch}},
		{name: "code from another country", code: "CA-ON", country: us, sname: "n", errs: []error{iso3166.ErrCountryMismatch}},
		{name: "invalid country", code: "US-AL", country: iso3166.CountryCode{Alpha2: "US"}, sname: "n", errs: []error{iso3166.ErrCountry, iso3166.ErrInvalidCountryCode, iso3166.ErrAlpha3}},
		{name: "name empty", code: "US-AL", country: us, errs: []error{iso3166.ErrName}},
		{name: "name too long", code: "US-AL", country: us, sname: strings.Repeat("a", 101), errs: []error{iso3166.ErrName}},
		{name: "bad reference", code: "US-AL", country: us, sname: "n", reference: "nope", errs: []error{iso3166.ErrReference}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := iso3166.NewSubdivision(tc.code, tc.country, tc.sname, tc.reference)

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, iso3166.ErrInvalidSubdivision)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, iso3166.Subdivision{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.code, actual.Code)
			assert.Equal(tt, tc.country, actual.Country)
			assert.Equal(tt, tc.sname, actual.Name)
			assert.Equal(tt, optioner.OfZero(tc.reference), actual.Reference)
		})
	}
}

func Test_Subdivision_Validate(t *testing.T) {
	alabama := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "")
	county := iso3166.MustNewSubdivision("US-AL1", us, "County", "")
	ontario := iso3166.MustNewSubdivision("CA-ON", ca, "Ontario", "")

	withParent := func(child iso3166.Subdivision, parent iso3166.Subdivision) iso3166.Subdivision {
		child.Parent = optioner.Some(parent)
		return child
	}

	tcs := []struct {
		name  string
		input iso3166.Subdivision
		errs  []error
	}{
		{name: "parent", input: withParent(county, alabama)},
		{name: "grandparent", input: withParent(iso3166.MustNewSubdivision("US-AL2", us, "Town", ""), withParent(county, alabama))},
		{
			name: "all optional fields",
			input: iso3166.Subdivision{
				Code: "US-AL", Country: us, Name: "Alabama",
				Category:      optioner.Some("state"),
				AlternateCode: optioner.Some("AL"),
				PreviousCode:  optioner.Some("US-ALA"),
				Reference:     optioner.Some("https://example.com"),
			},
		},
		{name: "parent in another country", input: withParent(county, ontario), errs: []error{iso3166.ErrParentCountry}},
		{name: "parent is itself", input: withParent(alabama, alabama), errs: []error{iso3166.ErrParentSelf}},
		{name: "parent chain repeats a code", input: withParent(alabama, withParent(county, alabama)), errs: []error{iso3166.ErrParentCycle}},
		{name: "invalid parent", input: withParent(county, iso3166.Subdivision{Code: "US-AL", Country: us}), errs: []error{iso3166.ErrParent, iso3166.ErrName}},
		{name: "category empty", input: iso3166.Subdivision{Code: "US-AL", Country: us, Name: "n", Category: optioner.Some("")}, errs: []error{iso3166.ErrCategory}},
		{name: "category too long", input: iso3166.Subdivision{Code: "US-AL", Country: us, Name: "n", Category: optioner.Some(strings.Repeat("a", 51))}, errs: []error{iso3166.ErrCategory}},
		{name: "alternate code too long", input: iso3166.Subdivision{Code: "US-AL", Country: us, Name: "n", AlternateCode: optioner.Some(strings.Repeat("a", 11))}, errs: []error{iso3166.ErrAlternateCode}},
		{name: "previous code bad format", input: iso3166.Subdivision{Code: "US-AL", Country: us, Name: "n", PreviousCode: optioner.Some("ALA")}, errs: []error{iso3166.ErrPreviousCode}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, iso3166.ErrInvalidSubdivision)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				return
			}

			assert.NoError(tt, err)
		})
	}
}

func Test_Subdivision_JSON(t *testing.T) {
	alabama := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "https://en.wikipedia.org/wiki/Alabama")
	county := iso3166.MustNewSubdivision("US-AL1", us, "County", "")
	county.Category = optioner.Some("county")
	county.Parent = optioner.Some(alabama)

	countryJSON := `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"United States of America","independent":true}`

	tcs := []struct {
		name     string
		input    iso3166.Subdivision
		expected string
	}{
		{
			name:     "absent fields are omitted",
			input:    alabama,
			expected: `{"code":"US-AL","country":` + countryJSON + `,"name":"Alabama","reference":"https://en.wikipedia.org/wiki/Alabama"}`,
		},
		{
			name:  "nested parent",
			input: county,
			expected: `{"code":"US-AL1","country":` + countryJSON + `,"name":"County","category":"county",` +
				`"parent":{"code":"US-AL","country":` + countryJSON + `,"name":"Alabama","reference":"https://en.wikipedia.org/wiki/Alabama"}}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			data, err := json.Marshal(tc.input)
			assert.NoError(tt, err)
			assert.JSONEq(tt, tc.expected, string(data))

			var actual iso3166.Subdivision
			assert.NoError(tt, json.Unmarshal(data, &actual))
			assert.Equal(tt, tc.input, actual)
		})
	}
}

func Test_Subdivision_UnmarshalJSON_invalid(t *testing.T) {
	const country = `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"n","independent":true}`

	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing code", input: `{"country":` + country + `,"name":"n"}`, err: iso3166.ErrCode},
		{name: "missing country", input: `{"code":"US-AL","name":"n"}`, err: iso3166.ErrCountry},
		{name: "country as code", input: `{"code":"US-AL","country":"US","name":"n"}`},
		{name: "invalid country", input: `{"code":"US-AL","country":{"alpha-2":"US"},"name":"n"}`, err: iso3166.ErrInvalidCountryCode},
		{name: "mismatched country", input: `{"code":"CA-ON","country":` + country + `,"name":"n"}`, err: iso3166.ErrCountryMismatch},
		{name: "missing name", input: `{"code":"US-AL","country":` + country + `}`, err: iso3166.ErrName},
		{name: "invalid parent", input: `{"code":"US-AL1","country":` + country + `,"name":"n","parent":{"code":"US-AL","country":` + country + `}}`, err: iso3166.ErrName},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual iso3166.Subdivision
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.Error(tt, err)
			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
			}
			assert.Equal(tt, iso3166.Subdivision{}, actual)
		})
	}
}

func Test_Subdivision_UnmarshalJSON_null(t *testing.T) {
	expected := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "")
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
