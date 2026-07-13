package label

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewGroup(name, description string, labels ...Label) Group {
	m := Group{
		Name:        name,
		Description: description,
	}

	m.Add(labels...)

	return m
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Group struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Labels      Labels  `json:"labels"`
}

func (this *Group) Add(labels ...Label) {
	this.Labels.Add(labels...)
}

func (this Group) Has(label Label) bool {
	return this.Labels.Has(label)
}

func (this Group) MarshalJSON() ([]byte, error) {
	v := struct {
		Kind string `json:"kind"`
		Group
	}{
		Kind:  "canonical.label.group",
		Group: this,
	}

	return json.Marshal(v)
}
