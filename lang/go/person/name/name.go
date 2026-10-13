// Package name implements the Name data type.
//
// Source of truth: domain-model/person/name/name.md
package name

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	partValidation   = validation.String[string]("").Min(1).Max(30)
	firstValidation  = validation.Slice[string]("first").Min(1).Each(partValidation)
	middleValidation = validation.Slice[string]("middle").Min(1).Each(partValidation)
	lastValidation   = validation.Slice[string]("last").Min(1).Each(partValidation)

	// ErrInvalid is wrapped by every Name validation error.
	ErrInvalid  = errorer.New("invalid name")
	ErrId       = id.ErrInvalid
	ErrRequired = errorer.New("at least one of first, middle, or last is required")
	ErrFirst    = errorer.New("first")
	ErrMiddle   = errorer.New("middle")
	ErrLast     = errorer.New("last")

	errIdFn       = errorer.Func(ErrInvalid, id.ErrInvalid)
	errRequiredFn = errorer.Func(ErrInvalid, ErrRequired)
	errFirstFn    = errorer.Func(ErrInvalid, ErrFirst)
	errMiddleFn   = errorer.Func(ErrInvalid, ErrMiddle)
	errLastFn     = errorer.Func(ErrInvalid, ErrLast)
)

// Name is the name of a person (real, fictional or otherwise).
//
// Each part is optional: a nil slice means not present. A present part
// must contain at least one name.
type Name struct {
	Id     id.ID                     `json:"id"`
	First  optioner.Option[[]string] `json:"first,omitempty"`
	Middle optioner.Option[[]string] `json:"middle,omitempty"`
	Last   optioner.Option[[]string] `json:"last,omitempty"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalid and the field's sentinel (ErrFirst, etc.).
// Lengths are counted in characters (runes), not bytes.
func (this Name) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errIdFn("%w", err))
	}

	if this.First.Empty() && this.Middle.Empty() && this.Last.Empty() {
		errs = append(errs, errRequiredFn(""))
	}

	if this.First.Defined() {
		if err := firstValidation.Validate(this.First.Get()); err != nil {
			errs = append(errs, errFirstFn("%w", err))
		}
	}

	if this.Middle.Defined() {
		if err := middleValidation.Validate(this.Middle.Get()); err != nil {
			errs = append(errs, errMiddleFn("%w", err))
		}
	}

	if this.Last.Defined() {
		if err := lastValidation.Validate(this.Last.Get()); err != nil {
			errs = append(errs, errLastFn("%w", err))
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this Name) MarshalJSON() ([]byte, error) {
	type alias Name
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded Name.
// A JSON null leaves the Name unchanged. An empty array (e.g. "first": [])
// counts as present and fails validation; omit the field or use null instead.
func (this *Name) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias Name
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	name := Name(decoded)

	if err := name.Validate(); err != nil {
		return err
	}

	*this = name
	return nil
}

// New returns a validated Name with a new ID.
// Pass nil for any part that is not present.
func New(first, middle, last []string) (Name, error) {
	return Create(id.New(), first, middle, last)
}

// Create returns a validated Name with the given  ID.
// Pass nil for any part that is not present.
func Create(id id.ID, first, middle, last []string) (Name, error) {
	name := Name{
		Id:     id,
		First:  optioner.OfSlice(first),
		Middle: optioner.OfSlice(middle),
		Last:   optioner.OfSlice(last),
	}

	if err := name.Validate(); err != nil {
		return Name{}, err
	}

	return name, nil
}

// MustNew is like New but panics if the Name is invalid.
func MustNew(first, middle, last []string) Name {
	name, err := New(first, middle, last)

	if err != nil {
		panic(err)
	}

	return name
}

func MustCreate(id id.ID, first, middle, last []string) Name {
	name, err := Create(id, first, middle, last)

	if err != nil {
		panic(err)
	}

	return name
}
