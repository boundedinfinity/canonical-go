package payment_card

import "github.com/boundedinfinity/canonical-go/idiomatic/bookmark"

type VirtualCard struct {
	Card
	LinkedTo  Card
	Dashboard bookmark.Bookmark
}
