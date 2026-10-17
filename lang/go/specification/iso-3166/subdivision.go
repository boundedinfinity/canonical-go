package iso3166

// Source of truth: domain-model/specification/iso-3166/subdivision.md

import (
	"encoding/json"
	"regexp"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	subdivisionCodePattern = regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)

	subdivisionCodeValidation          = validation.String[string]("code").Match(subdivisionCodePattern)
	subdivisionNameValidation          = validation.String[string]("name").Min(1).Max(100)
	subdivisionCategoryValidation      = validation.String[string]("category").Min(1).Max(50)
	subdivisionParentValidation        = validation.String[string]("parent").Match(subdivisionCodePattern)
	subdivisionAlternateCodeValidation = validation.String[string]("alternate-code").Min(1).Max(10)
	subdivisionPreviousCodeValidation  = validation.String[string]("previous-code").Match(subdivisionCodePattern)
	subdivisionReferenceValidation     = validation.String[string]("reference").URL("http", "https")

	// ErrInvalidSubdivision is wrapped by every Subdivision validation error.
	// The ErrName and ErrReference field sentinels are shared with CountryCode.
	ErrInvalidSubdivision = errorer.New("invalid subdivision")
	ErrCode               = errorer.New("code")
	ErrCountry            = errorer.New("country")
	ErrCategory           = errorer.New("category")
	ErrParent             = errorer.New("parent")
	ErrAlternateCode      = errorer.New("alternate-code")
	ErrPreviousCode       = errorer.New("previous-code")
	ErrCountryMismatch    = errorer.New("code does not start with the country's alpha-2")
	ErrParentCountry      = errorer.New("parent belongs to a different country")
	ErrParentSelf         = errorer.New("parent is the subdivision itself")
	ErrParentCycle        = errorer.New("parent chain repeats a code")

	errSubdivisionCodeFn            = errorer.Func(ErrInvalidSubdivision, ErrCode)
	errSubdivisionCountryFn         = errorer.Func(ErrInvalidSubdivision, ErrCountry)
	errSubdivisionNameFn            = errorer.Func(ErrInvalidSubdivision, ErrName)
	errSubdivisionCategoryFn        = errorer.Func(ErrInvalidSubdivision, ErrCategory)
	errSubdivisionParentFn          = errorer.Func(ErrInvalidSubdivision, ErrParent)
	errSubdivisionAlternateCodeFn   = errorer.Func(ErrInvalidSubdivision, ErrAlternateCode)
	errSubdivisionPreviousCodeFn    = errorer.Func(ErrInvalidSubdivision, ErrPreviousCode)
	errSubdivisionReferenceFn       = errorer.Func(ErrInvalidSubdivision, ErrReference)
	errSubdivisionCountryMismatchFn = errorer.Func(ErrInvalidSubdivision, ErrCountryMismatch)
	errSubdivisionParentCountryFn   = errorer.Func(ErrInvalidSubdivision, ErrParentCountry)
	errSubdivisionParentSelfFn      = errorer.Func(ErrInvalidSubdivision, ErrParentSelf)
	errSubdivisionParentCycleFn     = errorer.Func(ErrInvalidSubdivision, ErrParentCycle)
)

// Subdivision is an ISO 3166-2 principal subdivision of a country (a state,
// province, region, etc.). Its natural key is Code; it has no Id.
type Subdivision struct {
	Code          string                       `json:"code"`
	Country       CountryCode                  `json:"country"`
	Name          string                       `json:"name"`
	Category      optioner.Option[string]      `json:"category,omitzero"`
	Parent        optioner.Option[Subdivision] `json:"parent,omitzero"`
	AlternateCode optioner.Option[string]      `json:"alternate-code,omitzero"`
	PreviousCode  optioner.Option[string]      `json:"previous-code,omitzero"`
	Reference     optioner.Option[string]      `json:"reference,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidSubdivision and the field's sentinel (ErrCode,
// etc.). Country and Parent failures also wrap their own errors.
// Lengths are counted in characters (runes), not bytes.
func (this Subdivision) Validate() error {
	var errs []error

	if err := subdivisionCodeValidation.Validate(this.Code); err != nil {
		errs = append(errs, errSubdivisionCodeFn("%w", err))
	}

	if err := this.Country.Validate(); err != nil {
		errs = append(errs, errSubdivisionCountryFn("%w", err))
	}

	if len(this.Code) < 2 || this.Code[:2] != this.Country.Alpha2 {
		errs = append(errs, errSubdivisionCountryMismatchFn("%q, %q", this.Code, this.Country.Alpha2))
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

		if err := parent.Validate(); err != nil {
			errs = append(errs, errSubdivisionParentFn("%w", err))
		}

		if parent.Country.Alpha2 != this.Country.Alpha2 {
			errs = append(errs, errSubdivisionParentCountryFn("%q, %q", parent.Country.Alpha2, this.Country.Alpha2))
		}

		if parent.Code == this.Code {
			errs = append(errs, errSubdivisionParentSelfFn("%q", this.Code))
		} else {
			seen := map[string]bool{this.Code: true}

			for p := this.Parent; p.Defined(); p = p.Get().Parent {
				code := p.Get().Code

				if seen[code] {
					errs = append(errs, errSubdivisionParentCycleFn("%q", code))
					break
				}

				seen[code] = true
			}
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
func (this Subdivision) MarshalJSON() ([]byte, error) {
	type alias Subdivision
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// Subdivision. A JSON null leaves the Subdivision unchanged.
func (this *Subdivision) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Subdivision
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	subdivision := Subdivision(decoded)

	if err := subdivision.Validate(); err != nil {
		return err
	}

	*this = subdivision
	return nil
}

// NewSubdivision returns a validated Subdivision with no Category, Parent,
// AlternateCode, or PreviousCode; set those on the result and call Validate
// if they are needed. An empty reference is treated as not present.
func NewSubdivision(code string, country CountryCode, name, reference string) (Subdivision, error) {
	subdivision := Subdivision{
		Code:      code,
		Country:   country,
		Name:      name,
		Reference: optioner.OfZero(reference),
	}

	if err := subdivision.Validate(); err != nil {
		return Subdivision{}, err
	}

	return subdivision, nil
}

// MustNewSubdivision is like NewSubdivision but panics if the Subdivision is
// invalid.
func MustNewSubdivision(code string, country CountryCode, name, reference string) Subdivision {
	subdivision, err := NewSubdivision(code, country, name, reference)

	if err != nil {
		panic(err)
	}

	return subdivision
}
