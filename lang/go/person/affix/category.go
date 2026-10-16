package affix

// Source of truth: domain-model/person/affix/category.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	categoryNameValidation        = validation.String[string]("name").Min(1).Max(30)
	categoryDescriptionValidation = validation.String[string]("description").Min(2).Max(500)

	// ErrInvalidCategory is wrapped by every Category validation error. The
	// field sentinels (ErrId, ErrName, ErrDescription) are shared with Affix.
	ErrInvalidCategory = errorer.New("invalid affix category")

	errCategoryIdFn          = errorer.Func(ErrInvalidCategory, id.ErrInvalid)
	errCategoryNameFn        = errorer.Func(ErrInvalidCategory, ErrName)
	errCategoryDescriptionFn = errorer.Func(ErrInvalidCategory, ErrDescription)
)

// Category holds the fields shared by PrefixCategory and SuffixCategory.
type Category struct {
	Id          id.ID                   `json:"id"`
	Name        string                  `json:"name"`
	Description optioner.Option[string] `json:"description,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidCategory and the field's sentinel (ErrName, etc.).
// Lengths are counted in characters (runes), not bytes.
func (this Category) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errCategoryIdFn("%w", err))
	}

	if err := categoryNameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errCategoryNameFn("%w", err))
	}

	if this.Description.Defined() {
		if err := categoryDescriptionValidation.Validate(this.Description.Get()); err != nil {
			errs = append(errs, errCategoryDescriptionFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Category) MarshalJSON() ([]byte, error) {
	type alias Category
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Category.
// A JSON null leaves the Category unchanged.
func (this *Category) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Category
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	category := Category(decoded)

	if err := category.Validate(); err != nil {
		return err
	}

	*this = category
	return nil
}

// NewCategory returns a validated Category with a new ID.
// An empty description is treated as not present.
func NewCategory(name, description string) (Category, error) {
	category := Category{
		Id:          id.New(),
		Name:        name,
		Description: optioner.OfZero(description),
	}

	if err := category.Validate(); err != nil {
		return Category{}, err
	}

	return category, nil
}

// MustNewCategory is like NewCategory but panics if the Category is invalid.
func MustNewCategory(name, description string) Category {
	category, err := NewCategory(name, description)

	if err != nil {
		panic(err)
	}

	return category
}
