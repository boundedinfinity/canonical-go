package label

// Source of truth: domain-model/label/label.md, domain-model/global-notes/dto.md

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
)

// LabelDto is the data transfer object for Label. Label holds no other Top
// Level Models, so LabelDto has the same fields.
type LabelDto struct {
	Id          id.ID                     `json:"id"`
	Name        string                    `json:"name"`
	Acronyms    optioner.Option[[]string] `json:"acronyms,omitzero"`
	Description optioner.Option[string]   `json:"description,omitzero"`
}

// Validate applies the Label rules; see Label.Validate.
func (this LabelDto) Validate() error {
	return Label(this).Validate()
}

// MarshalJSON implements json.Marshaler.
func (this LabelDto) MarshalJSON() ([]byte, error) {
	type alias LabelDto
	return json.Marshal(alias(this))
}

// UnmarshalJSON implements json.Unmarshaler and validates the decoded LabelDto.
// A JSON null leaves the LabelDto unchanged.
func (this *LabelDto) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	type alias LabelDto
	var decoded alias

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	dto := LabelDto(decoded)

	if err := dto.Validate(); err != nil {
		return err
	}

	*this = dto
	return nil
}

// ToModel returns the validated Label for this LabelDto.
func (this LabelDto) ToModel() (Label, error) {
	label := Label(this)

	if err := label.Validate(); err != nil {
		return Label{}, err
	}

	return label, nil
}

// ToDto returns the LabelDto for this Label.
func (this Label) ToDto() LabelDto {
	return LabelDto(this)
}
