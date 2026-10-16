package affix

// Source of truth: domain-model/person/affix/prefix.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

// Prefix is a name prefix, such as "Dr.". It embeds the Affix fields and adds
// a required PrefixCategory. JSON is flat: {"id", "name", "description",
// "category"}.
type Prefix struct {
	Affix
	Category PrefixCategory `json:"category"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
// A category failure wraps ErrCategory and the Category's own errors.
func (this Prefix) Validate() error {
	errs := []error{this.Affix.Validate()}

	if err := this.Category.Validate(); err != nil {
		errs = append(errs, errCategoryFn("%w", err))
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Prefix) MarshalJSON() ([]byte, error) {
	type alias Prefix
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Prefix.
// A JSON null leaves the Prefix unchanged.
func (this *Prefix) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Prefix
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	prefix := Prefix(decoded)

	if err := prefix.Validate(); err != nil {
		return err
	}

	*this = prefix
	return nil
}

// NewPrefix returns a validated Prefix with a new ID.
// An empty description is treated as not present.
func NewPrefix(name, description string, category PrefixCategory) (Prefix, error) {
	prefix := Prefix{
		Id:          id.New(),
		Name:        name,
		Description: optioner.OfZero(description),
		Category:    category,
	}

	if err := prefix.Validate(); err != nil {
		return Prefix{}, err
	}

	return prefix, nil
}

// MustNewPrefix is like NewPrefix but panics if the Prefix is invalid.
func MustNewPrefix(name, description string, category PrefixCategory) Prefix {
	prefix, err := NewPrefix(name, description, category)

	if err != nil {
		panic(err)
	}

	return prefix
}
