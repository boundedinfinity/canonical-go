package temperature

import "github.com/boundedinfinity/go-commoner/idiomatic/stringer"

type Scale string

func (s Scale) String() string {
	return string(s)
}

var Scales = scales{
	Celsius:    "measurement.scale.celsius",
	Kelvin:     "measurement.scale.kelvin",
	Fahrenheit: "measurement.scale.fahrenheit",
	Rankine:    "measurement.scale.rankine",
	Romer:      "measurement.scale.romer",
	Newton:     "measurement.scale.newton",
	Delisle:    "measurement.scale.delisle",
	Reaumur:    "measurement.scale.reaumur",
	GasMark:    "measurement.scale.gas-mark",
	Leiden:     "measurement.scale.leiden",
	Wedgwood:   "measurement.scale.wedgwood",
}

type scales struct {
	Celsius    Scale
	Kelvin     Scale
	Fahrenheit Scale
	Rankine    Scale
	Romer      Scale
	Newton     Scale
	Delisle    Scale
	Reaumur    Scale
	GasMark    Scale
	Leiden     Scale
	Wedgwood   Scale
}

func clean(v Scale) string {
	return stringer.ReplaceAll(v, "measurement.scale.", "")
}
