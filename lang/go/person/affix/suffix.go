package affix

// Source of truth: domain-model/person/affix/suffix.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

// Suffix is a name suffix, such as "Jr.". It embeds the Affix fields and adds
// a required SuffixCategory. JSON is flat: {"id", "name", "description",
// "category"}.
type Suffix struct {
	Affix
	Category SuffixCategory `json:"category"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
// A category failure wraps ErrCategory and the Category's own errors.
func (this Suffix) Validate() error {
	errs := []error{this.Affix.Validate()}

	if err := this.Category.Validate(); err != nil {
		errs = append(errs, errCategoryFn("%w", err))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Suffix) MarshalJSON() ([]byte, error) {
	type alias Suffix
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Suffix.
// A JSON null leaves the Suffix unchanged.
func (this *Suffix) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Suffix
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	suffix := Suffix(decoded)

	if err := suffix.Validate(); err != nil {
		return err
	}

	*this = suffix
	return nil
}

// NewSuffix returns a validated Suffix with a new ID.
// An empty description is treated as not present.
func NewSuffix(name, description string, category SuffixCategory) (Suffix, error) {
	suffix := Suffix{
		Id:          id.New(),
		Name:        name,
		Description: optioner.OfZero(description),
		Category:    category,
	}

	if err := suffix.Validate(); err != nil {
		return Suffix{}, err
	}

	return suffix, nil
}

// MustNewSuffix is like NewSuffix but panics if the Suffix is invalid.
func MustNewSuffix(name, description string, category SuffixCategory) Suffix {
	suffix, err := NewSuffix(name, description, category)

	if err != nil {
		panic(err)
	}

	return suffix
}
