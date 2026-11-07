package note

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/label"
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/time"
)

type Note struct {
	Id        ider.Id        `json:"id"`
	Text      string         `json:"text"`
	CreatedAt time.Timestamp `json:"created-at"`
	Labels    []label.Label  `json:"labels"`
}
