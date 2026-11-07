package label

import (
	"encoding/json"

	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func NewGroup(name, description string, labels ...*Label) Group {
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

func (_ Group) Kind() string {
	return "canonical.label.group"
}

func (this *Group) Add(labels ...*Label) {
	this.Labels.Add(labels...)
}

func (this Group) Has(label *Label) bool {
	return this.Labels.Has(label)
}

func (this Group) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind string `json:"kind"`
		Group
	}{
		Kind:  this.Kind(),
		Group: this,
	})
}

type DehydratedGroup struct {
	Kind   string    `json:"kind"`
	Id     ider.Id   `json:"id"`
	Labels []ider.Id `json:"labels"`
}

func (this Group) DehydrateMarshalJSON() ([]byte, error) {
	return json.Marshal(DehydratedGroup{
		Kind:   this.Kind(),
		Id:     this.Id,
		Labels: this.Labels.Ids(),
	})
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ByGroupId(id ider.Id) func(resource *Group) bool {
	return func(resource *Group) bool {
		return resource.Id == id
	}
}

func ByGroupTerm(term string) func(group *Group) bool {
	return func(group *Group) bool {
		return stringer.ContainsIgnoreCase(group.Name, term) ||
			stringer.ContainsIgnoreCase(group.Description, term) ||
			slicer.ContainsFunc(ByTerm(term), group.Labels...)
	}
}
