package role

import (
	"image"

	"github.com/boundedinfinity/canonical-go/idiomatic/contact"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Creator struct {
	Id       ider.Id         `json:"id"`
	Contact  contact.Contact `json:"contact"`
	Image    []image.Image   `json:"image"`
	Synopsis string          `json:"synopsis"`
}
