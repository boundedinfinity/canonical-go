package book_keeping

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/business"
	"github.com/boundedinfinity/canonical-go/idiomatic/contact"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/label"
)

type VendorModel struct {
	Id          ider.Id           `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Contacts    contact.Contact   `json:"contacts"`
	Business    business.Business `json:"business"`
	Labels      []label.Label     `json:"labels"`
}
