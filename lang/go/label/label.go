// Package label implements the Label data type.
//
// Source of truth: domain-model/label/label.md
package label

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	nameValidation        = validation.String[string]("name").Min(1).Max(20)
	acronymValidation     = validation.String[string]("acronym").Min(2).Max(10)
	descriptionValidation = validation.String[string]("description").Min(1).Max(250)

	// ErrInvalid is wrapped by every Label validation error.
	ErrInvalid     = errorer.New("invalid label")
	ErrId          = id.ErrInvalid
	ErrName        = errorer.New("name")
	ErrAcronym     = errorer.New("acronym")
	ErrDescription = errorer.New("description")

	errIdFn          = errorer.Func(ErrInvalid, id.ErrInvalid)
	errNameFn        = errorer.Func(ErrInvalid, ErrName)
	errAcronymFn     = errorer.Func(ErrInvalid, ErrAcronym)
	errDescriptionFn = errorer.Func(ErrInvalid, ErrDescription)
)

// Label is a user defined label for a document.
type Label struct {
	Id          id.ID                   `json:"id"`
	Name        string                  `json:"name"`
	Acronym     optioner.Option[string] `json:"acronym,omitzero"`
	Description optioner.Option[string] `json:"description,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalid and the field's sentinel (ErrName, etc.).
// Lengths are counted in characters (runes), not bytes.
func (this Label) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errIdFn("%w", err))
	}

	if err := nameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errNameFn("%w", err))
	}

	if this.Acronym.Defined() {
		if err := acronymValidation.Validate(this.Acronym.Get()); err != nil {
			errs = append(errs, errAcronymFn("%w", err))
		}
	}

	if this.Description.Defined() {
		if err := descriptionValidation.Validate(this.Description.Get()); err != nil {
			errs = append(errs, errDescriptionFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Label) MarshalJSON() ([]byte, error) {
	type alias Label
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Label.
// A JSON null leaves the Label unchanged.
func (this *Label) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Label
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	label := Label(decoded)

	if err := label.Validate(); err != nil {
		return err
	}

	*this = label
	return nil
}

// New returns a validated Label with a new ID.
// An empty acronym or description is treated as not present.
func New(name, acronym, description string) (Label, error) {
	label := Label{
		Id:          id.New(),
		Name:        name,
		Acronym:     optioner.OfZero(acronym),
		Description: optioner.OfZero(description),
	}

	if err := label.Validate(); err != nil {
		return Label{}, err
	}

	return label, nil
}

// MustNew is like New but panics if the Label is invalid.
func MustNew(name, acronym, description string) Label {
	label, err := New(name, acronym, description)

	if err != nil {
		panic(err)
	}

	return label
}
