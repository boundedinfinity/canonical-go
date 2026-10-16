package affix

// Source of truth: domain-model/person/affix/suffix.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/errorer"
)

// SuffixDto is the data transfer object for Suffix. Category holds the ID of
// the SuffixCategory instead of the category itself. JSON is flat:
// {"id", "name", "description", "category"}.
type SuffixDto struct {
	Affix
	Category id.ID `json:"category"`
}

// Validate applies the Suffix rules, with Category checked as an ID.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
func (this SuffixDto) Validate() error {
	errs := []error{this.Affix.Validate()}

	if err := this.Category.Validate(); err != nil {
		errs = append(errs, errCategoryFn("%w", err))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this SuffixDto) MarshalJSON() ([]byte, error) {
	type alias SuffixDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// SuffixDto. A JSON null leaves the SuffixDto unchanged.
func (this *SuffixDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias SuffixDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := SuffixDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel resolves Category and returns the validated Suffix.
// A resolution failure wraps resolver.ErrResolve.
func (this SuffixDto) ToModel(categories resolver.Resolver[SuffixCategory]) (Suffix, error) {
	if err := this.Validate(); err != nil {
		return Suffix{}, err
	}

	category, err := resolver.One(categories, this.Category)

	if err != nil {
		return Suffix{}, fmt.Errorf("category %w", err)
	}

	suffix := Suffix{Affix: this.Affix, Category: category}

	if err := suffix.Validate(); err != nil {
		return Suffix{}, err
	}

	return suffix, nil
}

// ToDto returns the SuffixDto for this Suffix.
func (this Suffix) ToDto() SuffixDto {
	return SuffixDto{Affix: this.Affix, Category: this.Category.Id}
}
