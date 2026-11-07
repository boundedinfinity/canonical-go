package role

import (
	"image"

	"github.com/boundedinfinity/canonical-go/idiomatic/contact"
	"github.com/boundedinfinity/canonical-go/idiomatic/entertainment/fiction"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Actor struct {
	Id         ider.Id             `json:"id"`
	Contact    contact.Contact     `json:"contact"`
	Image      []image.Image       `json:"image"`
	Synopsis   string              `json:"synopsis"`
	Characters []fiction.Character `json:"characters"`
}
