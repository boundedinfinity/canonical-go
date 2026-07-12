package role_based_access_control

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/permission"
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/resource"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type PolicyModel struct {
	ID          ider.Id                 `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Resources   []resource.Resource     `json:"resources,omitempty"`
	Roles       []Role                  `json:"roles,omitempty"`
	Permissions []permission.Permission `json:"permissions,omitempty"`
}
