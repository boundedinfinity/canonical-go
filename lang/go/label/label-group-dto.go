package label

// Source of truth: domain-model/label/label-group.md,
// domain-model/global-notes/dto.md

import (
	"encoding/json"
	"fmt"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/boundedinfinity/go-commoner/errorer"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

var groupLabelIdsValidation = validation.Slice[id.ID]("labels").Min(1)

// LabelGroupDto is the data transfer object for LabelGroup. Labels holds the
// IDs of the Labels instead of the Labels themselves.
type LabelGroupDto struct {
	Id          id.ID                    `json:"id"`
	Name        string                   `json:"name"`
	Acronym     optioner.Option[string]  `json:"acronym,omitzero"`
	Description optioner.Option[string]  `json:"description,omitzero"`
	Labels      optioner.Option[[]id.ID] `json:"labels,omitzero"`
}

// Validate applies the LabelGroup rules, with Labels checked as IDs.
// Each failure wraps ErrInvalidLabelGroup and the field's sentinel (ErrName,
// etc.).
func (this LabelGroupDto) Validate() error {
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
		if err := groupLabelIdsValidation.Validate(this.Labels.Get()); err != nil {
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
func (this LabelGroupDto) MarshalJSON() ([]byte, error) {
	type alias LabelGroupDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded
// LabelGroupDto. A JSON null leaves the LabelGroupDto unchanged.
func (this *LabelGroupDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias LabelGroupDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := LabelGroupDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel resolves Labels and returns the validated LabelGroup.
// Every resolution failure is reported and wraps resolver.ErrResolve.
// The resolver is only called when Labels is present, so it may be nil
// otherwise.
func (this LabelGroupDto) ToModel(labels resolver.Resolver[id.ID, Label]) (LabelGroup, error) {
	if err := this.Validate(); err != nil {
		return LabelGroup{}, err
	}

	resolved, err := resolver.OptionAll(labels, this.Labels)

	if err != nil {
		return LabelGroup{}, fmt.Errorf("labels %w", err)
	}

	group := LabelGroup{
		Id:          this.Id,
		Name:        this.Name,
		Acronym:     this.Acronym,
		Description: this.Description,
		Labels:      resolved,
	}

	if err := group.Validate(); err != nil {
		return LabelGroup{}, err
	}

	return group, nil
}

// ToDto returns the LabelGroupDto for this LabelGroup.
func (this LabelGroup) ToDto() LabelGroupDto {
	return LabelGroupDto{
		Id:          this.Id,
		Name:        this.Name,
		Acronym:     this.Acronym,
		Description: this.Description,
		Labels:      optioner.Map(this.Labels, labelIds),
	}
}

func labelIds(labels []Label) []id.ID {
	ids := make([]id.ID, len(labels))

	for i, label := range labels {
		ids[i] = label.Id
	}

	return ids
}
