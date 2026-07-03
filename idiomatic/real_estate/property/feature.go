package property

import "github.com/boundedinfinity/canonical-go/idiomatic/ider"

type Feature struct {
	Id   ider.Id `json:"id"`
	Name string  `json:"name"`
}
