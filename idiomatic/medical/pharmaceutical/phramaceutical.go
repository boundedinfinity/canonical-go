package pharmaceutical

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/medical/condition"
	"github.com/boundedinfinity/canonical-go/idiomatic/modeller"
)

type Phramaceutical struct {
	Id          ider.Id               `json:"id"`
	Name        string                `json:"name"`
	Alias       []string              `json:"aliases"`
	Description string                `json:"description"`
	Treats      []condition.Condition `json:"treats"`
}

func (this Phramaceutical) Kind() string {
	return "canonical.medical.pharmacetical"
}

func (this Phramaceutical) MarshalJSON() ([]byte, error) {
	return modeller.CanonicalJsonMarshal(this)
}
