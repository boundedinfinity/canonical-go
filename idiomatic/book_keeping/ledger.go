package book_keeping

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/label"
)

type Ledger struct {
	Id          ider.Id       `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Accounts    []Account     `json:"accounts"`
	Labels      []label.Label `json:"labels"`
}
