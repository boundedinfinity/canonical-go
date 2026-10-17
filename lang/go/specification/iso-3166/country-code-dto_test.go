package iso3166_test

import (
	"encoding/json"
	"testing"

	iso3166 "github.com/boundedinfinity/canonical-go/lang/go/specification/iso-3166"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_CountryCodeDto_round_trip(t *testing.T) {
	expected := iso3166.MustNewCountryCode("US", "USA", "840", "United States of America", true, "https://en.wikipedia.org/wiki/United_States")

	dto := expected.ToDto()
	assert.Equal(t, expected.Alpha2, dto.Alpha2)
	assert.Equal(t, expected.Alpha3, dto.Alpha3)
	assert.Equal(t, expected.Numeric, dto.Numeric)
	assert.Equal(t, expected.Name, dto.Name)
	assert.Equal(t, expected.Independent, dto.Independent)
	assert.Equal(t, expected.Reference, dto.Reference)

	actual, err := dto.ToModel()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_CountryCodeDto_ToModel_invalid(t *testing.T) {
	actual, err := iso3166.CountryCodeDto{Alpha2: "us", Alpha3: "USA", Numeric: "840", Name: "n"}.ToModel()

	assert.ErrorIs(t, err, iso3166.ErrInvalidCountryCode)
	assert.ErrorIs(t, err, iso3166.ErrAlpha2)
	assert.Equal(t, iso3166.CountryCode{}, actual)
}

func Test_CountryCodeDto_JSON(t *testing.T) {
	input := iso3166.CountryCodeDto{Alpha2: "US", Alpha3: "USA", Numeric: "840", Name: "United States of America", Independent: true, Reference: optioner.Some("https://example.com")}
	expected := `{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"United States of America","independent":true,"reference":"https://example.com"}`

	data, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(data))

	var actual iso3166.CountryCodeDto
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, input, actual)
}

func Test_CountryCodeDto_UnmarshalJSON_invalid(t *testing.T) {
	var actual iso3166.CountryCodeDto
	err := json.Unmarshal([]byte(`{"alpha-2":"US","alpha-3":"USA","numeric":"840","name":"n"}`), &actual)

	assert.ErrorIs(t, err, iso3166.ErrIndependent)
	assert.Equal(t, iso3166.CountryCodeDto{}, actual)
}

func Test_CountryCodeDto_UnmarshalJSON_null(t *testing.T) {
	expected := iso3166.MustNewCountryCode("US", "USA", "840", "n", true, "").ToDto()
	actual := expected

	assert.NoError(t, json.Unmarshal([]byte(`null`), &actual))
	assert.Equal(t, expected, actual)
}
