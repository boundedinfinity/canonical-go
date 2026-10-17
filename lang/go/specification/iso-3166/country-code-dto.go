package iso3166

// Source of truth: domain-model/specification/iso-3166/country-code.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"

	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

// CountryCodeDto is the data transfer object for CountryCode. CountryCode
// holds no other Top Level Models, so CountryCodeDto has the same fields.
type CountryCodeDto struct {
	Alpha2      string                  `json:"alpha-2"`
	Alpha3      string                  `json:"alpha-3"`
	Numeric     string                  `json:"numeric"`
	Name        string                  `json:"name"`
	Independent bool                    `json:"independent"`
	Reference   optioner.Option[string] `json:"reference,omitzero"`
}

// Validate applies the CountryCode rules; see CountryCode.Validate.
func (this CountryCodeDto) Validate() error {
	return CountryCode(this).Validate()
}

// MarshalJSON implements json.Marshaler.
func (this CountryCodeDto) MarshalJSON() ([]byte, error) {
	return json.Marshal(CountryCode(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// CountryCodeDto, with the same rules as CountryCode.UnmarshalJSON.
// A JSON null leaves the CountryCodeDto unchanged.
func (this *CountryCodeDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var country CountryCode

	if err := json.Unmarshal(data, &country); err != nil {
		return err
	}

	*this = CountryCodeDto(country)
	return nil
}

// ToModel returns the validated CountryCode for this CountryCodeDto.
func (this CountryCodeDto) ToModel() (CountryCode, error) {
	country := CountryCode(this)

	if err := country.Validate(); err != nil {
		return CountryCode{}, err
	}

	return country, nil
}

// ToDto returns the CountryCodeDto for this CountryCode.
func (this CountryCode) ToDto() CountryCodeDto {
	return CountryCodeDto(this)
}
