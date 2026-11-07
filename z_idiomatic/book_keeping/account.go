package book_keeping

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/label"
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/time"
)

type Account struct {
	Id           ider.Id            `json:"id"`
	Name         string             `json:"name"`
	Kind         Kind               `json:"kind"`
	Description  string             `json:"description"`
	Transactions []TransactionModel `json:"transactions"`
	Parent       *Account           `json:"parent"`
	Children     []*Account         `json:"children"`
	Labels       []label.Label      `json:"labels"`
	Vendor       *Vendor            `json:"vendor"`
	Client       *Client            `json:"client"`
	CreatedAt    time.Timestamp     `json:"created-at"`
}
