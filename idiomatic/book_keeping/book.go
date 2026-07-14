package book_keeping

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
)

type Book struct {
	Id          ider.Id  `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ledgers     []Ledger `json:"ledgers"`
}
