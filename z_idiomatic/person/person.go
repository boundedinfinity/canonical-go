package person

import (
	"github.com/boundedinfinity/canonical-go/idiomatic/canonical"
	"github.com/boundedinfinity/canonical-go/idiomatic/ider"
	"github.com/boundedinfinity/canonical-go/idiomatic/label"
	"github.com/boundedinfinity/canonical-go/idiomatic/location/country"
	"github.com/boundedinfinity/canonical-go/idiomatic/measurement/time"
	"github.com/boundedinfinity/canonical-go/idiomatic/person/name"
	"github.com/boundedinfinity/canonical-go/idiomatic/professional"
	"github.com/boundedinfinity/rfc3339date"
)

type Person struct {
	Id                      ider.Id                             `json:"id"`
	LegalName               name.Name                           `json:"legal-name"`
	Aliases                 []name.Name                         `json:"aliases"`
	ProfessionalCredentials professional.Credentials            `json:"professional-credentials"`
	DateOfBirth             []rfc3339date.Rfc3339DateTime       `json:"date-of-birth"`
	DateOfDeath             []rfc3339date.Rfc3339DateTime       `json:"date-of-death"`
	CountryOfBirth          canonical.OneToOne[country.Country] `json:"country-of-birth"`
	Labels                  canonical.OneToMany[label.Label]    `json:"labels"`
}

func (p Person) String() string {
	return p.LegalName.String()
}

type personLoader struct {
	dto                     PersonDto
	LegalName               name.Name                     `json:"legal-name"`
	Aliases                 []name.Name                   `json:"aliases"`
	ProfessionalCredentials professional.Credentials      `json:"professional-credentials"`
	DateOfBirth             []rfc3339date.Rfc3339DateTime `json:"date-of-birth"`
	DateOfDeath             []rfc3339date.Rfc3339DateTime `json:"date-of-death"`
	countryOfBirthId        ider.Id
	CountryOfBirth          canonical.OneToOne[country.Country] `json:"country-of-birth"`
	Labels                  canonical.OneToMany[label.Label]    `json:"labels"`
}

type PersonDto struct {
	Id                      ider.Id     `json:"id"`
	LegalName               []ider.Id   `json:"legal-name"`
	Aliases                 []ider.Id   `json:"aliases"`
	ProfessionalCredentials ider.Id     `json:"professional-credentials"`
	DateOfBirth             []time.Date `json:"date-of-birth"`
	DateOfDeath             []time.Date `json:"date-of-death"`
	CountryOfBirth          ider.Id     `json:"country-of-birth"`
	Labels                  []ider.Id   `json:"labels"`
}
