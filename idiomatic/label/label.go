package label

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewLabel(name, description string) Label {
	m := Label{
		Name:        name,
		Description: description,
	}

	return m
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Label struct {
	Id           ider.Id `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Abbreviation string  `json:"abbreviation,omitempty"`
}

func (this Label) MarshalJSON() ([]byte, error) {
	v := struct {
		Kind string `json:"kind"`
		Label
	}{
		Kind:  "canonical.label.label",
		Label: this,
	}

	return json.Marshal(v)
}
