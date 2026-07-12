package permission

import (
	"github.com/boundedinfinity/canonical-go/avatar/repository"
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/permission"
)

type Manager repository.Repository[permission.Permission]

func Memory(permissions []*permission.Permission) Manager {
	return repository.Memory("permission", permissions, permission.ById, permission.ByTerm, permission.ById)
}
