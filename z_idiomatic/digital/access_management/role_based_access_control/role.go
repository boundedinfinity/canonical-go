package role_based_access_control

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/access_management/subject"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Role struct {
	ID          ider.Id           `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Subjects    []subject.Subject `json:"subjects,omitempty"`
}
