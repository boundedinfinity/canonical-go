package condition

import "github.com/boundedinfinity/canonical-model/go/idiomatic/ider"

type Condition struct {
	Id          ider.Id  `json:"id"`
	Name        string   `json:"name"`
	Alias       []string `json:"aliases"`
	Description string   `json:"description"`
}
