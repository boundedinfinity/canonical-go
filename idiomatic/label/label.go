package label

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewLabel(name, description string) *Label {
	m := &Label{
		Id:          ider.New(),
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

func (_ Label) Kind() string {
	return "canonical.label.label"
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

type DehydratedLabel struct {
	Kind string  `json:"kind"`
	Id   ider.Id `json:"id"`
}

func (this Label) DehydrateMarshalJSON() ([]byte, error) {
	return json.Marshal(DehydratedLabel{
		Kind: this.Kind(),
		Id:   this.Id,
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ById(id ider.Id) func(resource *Label) bool {
	return func(resource *Label) bool {
		return resource.Id == id
	}
}

func ByTerm(term string) func(resource *Label) bool {
	return func(resource *Label) bool {
		return stringer.ContainsIgnoreCase(resource.Name, term) ||
			stringer.ContainsIgnoreCase(resource.Abbreviation, term) ||
			stringer.ContainsIgnoreCase(resource.Description, term)
	}
}
