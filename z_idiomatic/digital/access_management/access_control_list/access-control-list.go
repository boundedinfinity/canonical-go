package access_control_list

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/permission"
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/resource"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

type List struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Items       []Item  `json:"items,omitempty"`
}

type Item struct {
	Resource   *resource.Resource     `json:"resource"`
	Permission *permission.Permission `json:"permission"`
}

func ById(id ider.Id) func(*List) bool {
	return func(item *List) bool {
		return item.Id == id
	}
}

func ByTerm(term string) func(*List) bool {
	return func(item *List) bool {
		return stringer.ContainsIgnoreCase(item.Name, term) ||
			stringer.ContainsIgnoreCase(item.Description, term)
	}
}

func ByMatch(id ider.Id) func(*List) bool {
	return func(list *List) bool {
		if list.Id == id {
			return true
		}

		for _, item := range list.Items {
			if item.Resource.Id == id || item.Permission.Id == id {
				return true
			}
		}

		return false
	}
}
