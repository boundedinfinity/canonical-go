package authenticator_application

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/digital/bookmark"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Application struct {
	Id        ider.Id             `json:"id"`
	Name      string              `json:"name"`
	Bookmarks []bookmark.Bookmark `json:"bookmarks"`
}
