# Subdivision

**Top Level Model**

**Natural Key**: `Code`

A principal subdivision of a country (e.g. a state, province, or region), as
defined by [ISO 3166-2](https://www.iso.org/iso-3166-country-codes.html). Each
subdivision belongs to one [Country](country-code.md).

## Programmatic Information

- Code
    type: string
    desc: >
        The ISO 3166-2 code (e.g. `US-AL`). It is the country's alpha-2 code, a
        hyphen, and 1 to 3 letters or digits. This is the model's natural key.
    validation:
        - 2 uppercase letters `A`-`Z`, then `-`, then 1 to 3 characters that are each an uppercase letter `A`-`Z` or a digit `0`-`9`
- Country
    type: [Country Code](country-code.md)
    desc: The country this subdivision belongs to.
    validation:
        - call the validation on the [Country Code](country-code.md)
- Name
    type: string
    desc: The name of the subdivision (e.g. `Alabama`).
    validation:
        - length >= 1 and <= 100
- Category
    type: string
    optional: true
    desc: The kind of subdivision (e.g. `state`, `province`, `region`).
    validation:
        - length >= 1 and <= 50
- Parent
    type: [Subdivision](subdivision.md)
    optional: true
    desc: >
        The subdivision this one is part of, for countries whose subdivisions
        have more than one level.
    validation:
        - call the validation on the [Subdivision](subdivision.md)
- AlternateCode
    type: string
    optional: true
    desc: Another code commonly used for the subdivision.
    validation:
        - length >= 1 and <= 10
- PreviousCode
    type: string
    optional: true
    desc: The subdivision's code before ISO last changed it.
    validation:
        - same format as `Code`
- Reference
    type: string
    optional: true
    desc: A link to more information about the subdivision.
    validation:
        - must be an absolute URL with an `http` or `https` scheme

## Validation

- All required fields (see above) must be present.
- `Code` must be 2 uppercase letters, `-`, and 1 to 3 uppercase letters or digits.
- The first 2 characters of `Code` must equal the `Alpha2` of `Country`.
- `Name` length must be >= 1 and <= 100.
- `Category`, if present, length must be >= 1 and <= 50.
- `Parent`, if present, must belong to the same `Country`, and must not be the
  subdivision itself (its `Code` must differ from this `Code`).
- Following `Parent` from subdivision to subdivision must never reach a `Code`
  that was already seen (no cycles).
- `AlternateCode`, if present, length must be >= 1 and <= 10.
- `PreviousCode`, if present, must have the same format as `Code`.
- `Reference`, if present, must be an absolute `http` or `https` URL.

## Lookup

Every comparison ignores case.

- **Get** finds the single subdivision whose `Code` equals the search term. The
  term also matches without the hyphen, so `USAL` finds `US-AL`.
- **Find by name** returns every subdivision whose `Name` equals the search
  term. Names are not unique (e.g. `Western` is a subdivision of 8 countries),
  so it may return several subdivisions, or none.
- **By country** returns every subdivision whose `Country` is the given country.

## Data

The full list of subdivisions is in [subdivision.json](subdivision.json), as an
array of subdivision DTOs (see [DTO](../../global-notes/dto.md)) sorted by
`Code`. `Country` holds the country's `Alpha2`, and `Parent`, when present, holds
the parent's `Code`:

```json
{
    "code": "US-AL",
    "country": "US",
    "name": "Alabama",
    "reference": "https://en.wikipedia.org/wiki/Alabama"
}
```

Across the data set, `Code` is unique, and every `Country` exists in
[country-code.json](country-code.json). `Category`, `Parent`, `AlternateCode`,
and `PreviousCode` are not yet filled in for any subdivision.

The data is incomplete: 56 countries have no subdivisions in it. Many of them
have none in ISO 3166-2 either (e.g. `GI` Gibraltar, `AQ` Antarctica), but some
do and are missing, including `GB`, `IN`, `SE`, `EE`, `LT`, `MZ`, and `HT`.

## Notes

- https://www.iso.org/iso-3166-country-codes.html
- https://en.wikipedia.org/wiki/ISO_3166-2
- The data was extracted from the legacy implementation in
  `idiomatic/specification/iso/iso3166`. The legacy code nested each country's
  subdivisions inside the country. Here a subdivision refers to its country
  instead, which avoids a circular reference between the two models.
- The legacy name lookup kept only one subdivision per name, so it could not
  find most of the subdivisions that share a name. **Find by name** returns all
  of them.
