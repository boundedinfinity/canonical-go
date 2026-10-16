// Package name implements the Name data type.
//
// Source of truth: domain-model/person/name/name.md
package name

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	prefixValidation = validation.Slice[affix.Prefix]("prefix").Min(1)
	suffixValidation = validation.Slice[affix.Suffix]("suffix").Min(1)
	partValidation   = validation.String[string]("").Min(1).Max(30)
	firstValidation  = validation.Slice[string]("first").Min(1).Each(partValidation)
	middleValidation = validation.Slice[string]("middle").Min(1).Each(partValidation)
	lastValidation   = validation.Slice[string]("last").Min(1).Each(partValidation)

	// ErrInvalid is wrapped by every Name validation error.
	ErrInvalid  = errorer.New("invalid name")
	ErrId       = id.ErrInvalid
	ErrRequired = errorer.New("at least one of first, middle, or last is required")
	ErrPrefix   = errorer.New("prefix")
	ErrFirst    = errorer.New("first")
	ErrMiddle   = errorer.New("middle")
	ErrLast     = errorer.New("last")
	ErrSuffix   = errorer.New("suffix")

	errIdFn       = errorer.Func(ErrInvalid, id.ErrInvalid)
	errRequiredFn = errorer.Func(ErrInvalid, ErrRequired)
	errPrefixFn   = errorer.Func(ErrInvalid, ErrPrefix)
	errFirstFn    = errorer.Func(ErrInvalid, ErrFirst)
	errMiddleFn   = errorer.Func(ErrInvalid, ErrMiddle)
	errLastFn     = errorer.Func(ErrInvalid, ErrLast)
	errSuffixFn   = errorer.Func(ErrInvalid, ErrSuffix)
)

// Name is the name of a person (real, fictional or otherwise).
//
// Each part is optional: None means not present, and is omitted from JSON
// (omitzero; omitempty does not omit structs such as optioner.Option).
// A present part must contain at least one element. Prefix and Suffix hold
// the full affix.Prefix and affix.Suffix records.
type Name struct {
	Id     id.ID                           `json:"id"`
	Prefix optioner.Option[[]affix.Prefix] `json:"prefix,omitzero"`
	First  optioner.Option[[]string]       `json:"first,omitzero"`
	Middle optioner.Option[[]string]       `json:"middle,omitzero"`
	Last   optioner.Option[[]string]       `json:"last,omitzero"`
	Suffix optioner.Option[[]affix.Suffix] `json:"suffix,omitzero"`
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

	if this.Prefix.Defined() {
		if err := prefixValidation.Validate(this.Prefix.Get()); err != nil {
			errs = append(errs, errPrefixFn("%w", err))
		}

		for i, prefix := range this.Prefix.Get() {
			if err := prefix.Validate(); err != nil {
				errs = append(errs, errPrefixFn("[%d] %w", i, err))
			}
		}
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

	if this.Suffix.Defined() {
		if err := suffixValidation.Validate(this.Suffix.Get()); err != nil {
			errs = append(errs, errSuffixFn("%w", err))
		}

		for i, suffix := range this.Suffix.Get() {
			if err := suffix.Validate(); err != nil {
				errs = append(errs, errSuffixFn("[%d] %w", i, err))
			}
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
func New(prefix []affix.Prefix, first, middle, last []string, suffix []affix.Suffix) (Name, error) {
	return Create(id.New(), prefix, first, middle, last, suffix)
}

// Create returns a validated Name with the given  ID.
// Pass nil for any part that is not present.
func Create(id id.ID, prefix []affix.Prefix, first, middle, last []string, suffix []affix.Suffix) (Name, error) {
	name := Name{
		Id:     id,
		Prefix: optioner.OfSlice(prefix),
		First:  optioner.OfSlice(first),
		Middle: optioner.OfSlice(middle),
		Last:   optioner.OfSlice(last),
		Suffix: optioner.OfSlice(suffix),
	}

	if err := name.Validate(); err != nil {
		return Name{}, err
	}

	return name, nil
}

// MustNew is like New but panics if the Name is invalid.
func MustNew(prefix []affix.Prefix, first, middle, last []string, suffix []affix.Suffix) Name {
	name, err := New(prefix, first, middle, last, suffix)

	if err != nil {
		panic(err)
	}

	return name
}

func MustCreate(id id.ID, prefix []affix.Prefix, first, middle, last []string, suffix []affix.Suffix) Name {
	name, err := Create(id, prefix, first, middle, last, suffix)

	if err != nil {
		panic(err)
	}

	return name
}
