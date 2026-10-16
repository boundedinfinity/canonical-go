// Package affix implements the Affix, Prefix, Suffix, and affix Category data
// types.
//
// Source of truth: domain-model/person/affix/affix.md,
// domain-model/person/affix/prefix.md, domain-model/person/affix/suffix.md,
// domain-model/person/affix/category.md,
// domain-model/person/affix/prefix-category.md,
// domain-model/person/affix/suffix-category.md
package affix

import (
	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	nameValidation        = validation.String[string]("name").Min(1).Max(30)
	descriptionValidation = validation.String[string]("description").Min(2).Max(500)

	// ErrInvalid is wrapped by every Affix, Prefix, and Suffix validation error.
	ErrInvalid     = errorer.New("invalid affix")
	ErrId          = id.ErrInvalid
	ErrName        = errorer.New("name")
	ErrDescription = errorer.New("description")
	ErrCategory    = errorer.New("category")

	errIdFn          = errorer.Func(ErrInvalid, id.ErrInvalid)
	errNameFn        = errorer.Func(ErrInvalid, ErrName)
	errDescriptionFn = errorer.Func(ErrInvalid, ErrDescription)
	errCategoryFn    = errorer.Func(ErrInvalid, ErrCategory)
)

// Affix holds the fields shared by Prefix and Suffix.
//
// The spec's Category field is declared on Prefix (as a PrefixCategory) and
// Suffix (as a SuffixCategory) so the compiler enforces which category each
// one uses. Affix has no JSON methods of its own: they would be promoted to
// Prefix and Suffix and drop their category.
type Affix struct {
	Id          id.ID                   `json:"id"`
	Name        string                  `json:"name"`
	Description optioner.Option[string] `json:"description,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
// Lengths are counted in characters (runes), not bytes.
func (this Affix) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errIdFn("%w", err))
	}

	if err := nameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errNameFn("%w", err))
	}

	if this.Description.Defined() {
		if err := descriptionValidation.Validate(this.Description.Get()); err != nil {
			errs = append(errs, errDescriptionFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}
