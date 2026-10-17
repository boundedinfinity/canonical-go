package name

// Source of truth: domain-model/person/name/name.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	prefixIdValidation = validation.Slice[id.ID]("prefix").Min(1)
	suffixIdValidation = validation.Slice[id.ID]("suffix").Min(1)
)

// NameDto is the data transfer object for Name. Prefix and Suffix hold the
// IDs of the affix.Prefix and affix.Suffix records instead of the records.
type NameDto struct {
	Id     id.ID                     `json:"id"`
	Prefix optioner.Option[[]id.ID]  `json:"prefix,omitzero"`
	First  optioner.Option[[]string] `json:"first,omitzero"`
	Middle optioner.Option[[]string] `json:"middle,omitzero"`
	Last   optioner.Option[[]string] `json:"last,omitzero"`
	Suffix optioner.Option[[]id.ID]  `json:"suffix,omitzero"`
}

// Validate applies the Name rules, with Prefix and Suffix checked as IDs.
// Each failure wraps ErrInvalid and the field's sentinel (ErrPrefix, etc.).
func (this NameDto) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errIdFn("%w", err))
	}

	if this.First.Empty() && this.Middle.Empty() && this.Last.Empty() {
		errs = append(errs, errRequiredFn(""))
	}

	if this.Prefix.Defined() {
		if err := prefixIdValidation.Validate(this.Prefix.Get()); err != nil {
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
		if err := suffixIdValidation.Validate(this.Suffix.Get()); err != nil {
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
func (this NameDto) MarshalJSON() ([]byte, error) {
	type alias NameDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded NameDto.
// A JSON null leaves the NameDto unchanged.
func (this *NameDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias NameDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := NameDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel resolves Prefix and Suffix and returns the validated Name.
// Every resolution failure is reported and wraps resolver.ErrResolve.
// A resolver is only called when its field is present, so it may be nil
// otherwise.
func (this NameDto) ToModel(
	prefixes resolver.Resolver[id.ID, affix.Prefix],
	suffixes resolver.Resolver[id.ID, affix.Suffix],
) (Name, error) {
	if err := this.Validate(); err != nil {
		return Name{}, err
	}

	var errs []error

	prefix, err := resolver.OptionAll(prefixes, this.Prefix)

	if err != nil {
		errs = append(errs, fmt.Errorf("prefix %w", err))
	}

	suffix, err := resolver.OptionAll(suffixes, this.Suffix)

	if err != nil {
		errs = append(errs, fmt.Errorf("suffix %w", err))
	}

	if err := errorer.Join(errs...); err != nil {
		return Name{}, err
	}

	name := Name{
		Id:     this.Id,
		Prefix: prefix,
		First:  this.First,
		Middle: this.Middle,
		Last:   this.Last,
		Suffix: suffix,
	}

	if err := name.Validate(); err != nil {
		return Name{}, err
	}

	return name, nil
}

// ToDto returns the NameDto for this Name.
func (this Name) ToDto() NameDto {
	return NameDto{
		Id:     this.Id,
		Prefix: optioner.Map(this.Prefix, prefixIds),
		First:  this.First,
		Middle: this.Middle,
		Last:   this.Last,
		Suffix: optioner.Map(this.Suffix, suffixIds),
	}
}

func prefixIds(prefixes []affix.Prefix) []id.ID {
	ids := make([]id.ID, len(prefixes))

	for i, prefix := range prefixes {
		ids[i] = prefix.Id
	}

	return ids
}

func suffixIds(suffixes []affix.Suffix) []id.ID {
	ids := make([]id.ID, len(suffixes))

	for i, suffix := range suffixes {
		ids[i] = suffix.Id
	}

	return ids
}
