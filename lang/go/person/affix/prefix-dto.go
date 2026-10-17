package affix

// Source of truth: domain-model/person/affix/prefix.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/errorer"
)

// PrefixDto is the data transfer object for Prefix. Category holds the ID of
// the PrefixCategory instead of the category itself. JSON is flat:
// {"id", "name", "description", "category"}.
type PrefixDto struct {
	Affix
	Category id.ID `json:"category"`
}

// Validate applies the Prefix rules, with Category checked as an ID.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
func (this PrefixDto) Validate() error {
	errs := []error{this.Affix.Validate()}

	if err := this.Category.Validate(); err != nil {
		errs = append(errs, errCategoryFn("%w", err))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this PrefixDto) MarshalJSON() ([]byte, error) {
	type alias PrefixDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// PrefixDto. A JSON null leaves the PrefixDto unchanged.
func (this *PrefixDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias PrefixDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := PrefixDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel resolves Category and returns the validated Prefix.
// A resolution failure wraps resolver.ErrResolve.
func (this PrefixDto) ToModel(categories resolver.Resolver[id.ID, PrefixCategory]) (Prefix, error) {
	if err := this.Validate(); err != nil {
		return Prefix{}, err
	}

	category, err := resolver.One(categories, this.Category)

	if err != nil {
		return Prefix{}, fmt.Errorf("category %w", err)
	}

	prefix := Prefix{Affix: this.Affix, Category: category}

	if err := prefix.Validate(); err != nil {
		return Prefix{}, err
	}

	return prefix, nil
}

// ToDto returns the PrefixDto for this Prefix.
func (this Prefix) ToDto() PrefixDto {
	return PrefixDto{Affix: this.Affix, Category: this.Category.Id}
}
