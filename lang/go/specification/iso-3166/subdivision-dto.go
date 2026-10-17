package iso3166

// Source of truth: domain-model/specification/iso-3166/subdivision.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

// SubdivisionDto is the data transfer object for Subdivision. Country holds
// the country's Alpha2, and Parent holds the parent's Code, instead of the
// models themselves.
type SubdivisionDto struct {
	Code          string                  `json:"code"`
	Country       string                  `json:"country"`
	Name          string                  `json:"name"`
	Category      optioner.Option[string] `json:"category,omitzero"`
	Parent        optioner.Option[string] `json:"parent,omitzero"`
	AlternateCode optioner.Option[string] `json:"alternate-code,omitzero"`
	PreviousCode  optioner.Option[string] `json:"previous-code,omitzero"`
	Reference     optioner.Option[string] `json:"reference,omitzero"`
}

// Validate applies the Subdivision rules, with Country checked as an Alpha2
// and Parent checked as a Code. A parent belongs to the same country when its
// Code starts with Country. The no-cycles rule needs the whole parent chain,
// so it is checked by ToModel.
// Each failure wraps ErrInvalidSubdivision and the field's sentinel.
func (this SubdivisionDto) Validate() error {
	var errs []error

	if err := subdivisionCodeValidation.Validate(this.Code); err != nil {
		errs = append(errs, errSubdivisionCodeFn("%w", err))
	}

	if err := countryAlpha2Validation.Validate(this.Country); err != nil {
		errs = append(errs, errSubdivisionCountryFn("%w", err))
	}

	if len(this.Code) < 2 || this.Code[:2] != this.Country {
		errs = append(errs, errSubdivisionCountryMismatchFn("%q, %q", this.Code, this.Country))
	}

	if err := subdivisionNameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errSubdivisionNameFn("%w", err))
	}

	if this.Category.Defined() {
		if err := subdivisionCategoryValidation.Validate(this.Category.Get()); err != nil {
			errs = append(errs, errSubdivisionCategoryFn("%w", err))
		}
	}

	if this.Parent.Defined() {
		parent := this.Parent.Get()

		if err := subdivisionParentValidation.Validate(parent); err != nil {
			errs = append(errs, errSubdivisionParentFn("%w", err))
		}

		if len(parent) < 2 || parent[:2] != this.Country {
			errs = append(errs, errSubdivisionParentCountryFn("%q, %q", parent, this.Country))
		}

		if parent == this.Code {
			errs = append(errs, errSubdivisionParentSelfFn("%q", this.Code))
		}
	}

	if this.AlternateCode.Defined() {
		if err := subdivisionAlternateCodeValidation.Validate(this.AlternateCode.Get()); err != nil {
			errs = append(errs, errSubdivisionAlternateCodeFn("%w", err))
		}
	}

	if this.PreviousCode.Defined() {
		if err := subdivisionPreviousCodeValidation.Validate(this.PreviousCode.Get()); err != nil {
			errs = append(errs, errSubdivisionPreviousCodeFn("%w", err))
		}
	}

	if this.Reference.Defined() {
		if err := subdivisionReferenceValidation.Validate(this.Reference.Get()); err != nil {
			errs = append(errs, errSubdivisionReferenceFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this SubdivisionDto) MarshalJSON() ([]byte, error) {
	type alias SubdivisionDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// SubdivisionDto. A JSON null leaves the SubdivisionDto unchanged.
func (this *SubdivisionDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias SubdivisionDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := SubdivisionDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel resolves Country (by Alpha2) and Parent (by Code) and returns the
// validated Subdivision. Every resolution failure is reported and wraps
// resolver.ErrResolve. The subdivisions resolver is only called when Parent is
// present, so it may be nil otherwise.
func (this SubdivisionDto) ToModel(
	countries resolver.Resolver[string, CountryCode],
	subdivisions resolver.Resolver[string, Subdivision],
) (Subdivision, error) {
	if err := this.Validate(); err != nil {
		return Subdivision{}, err
	}

	var errs []error

	country, err := resolver.One(countries, this.Country)

	if err != nil {
		errs = append(errs, fmt.Errorf("country %w", err))
	}

	parent, err := resolver.Option(subdivisions, this.Parent)

	if err != nil {
		errs = append(errs, fmt.Errorf("parent %w", err))
	}

	if err := errorer.Join(errs...); err != nil {
		return Subdivision{}, err
	}

	subdivision := Subdivision{
		Code:          this.Code,
		Country:       country,
		Name:          this.Name,
		Category:      this.Category,
		Parent:        parent,
		AlternateCode: this.AlternateCode,
		PreviousCode:  this.PreviousCode,
		Reference:     this.Reference,
	}

	if err := subdivision.Validate(); err != nil {
		return Subdivision{}, err
	}

	return subdivision, nil
}

// ToDto returns the SubdivisionDto for this Subdivision.
func (this Subdivision) ToDto() SubdivisionDto {
	return SubdivisionDto{
		Code:          this.Code,
		Country:       this.Country.Alpha2,
		Name:          this.Name,
		Category:      this.Category,
		Parent:        optioner.Map(this.Parent, subdivisionCode),
		AlternateCode: this.AlternateCode,
		PreviousCode:  this.PreviousCode,
		Reference:     this.Reference,
	}
}

func subdivisionCode(subdivision Subdivision) string {
	return subdivision.Code
}
