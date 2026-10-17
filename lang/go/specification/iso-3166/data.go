package iso3166

// Source of truth: the ## Lookup and ## Data sections of
// domain-model/specification/iso-3166/country-code.md and
// domain-model/specification/iso-3166/subdivision.md

//go:generate cp ../../../../domain-model/specification/iso-3166/country-code.json ../../../../domain-model/specification/iso-3166/subdivision.json .

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/boundedinfinity/canonical-go/lang/go/resolver"
	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	//go:embed country-code.json
	countryCodeJSON []byte

	//go:embed subdivision.json
	subdivisionJSON []byte

	// ErrNotFound is returned by Resolve when no model has the given key.
	ErrNotFound = errorer.New("not found")

	// CountryCodes holds every ISO 3166-1 country in the embedded data.
	CountryCodes = mustLoadCountryCodes()

	// Subdivisions holds every ISO 3166-2 subdivision in the embedded data.
	Subdivisions = mustLoadSubdivisions(CountryCodes)
)

// CountryCodeList is a set of CountryCodes with the lookups from the spec.
// Every lookup ignores case.
type CountryCodeList struct {
	all       []CountryCode
	byName    map[string]int
	byAlpha2  map[string]int
	byAlpha3  map[string]int
	byNumeric map[string]int
}

// All returns every country, sorted by Alpha2.
func (this CountryCodeList) All() []CountryCode {
	return slices.Clone(this.all)
}

// Get returns the single country whose Name, Alpha2, Alpha3, or Numeric equals
// term, checked in that order. A numeric term also matches without its leading
// zeros, so "4" and "004" both find Afghanistan.
func (this CountryCodeList) Get(term string) (CountryCode, bool) {
	lc := strings.ToLower(term)

	for _, index := range []map[string]int{this.byName, this.byAlpha2, this.byAlpha3} {
		if i, ok := index[lc]; ok {
			return this.all[i], true
		}
	}

	if i, ok := this.byNumeric[strings.TrimLeft(lc, "0")]; ok {
		return this.all[i], true
	}

	return CountryCode{}, false
}

// Find returns the result of Get if there is one. Otherwise it returns every
// country whose Name, Alpha2, Alpha3, or Numeric contains term, in Alpha2
// order. It may return several countries, or none.
func (this CountryCodeList) Find(term string) []CountryCode {
	if country, ok := this.Get(term); ok {
		return []CountryCode{country}
	}

	lc := strings.ToLower(term)
	var found []CountryCode

	for _, country := range this.all {
		if strings.Contains(strings.ToLower(country.Name), lc) ||
			strings.Contains(strings.ToLower(country.Alpha2), lc) ||
			strings.Contains(strings.ToLower(country.Alpha3), lc) ||
			strings.Contains(country.Numeric, lc) {
			found = append(found, country)
		}
	}

	return found
}

// Resolve returns the country with the given Alpha2, so a CountryCodeList is a
// resolver.Resolver[string, CountryCode]. It fails with ErrNotFound.
func (this CountryCodeList) Resolve(alpha2 string) (CountryCode, error) {
	if i, ok := this.byAlpha2[strings.ToLower(alpha2)]; ok {
		return this.all[i], nil
	}

	return CountryCode{}, fmt.Errorf("country code %q %w", alpha2, ErrNotFound)
}

// SubdivisionList is a set of Subdivisions with the lookups from the spec.
// Every lookup ignores case.
type SubdivisionList struct {
	all       []Subdivision
	byCode    map[string]int
	byName    map[string][]int
	byCountry map[string][]int
}

// All returns every subdivision, sorted by Code.
func (this SubdivisionList) All() []Subdivision {
	return slices.Clone(this.all)
}

// Get returns the subdivision whose Code equals term. The term also matches
// without the hyphen, so "USAL" finds "US-AL".
func (this SubdivisionList) Get(term string) (Subdivision, bool) {
	if i, ok := this.byCode[subdivisionCodeKey(term)]; ok {
		return this.all[i], true
	}

	return Subdivision{}, false
}

// FindByName returns every subdivision whose Name equals name, in Code order.
// Names are not unique, so it may return several subdivisions, or none.
func (this SubdivisionList) FindByName(name string) []Subdivision {
	return this.pick(this.byName[strings.ToLower(name)])
}

// ByCountry returns every subdivision of country, in Code order.
func (this SubdivisionList) ByCountry(country CountryCode) []Subdivision {
	return this.pick(this.byCountry[strings.ToLower(country.Alpha2)])
}

// Resolve returns the subdivision with the given Code, so a SubdivisionList is
// a resolver.Resolver[string, Subdivision]. It fails with ErrNotFound.
func (this SubdivisionList) Resolve(code string) (Subdivision, error) {
	if subdivision, ok := this.Get(code); ok {
		return subdivision, nil
	}

	return Subdivision{}, fmt.Errorf("subdivision %q %w", code, ErrNotFound)
}

func (this SubdivisionList) pick(indexes []int) []Subdivision {
	if len(indexes) == 0 {
		return nil
	}

	subdivisions := make([]Subdivision, len(indexes))

	for i, index := range indexes {
		subdivisions[i] = this.all[index]
	}

	return subdivisions
}

// subdivisionCodeKey is the lookup key for a code: lowercase, without the
// hyphen.
func subdivisionCodeKey(code string) string {
	return strings.ReplaceAll(strings.ToLower(code), "-", "")
}

// mustLoadCountryCodes decodes the embedded country DTOs (validating each one)
// and indexes them. The data is checked by the package tests, so a failure here
// is a bug in the embedded data.
func mustLoadCountryCodes() CountryCodeList {
	var dtos []CountryCodeDto

	if err := json.Unmarshal(countryCodeJSON, &dtos); err != nil {
		panic(fmt.Errorf("iso3166: loading country-code.json: %w", err))
	}

	list := CountryCodeList{
		all:       make([]CountryCode, len(dtos)),
		byName:    map[string]int{},
		byAlpha2:  map[string]int{},
		byAlpha3:  map[string]int{},
		byNumeric: map[string]int{},
	}

	for i, dto := range dtos {
		country, err := dto.ToModel()

		if err != nil {
			panic(fmt.Errorf("iso3166: loading country-code.json: %w", err))
		}

		list.all[i] = country
		list.byName[strings.ToLower(country.Name)] = i
		list.byAlpha2[strings.ToLower(country.Alpha2)] = i
		list.byAlpha3[strings.ToLower(country.Alpha3)] = i
		list.byNumeric[strings.TrimLeft(country.Numeric, "0")] = i
	}

	return list
}

// mustLoadSubdivisions decodes the embedded subdivision DTOs (validating each
// one), resolves each Country from countries and each Parent from the same
// data, and indexes them.
func mustLoadSubdivisions(countries CountryCodeList) SubdivisionList {
	var dtos []SubdivisionDto

	if err := json.Unmarshal(subdivisionJSON, &dtos); err != nil {
		panic(fmt.Errorf("iso3166: loading subdivision.json: %w", err))
	}

	dtosByCode := make(map[string]SubdivisionDto, len(dtos))

	for _, dto := range dtos {
		dtosByCode[dto.Code] = dto
	}

	// parents converts a parent's DTO on demand, so parents may appear in any
	// order. inProgress stops a parent chain that loops back on itself.
	loaded := map[string]Subdivision{}
	inProgress := map[string]bool{}
	var parents resolver.Func[string, Subdivision]

	parents = func(code string) (Subdivision, error) {
		if subdivision, ok := loaded[code]; ok {
			return subdivision, nil
		}

		dto, ok := dtosByCode[code]

		if !ok {
			return Subdivision{}, fmt.Errorf("subdivision %q %w", code, ErrNotFound)
		}

		if inProgress[code] {
			return Subdivision{}, errorer.Func(ErrInvalidSubdivision, ErrParentCycle)("%q", code)
		}

		inProgress[code] = true
		subdivision, err := dto.ToModel(countries, parents)
		delete(inProgress, code)

		if err != nil {
			return Subdivision{}, err
		}

		loaded[code] = subdivision
		return subdivision, nil
	}

	list := SubdivisionList{
		all:       make([]Subdivision, len(dtos)),
		byCode:    map[string]int{},
		byName:    map[string][]int{},
		byCountry: map[string][]int{},
	}

	for i, dto := range dtos {
		subdivision, err := parents(dto.Code)

		if err != nil {
			panic(fmt.Errorf("iso3166: loading subdivision.json: %w", err))
		}

		list.all[i] = subdivision
		list.byCode[subdivisionCodeKey(subdivision.Code)] = i

		name := strings.ToLower(subdivision.Name)
		list.byName[name] = append(list.byName[name], i)

		country := strings.ToLower(subdivision.Country.Alpha2)
		list.byCountry[country] = append(list.byCountry[country], i)
	}

	return list
}
