package access_control_list

import (
	"github.com/boundedinfinity/canonical-go/avatar/repository"
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/access_control_list"
)

type Manager repository.Repository[access_control_list.List]

func Memory(lists []*access_control_list.List) Manager {
	return repository.Memory(
		"access control list",
		lists,
		access_control_list.ById,
		access_control_list.ByTerm,
		access_control_list.ByMatch,
	)
}
