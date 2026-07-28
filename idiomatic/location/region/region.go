package region

import "github.com/boundedinfinity/canonical-go/idiomatic/location/country"

type Region struct {
	Name         string
	Abbreviation string
	Contries     []country.Country
}
