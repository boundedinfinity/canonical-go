package resource

import (
	"github.com/boundedinfinity/canonical-go/avatar/repository"
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/resource"
)

type Manager repository.Repository[resource.Resource]

func Memory(resources []*resource.Resource) Manager {
	return repository.Memory("resource", resources, resource.ById, resource.ByTerm, resource.ById)
}
