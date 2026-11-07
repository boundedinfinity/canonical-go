package movie

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/location/country"
)

type TheatricalRelease struct {
	// ReleaseDate time.Date       `json:"release-date"`
	Country country.Country `json:"countries"`
}
