package perscription

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/business"
	"github.com/boundedinfinity/canonical-go/idiomatic/currency"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/medical/amount"
	"github.com/boundedinfinity/canonical-go/idiomatic/medical/pharmaceutical"
	"github.com/boundedinfinity/canonical-go/idiomatic/modeller"
	"github.com/boundedinfinity/canonical-go/idiomatic/person"
	"github.com/boundedinfinity/canonical-go/idiomatic/phone"
	"github.com/boundedinfinity/rfc3339date"
)

type PerscriptionModel struct {
	Id             ider.Id                       `json:"id"`
	PerscriptionId string                        `json:"perscription-id"`
	Patient        person.Person                 `json:"patient"`
	Physician      person.Person                 `json:"physician"`
	Pharmachy      business.Business             `json:"pharmachy"`
	Directions     []string                      `json:"directions"`
	Phramaceutical pharmaceutical.Phramaceutical `json:"phramaceutical"`
	Count          int                           `json:"count"`
	Phone          phone.Phone                   `json:"phone"`
	PickedUpDate   rfc3339date.Rfc3339Date       `json:"picked-up-date"`
	PromisedDate   rfc3339date.Rfc3339Date       `json:"promised-date"`
	RefillCount    int                           `json:"refill-count"`
	RefillBy       rfc3339date.Rfc3339Date       `json:"refill-by"`
	FormFactor     FormFactor                    `json:"form-factor"`
	Amount         amount.Amount                 `json:"amount"`
	Price          currency.Amount               `json:"price"`
}

func (this PerscriptionModel) Kind() string {
	return "canonical.medical.prescription"
}

func (this PerscriptionModel) MarshalJSON() ([]byte, error) {
	return modeller.CanonicalJsonMarshal(this)
}
