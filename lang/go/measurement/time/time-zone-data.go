package time

// Source of truth: the ## Lookup and ## Data sections of
// domain-model/measurement/time/time-zone.md

//go:generate cp ../../../../domain-model/measurement/time/time-zone.json .

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/boundedinfinity/go-commoner/errorer"
)

var (
	//go:embed time-zone.json
	timeZoneJSON []byte

	// ErrNotFound is returned by Resolve when no model has the given key.
	ErrNotFound = errorer.New("not found")

	// TimeZones holds every time zone in the embedded data.
	TimeZones = mustLoadTimeZones()
)

// TimeZoneList is a set of TimeZones with the lookups from the spec.
// Every lookup ignores case.
type TimeZoneList struct {
	all          []TimeZone
	byIdentifier map[string]int
	byName       map[string]int
}

// All returns every time zone, sorted by Identifier.
func (this TimeZoneList) All() []TimeZone {
	return slices.Clone(this.all)
}

// Get returns the single time zone whose Identifier or Name equals term,
// checked in that order.
func (this TimeZoneList) Get(term string) (TimeZone, bool) {
	lc := strings.ToLower(term)

	for _, index := range []map[string]int{this.byIdentifier, this.byName} {
		if i, ok := index[lc]; ok {
			return this.all[i], true
		}
	}

	return TimeZone{}, false
}

// Find returns the result of Get if there is one. Otherwise it returns every
// time zone whose Abbreviation or UtcOffset equals term, in Identifier order.
// It may return several time zones, or none.
func (this TimeZoneList) Find(term string) []TimeZone {
	if timeZone, ok := this.Get(term); ok {
		return []TimeZone{timeZone}
	}

	var found []TimeZone

	for _, timeZone := range this.all {
		if strings.EqualFold(timeZone.Abbreviation, term) || strings.EqualFold(timeZone.UtcOffset, term) {
			found = append(found, timeZone)
		}
	}

	return found
}

// Resolve returns the time zone with the given Identifier, so a TimeZoneList
// is a resolver.Resolver[string, TimeZone]. It fails with ErrNotFound.
func (this TimeZoneList) Resolve(identifier string) (TimeZone, error) {
	if i, ok := this.byIdentifier[strings.ToLower(identifier)]; ok {
		return this.all[i], nil
	}

	return TimeZone{}, fmt.Errorf("time zone %q %w", identifier, ErrNotFound)
}

// mustLoadTimeZones decodes the embedded time zone DTOs (validating each one)
// and indexes them. The data is checked by the package tests, so a failure
// here is a bug in the embedded data.
func mustLoadTimeZones() TimeZoneList {
	var dtos []TimeZoneDto

	if err := json.Unmarshal(timeZoneJSON, &dtos); err != nil {
		panic(fmt.Errorf("time: loading time-zone.json: %w", err))
	}

	list := TimeZoneList{
		all:          make([]TimeZone, len(dtos)),
		byIdentifier: map[string]int{},
		byName:       map[string]int{},
	}

	for i, dto := range dtos {
		timeZone, err := dto.ToModel()

		if err != nil {
			panic(fmt.Errorf("time: loading time-zone.json: %w", err))
		}

		list.all[i] = timeZone
		list.byIdentifier[strings.ToLower(timeZone.Identifier)] = i
		list.byName[strings.ToLower(timeZone.Name)] = i
	}

	return list
}
