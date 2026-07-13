package ratio

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
)

type Ratio[T mather.Number] struct {
	Antecedent            float64 `json:"antecedent"`
	Consequent            float64 `json:"consequent"`
	AntecedentDescription string  `json:"antecedent-description"`
	ConsequentDescription string  `json:"consequent-description"`
}

func (this Ratio[T]) String() string {
	return fmt.Sprintf("%f:%f", this.Antecedent, this.Consequent)
}

func (this Ratio[T]) Description() string {
	if this.AntecedentDescription == "" || this.ConsequentDescription == "" {
		return ""
	}

	return fmt.Sprintf("%s:%s", this.AntecedentDescription, this.ConsequentDescription)
}
