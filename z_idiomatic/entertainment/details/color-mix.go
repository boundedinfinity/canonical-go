package details

import "github.com/boundedinfinity/canonical-go/idiomatic/ider"

type ColorMix struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}
