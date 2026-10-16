package label

// Source of truth: domain-model/label/label-group.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var (
	groupNameValidation        = validation.String[string]("name").Min(1).Max(20)
	groupAcronymValidation     = validation.String[string]("acronym").Min(2).Max(10)
	groupDescriptionValidation = validation.String[string]("description").Min(1).Max(250)
	groupLabelsValidation      = validation.Slice[Label]("labels").Min(1)

	// ErrInvalidLabelGroup is wrapped by every LabelGroup validation error.
	// The field sentinels (ErrId, ErrName, ErrAcronym, ErrDescription) are
	// shared with Label.
	ErrInvalidLabelGroup = errorer.New("invalid label group")
	ErrLabels            = errorer.New("labels")

	errGroupIdFn          = errorer.Func(ErrInvalidLabelGroup, id.ErrInvalid)
	errGroupNameFn        = errorer.Func(ErrInvalidLabelGroup, ErrName)
	errGroupAcronymFn     = errorer.Func(ErrInvalidLabelGroup, ErrAcronym)
	errGroupDescriptionFn = errorer.Func(ErrInvalidLabelGroup, ErrDescription)
	errGroupLabelsFn      = errorer.Func(ErrInvalidLabelGroup, ErrLabels)
)

// LabelGroup is a user defined group of Labels.
type LabelGroup struct {
	Id          id.ID                    `json:"id"`
	Name        string                   `json:"name"`
	Acronym     optioner.Option[string]  `json:"acronym,omitzero"`
	Description optioner.Option[string]  `json:"description,omitzero"`
	Labels      optioner.Option[[]Label] `json:"labels,omitzero"`
}

// Validate checks every field and returns all failures joined together.
// Each failure wraps ErrInvalidLabelGroup and the field's sentinel (ErrName,
// etc.). A Label failure also wraps the Label's own errors.
// Lengths are counted in characters (runes), not bytes.
func (this LabelGroup) Validate() error {
	var errs []error

	if err := this.Id.Validate(); err != nil {
		errs = append(errs, errGroupIdFn("%w", err))
	}

	if err := groupNameValidation.Validate(this.Name); err != nil {
		errs = append(errs, errGroupNameFn("%w", err))
	}

	if this.Acronym.Defined() {
		if err := groupAcronymValidation.Validate(this.Acronym.Get()); err != nil {
			errs = append(errs, errGroupAcronymFn("%w", err))
		}
	}

	if this.Description.Defined() {
		if err := groupDescriptionValidation.Validate(this.Description.Get()); err != nil {
			errs = append(errs, errGroupDescriptionFn("%w", err))
		}
	}

	if this.Labels.Defined() {
		if err := groupLabelsValidation.Validate(this.Labels.Get()); err != nil {
			errs = append(errs, errGroupLabelsFn("%w", err))
		}

		for i, label := range this.Labels.Get() {
			if err := label.Validate(); err != nil {
				errs = append(errs, errGroupLabelsFn("[%d] %w", i, err))
			}
		}
	}

	return errorer.Join(errs...)
}

// MarshalJSON implements json.Marshaler.
func (this LabelGroup) MarshalJSON() ([]byte, error) {
	type alias LabelGroup
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// LabelGroup. A JSON null leaves the LabelGroup unchanged. An empty array
// ("labels": []) counts as present and fails validation.
func (this *LabelGroup) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias LabelGroup
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	group := LabelGroup(decoded)

	if err := group.Validate(); err != nil {
		return err
	}

	*this = group
	return nil
}

// NewLabelGroup returns a validated LabelGroup with a new ID.
// An empty acronym or description, or a nil or empty labels slice, is treated
// as not present.
func NewLabelGroup(name, acronym, description string, labels []Label) (LabelGroup, error) {
	group := LabelGroup{
		Id:          id.New(),
		Name:        name,
		Acronym:     optioner.OfZero(acronym),
		Description: optioner.OfZero(description),
		Labels:      optioner.OfSlice(labels),
	}

	if err := group.Validate(); err != nil {
		return LabelGroup{}, err
	}

	return group, nil
}

// MustNewLabelGroup is like NewLabelGroup but panics if the LabelGroup is
// invalid.
func MustNewLabelGroup(name, acronym, description string, labels []Label) LabelGroup {
	group, err := NewLabelGroup(name, acronym, description, labels)

	if err != nil {
		panic(err)
	}

	return group
}
