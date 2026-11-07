package resource

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

type ResourceAble interface {
	AccessManagementResourceName() string
}

// Resource represents a resource that can be controlled by an access management system.
type Resource struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
}

func (this Resource) Matches(resource *Resource) bool {
	return ById(this.Id)(resource) || ByTerm(this.Name)(resource)
}

func ById(id ider.Id) func(resource *Resource) bool {
	return func(resource *Resource) bool {
		return resource.Id == id
	}
}

func ByTerm(term string) func(resource *Resource) bool {
	return func(resource *Resource) bool {
		return stringer.ContainsIgnoreCase(resource.Name, term) ||
			stringer.ContainsIgnoreCase(resource.Description, term)
	}
}
