package power_supply

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/amperage"
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/temperature"
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/voltage"
)

type PowerSupply struct {
	Voltage            voltage.Voltage         `json:"voltage"`
	Current            amperage.Amperage       `json:"current"`
	WorkingTemperature temperature.Temperature `json:"working-temperature"`
}
