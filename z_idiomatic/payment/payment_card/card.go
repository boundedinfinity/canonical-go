package payment_card

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/location/postal_code"
	"github.com/boundedinfinity/canonical-go/idiomatic/person"
	"github.com/boundedinfinity/canonical-go/idiomatic/person/name"
	"github.com/boundedinfinity/rfc3339date"
)

type Card struct {
	Id               ider.Id                 `json:"id"`
	Name             string                  `json:"name"`
	Person           person.Person           `json:"name-on-card"`
	NameFormatter    name.Formatter          `json:"name-formatter"`
	Kind             Kind                    `json:"kind"`
	Number           Number                  `json:"number"`
	IssueDate        rfc3339date.Rfc3339Date `json:"issue-date"`
	ActivationDate   rfc3339date.Rfc3339Date `json:"activation-date"`
	ExpirationDate   rfc3339date.Rfc3339Date `json:"expiration-date"`
	CardSecurityCode CardVerificationValue   `json:"card-security-code"`
	PersonalIdNumber string                  `json:"personal-id-number"`
	PostalCode       postal_code.PostalCode  `json:"postal-code"`
}

func (this Card) GetId() ider.Id {
	return this.Id
}

func (this Card) GetName() string {
	return this.Name
}

// SecurityCode is typically a 3 or 4 digit number printed on the back of the card
// used for "card not presnt" scenarios and entering a Personal Idenfiication
// Number (PIN) is possible.
func (this Card) SecurityCode() CardVerificationValue {
	return this.CardSecurityCode
}

// Ccv A common typo or alternate shorthand sometimes seen on payment forms, meaning the same thing
// as the [Card.SecurityCode].
func (this Card) Ccv() CardVerificationValue {
	return this.CardSecurityCode
}

// Csc is the Card Security Code.  This is another name for the [Card.SecurityCode].
func (this Card) Csc() CardVerificationValue {
	return this.CardSecurityCode
}

// Cvc is the Card Verification Code.  This is another name for the [Card.SecurityCode].
func (this Card) Cvc() CardVerificationValue {
	return this.CardSecurityCode
}

// Cid is the card Identifier.  This is another name for the [Card.SecurityCode].
func (this Card) Cid() CardVerificationValue {
	return this.CardSecurityCode
}

// Pin another name for the [Card.PersonalIdNumber]
func (this Card) Pin() string {
	return this.PersonalIdNumber
}
