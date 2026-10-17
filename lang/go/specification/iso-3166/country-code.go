// Package iso3166 implements the ISO 3166 Country Code and Subdivision data
// types, and lookups over the ISO 3166 data set.
//
// Source of truth: domain-model/specification/iso-3166/country-code.md,
// domain-model/specification/iso-3166/subdivision.md
package iso3166

import (
	"encoding/json"
	"regexp"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	countryAlpha2Validation    = validation.String[string]("alpha-2").Match(regexp.MustCompile(`^[A-Z]{2}$`))
	countryAlpha3Validation    = validation.String[string]("alpha-3").Match(regexp.MustCompile(`^[A-Z]{3}$`))
	countryNumericValidation   = validation.String[string]("numeric").Match(regexp.MustCompile(`^[0-9]{3}$`))
	countryNameValidation      = validation.String[string]("name").Min(1).Max(100)
	countryReferenceValidation = validation.String[string]("reference").URL("http", "https")

	// ErrInvalidCountryCode is wrapped by every CountryCode validation error.
	ErrInvalidCountryCode = errorer.New("invalid country code")
	ErrAlpha2             = errorer.New("alpha-2")
	ErrAlpha3             = errorer.New("alpha-3")
	ErrNumeric            = errorer.New("numeric")
	ErrName               = errorer.New("name")
	ErrIndependent        = errorer.New("independent")
	ErrReference          = errorer.New("reference")

	errCountryAlpha2Fn      = errorer.Func(ErrInvalidCountryCode, ErrAlpha2)
	errCountryAlpha3Fn      = errorer.Func(ErrInvalidCountryCode, ErrAlpha3)
	errCountryNumericFn     = errorer.Func(ErrInvalidCountryCode, ErrNumeric)
	errCountryNameFn        = errorer.Func(ErrInvalidCountryCode, ErrName)
	errCountryIndependentFn = errorer.Func(ErrInvalidCountryCode, ErrIndependent)
	errCountryReferenceFn   = errorer.Func(ErrInvalidCountryCode, ErrReference)
)

// CountryCode is an ISO 3166-1 country, dependent territory, or special area
// of geographical interest. Its natural key is Alpha2; it has no Id.
type CountryCode struct {
	Alpha2      string                  `json:"alpha-2"`
	Alpha3      string                  `json:"alpha-3"`
	Numeric     string                  `json:"numeric"`
	Name        string                  `json:"name"`
	Independent bool                    `json:"independent"`
	Reference   optioner.Option[string] `json:"reference,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidCountryCode and the field's sentinel
// (ErrAlpha2, etc.). Lengths are counted in characters (runes), not bytes.
func (this CountryCode) Validate() error {
	var errs []error

	if err := countryAlpha2Validation.Validate(this.Alpha2); err != nil {
		errs = append(errs, errCountryAlpha2Fn("%w", err))
	}

	if err := countryAlpha3Validation.Validate(this.Alpha3); err != nil {
		errs = append(errs, errCountryAlpha3Fn("%w", err))
	}

	if err := countryNumericValidation.Validate(this.Numeric); err != nil {
		errs = append(errs, errCountryNumericFn("%w", err))
	}

	if err := countryNameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errCountryNameFn("%w", err))
	}

	if this.Reference.Defined() {
		if err := countryReferenceValidation.Validate(this.Reference.Get()); err != nil {
			errs = append(errs, errCountryReferenceFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this CountryCode) MarshalJSON() ([]byte, error) {
	type alias CountryCode
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// CountryCode. A JSON null leaves the CountryCode unchanged.
//
// Independent is a required boolean, and a missing boolean would otherwise
// decode as false, so a missing "independent" key fails with ErrIndependent.
func (this *CountryCode) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias CountryCode
	var decoded struct {
		alias
		Independent *bool `json:"independent"`
	}

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	country := CountryCode(decoded.alias)
	var errs []error

	if decoded.Independent == nil {
		errs = append(errs, errCountryIndependentFn("required"))
	} else {
		country.Independent = *decoded.Independent
	}

	errs = append(errs, country.Validate())

	if err := errorer.Join(errs...); err != nil {
		return err
	}

	*this = country
	return nil
}

// NewCountryCode returns a validated CountryCode.
// An empty reference is treated as not present.
func NewCountryCode(alpha2, alpha3, numeric, name string, independent bool, reference string) (CountryCode, error) {
	country := CountryCode{
		Alpha2:      alpha2,
		Alpha3:      alpha3,
		Numeric:     numeric,
		Name:        name,
		Independent: independent,
		Reference:   optioner.OfZero(reference),
	}

	if err := country.Validate(); err != nil {
		return CountryCode{}, err
	}

	return country, nil
}

// MustNewCountryCode is like NewCountryCode but panics if the CountryCode is
// invalid.
func MustNewCountryCode(alpha2, alpha3, numeric, name string, independent bool, reference string) CountryCode {
	country, err := NewCountryCode(alpha2, alpha3, numeric, name, independent, reference)

	if err != nil {
		panic(err)
	}

	return country
}
