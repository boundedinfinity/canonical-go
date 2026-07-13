package ratio

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/boundedinfinity/go-commoner/idiomatic/mather"
)

type Ratio[T mather.Number] struct {
	Antecedent            float64                 `json:"antecedent"`
	Consequent            float64                 `json:"consequent"`
	AntecedentDescription optioner.Option[string] `json:"antecedent-description"`
	ConsequentDescription optioner.Option[string] `json:"consequent-description"`
}

func (this Ratio[T]) String() string {
	return fmt.Sprintf("%f:%f", this.Antecedent, this.Consequent)
}

func (this Ratio[T]) Description() optioner.Option[string] {
	if this.AntecedentDescription.Empty() || this.ConsequentDescription.Empty() {
		return optioner.None[string]()
	}

	return optioner.Some(fmt.Sprintf("%s:%s", this.AntecedentDescription.Get(), this.ConsequentDescription.Get()))
}
