package details

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/location"
)

type FilmingLocation struct {
	Id       ider.Id `json:"id"`
	Location location.Location
}
