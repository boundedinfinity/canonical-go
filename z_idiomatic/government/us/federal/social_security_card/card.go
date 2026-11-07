package social_security_card

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/time"
	"github.com/boundedinfinity/canonical-go/idiomatic/person/name"
)

type Card struct {
	Number    Number
	Name      name.Name
	Signature name.Name
	Date      time.Date
}
