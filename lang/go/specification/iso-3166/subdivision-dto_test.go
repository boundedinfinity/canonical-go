package iso3166_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	iso3166 "github.com/boundedinfinity/canonical-go/lang/go/specification/iso-3166"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_SubdivisionDto_ToDto(t *testing.T) {
	alabama := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "")
	county := iso3166.MustNewSubdivision("US-AL1", us, "County", "https://example.com")
	county.Category = optioner.Some("county")
	county.Parent = optioner.Some(alabama)

	dto := county.ToDto()

	assert.Equal(t, iso3166.SubdivisionDto{
		Code:      "US-AL1",
		Country:   "US",
		Name:      "County",
		Category:  optioner.Some("county"),
		Parent:    optioner.Some("US-AL"),
		Reference: optioner.Some("https://example.com"),
	}, dto)
	assert.True(t, alabama.ToDto().Parent.Empty())
}

func Test_SubdivisionDto_ToModel(t *testing.T) {
	errNotFound := errors.New("not found")
	alabama := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "")
	ontario := iso3166.MustNewSubdivision("CA-ON", ca, "Ontario", "")
	county := iso3166.MustNewSubdivision("US-AL1", us, "County", "")
	county.Parent = optioner.Some(alabama)

	countries := resolver.Func[string, iso3166.CountryCode](func(alpha2 string) (iso3166.CountryCode, error) {
		switch alpha2 {
		case "US":
			return us, nil
		case "CA":
			return ca, nil
		}

		return iso3166.CountryCode{}, errNotFound
	})

	// mislabeled returns Ontario for US-AL, so the resolved parent is in the
	// wrong country even though its DTO code looked right.
	mislabeled := resolver.Func[string, iso3166.Subdivision](func(string) (iso3166.Subdivision, error) {
		return ontario, nil
	})

	subdivisions := resolver.Func[string, iso3166.Subdivision](func(code string) (iso3166.Subdivision, error) {
		if code == alabama.Code {
			return alabama, nil
		}

		return iso3166.Subdivision{}, errNotFound
	})

	tcs := []struct {
		name         string
		input        iso3166.SubdivisionDto
		countries    resolver.Resolver[string, iso3166.CountryCode]
		subdivisions resolver.Resolver[string, iso3166.Subdivision]
		expected     iso3166.Subdivision
		errs         []error
	}{
		{name: "round trip", input: county.ToDto(), countries: countries, subdivisions: subdivisions, expected: county},
		{name: "no parent needs no subdivision resolver", input: alabama.ToDto(), countries: countries, expected: alabama},
		{
			name:  "invalid dto",
			input: iso3166.SubdivisionDto{Code: "US-AL", Country: "us"},
			errs:  []error{iso3166.ErrInvalidSubdivision, iso3166.ErrCountry, iso3166.ErrName},
		},
		{
			name:      "unknown country",
			input:     iso3166.SubdivisionDto{Code: "ZZ-AA", Country: "ZZ", Name: "n"},
			countries: countries,
			errs:      []error{resolver.ErrResolve, errNotFound},
		},
		{
			name:         "unknown parent",
			input:        iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "n", Parent: optioner.Some("US-ZZ")},
			countries:    countries,
			subdivisions: subdivisions,
			errs:         []error{resolver.ErrResolve, errNotFound},
		},
		{
			name:      "missing subdivision resolver",
			input:     county.ToDto(),
			countries: countries,
			errs:      []error{resolver.ErrResolve, resolver.ErrNoResolver},
		},
		{
			name:         "resolved parent in another country",
			input:        county.ToDto(),
			countries:    countries,
			subdivisions: mislabeled,
			errs:         []error{iso3166.ErrInvalidSubdivision, iso3166.ErrParentCountry},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			actual, err := tc.input.ToModel(tc.countries, tc.subdivisions)

			if len(tc.errs) > 0 {
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				assert.Equal(tt, iso3166.Subdivision{}, actual)
				return
			}

			assert.NoError(tt, err)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func Test_SubdivisionDto_Validate(t *testing.T) {
	tcs := []struct {
		name  string
		input iso3166.SubdivisionDto
		errs  []error
	}{
		{name: "valid", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US", Name: "n"}},
		{name: "with parent", input: iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "n", Parent: optioner.Some("US-AL")}},
		{name: "bad code", input: iso3166.SubdivisionDto{Code: "USAL", Country: "US", Name: "n"}, errs: []error{iso3166.ErrCode}},
		{name: "bad country", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "USA", Name: "n"}, errs: []error{iso3166.ErrCountry, iso3166.ErrCountryMismatch}},
		{name: "mismatched country", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "CA", Name: "n"}, errs: []error{iso3166.ErrCountryMismatch}},
		{name: "name empty", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US"}, errs: []error{iso3166.ErrName}},
		{name: "bad parent code", input: iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "n", Parent: optioner.Some("AL")}, errs: []error{iso3166.ErrParent, iso3166.ErrParentCountry}},
		{name: "parent in another country", input: iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "n", Parent: optioner.Some("CA-ON")}, errs: []error{iso3166.ErrParentCountry}},
		{name: "parent is itself", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US", Name: "n", Parent: optioner.Some("US-AL")}, errs: []error{iso3166.ErrParentSelf}},
		{name: "category empty", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US", Name: "n", Category: optioner.Some("")}, errs: []error{iso3166.ErrCategory}},
		{name: "previous code bad format", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US", Name: "n", PreviousCode: optioner.Some("x")}, errs: []error{iso3166.ErrPreviousCode}},
		{name: "bad reference", input: iso3166.SubdivisionDto{Code: "US-AL", Country: "US", Name: "n", Reference: optioner.Some("x")}, errs: []error{iso3166.ErrReference}},
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

func Test_SubdivisionDto_JSON(t *testing.T) {
	input := iso3166.SubdivisionDto{Code: "US-AL1", Country: "US", Name: "County", Parent: optioner.Some("US-AL"), Reference: optioner.Some("https://example.com")}
	expected := `{"code":"US-AL1","country":"US","name":"County","parent":"US-AL","reference":"https://example.com"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual iso3166.SubdivisionDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_SubdivisionDto_UnmarshalJSON_invalid(t *testing.T) {
	tcs := []struct {
		name  string
		input string
		err   error
	}{
		{name: "missing country", input: `{"code":"US-AL","name":"n"}`, err: iso3166.ErrCountry},
		{name: "country as object", input: `{"code":"US-AL","country":{"alpha-2":"US"},"name":"n"}`},
		{name: "parent as object", input: `{"code":"US-AL1","country":"US","name":"n","parent":{"code":"US-AL"}}`},
		{name: "parent in another country", input: `{"code":"US-AL1","country":"US","name":"n","parent":"CA-ON"}`, err: iso3166.ErrParentCountry},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			var actual iso3166.SubdivisionDto
			err := json.Unmarshal([]byte(tc.input), &actual)

			assert.Error(tt, err)
			if tc.err != nil {
				assert.ErrorIs(tt, err, tc.err)
			}
			assert.Equal(tt, iso3166.SubdivisionDto{}, actual)
		})
	}
}

func Test_SubdivisionDto_UnmarshalJSON_null(t *testing.T) {
	expected := iso3166.MustNewSubdivision("US-AL", us, "Alabama", "").ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
