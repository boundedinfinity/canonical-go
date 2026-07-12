package permission

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
)

type PermissionAble interface {
	AccessManagementPermissionName() string
}

type Permission struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
}

func (this Permission) Matches(permission *Permission) bool {
	return ById(this.Id)(permission)
}

func (this Permission) Copy() *Permission {
	return &Permission{
		Id:          ider.New(),
		Name:        this.Name,
		Description: this.Description,
	}
}

func ById(id ider.Id) func(permission *Permission) bool {
	return func(permission *Permission) bool {
		return permission.Id == id
	}
}

func ByTerm(term string) func(permission *Permission) bool {
	return func(permission *Permission) bool {
		return stringer.ContainsIgnoreCase(permission.Name, term) ||
			stringer.ContainsIgnoreCase(permission.Description, term)
	}
}
