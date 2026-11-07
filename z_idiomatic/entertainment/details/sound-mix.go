package details

import "github.com/boundedinfinity/canonical-go/idiomatic/ider"

type SoundMix struct {
	Id          ider.Id `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}
