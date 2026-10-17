# Country Code

**Top Level Model**

**Natural Key**: `Alpha2`

A country, dependent territory, or special area of geographical interest, as
defined by [ISO 3166-1](https://www.iso.org/iso-3166-country-codes.html). Each
country has three codes: a two-letter code, a three-letter code, and a
three-digit numeric code.

Its subdivisions (states, provinces, etc.) are described by
[Subdivision](subdivision.md).

## Programmatic Information

- Alpha2
    type: string
    desc: The ISO 3166-1 alpha-2 code (e.g. `US`). This is the model's natural key.
    validation:
        - exactly 2 characters
        - each character is an uppercase letter `A`-`Z`
- Alpha3
    type: string
    desc: The ISO 3166-1 alpha-3 code (e.g. `USA`).
    validation:
        - exactly 3 characters
        - each character is an uppercase letter `A`-`Z`
- Numeric
    type: string
    desc: >
        The ISO 3166-1 numeric code (e.g. `840`). It is a string, not a number,
        because ISO writes it with leading zeros (e.g. `004` for Afghanistan).
    validation:
        - exactly 3 characters
        - each character is a digit `0`-`9`
- Name
    type: string
    desc: The English short name of the country (e.g. `United States of America`).
    validation:
        - length >= 1 and <= 100
- Independent
    type: boolean
    desc: >
        Whether the country is an independent sovereign state (`true`), or a
        dependent territory or special area (`false`).
- Reference
    type: string
    optional: true
    desc: A link to more information about the country.
    validation:
        - must be an absolute URL with an `http` or `https` scheme

## Validation

- All required fields (see above) must be present.
- `Alpha2` must be exactly 2 uppercase letters.
- `Alpha3` must be exactly 3 uppercase letters.
- `Numeric` must be exactly 3 digits.
- `Name` length must be >= 1 and <= 100.
- `Reference`, if present, must be an absolute `http` or `https` URL.

## Lookup

A country can be looked up from its codes or name. Every comparison ignores case.

- **Get** finds the single country whose `Name`, `Alpha2`, `Alpha3`, or `Numeric`
  equals the search term, checking them in that order. A numeric term also
  matches without its leading zeros, so `4` and `004` both find Afghanistan.
- **Find** returns the result of **Get** if there is one. Otherwise it returns
  every country whose `Name`, `Alpha2`, `Alpha3`, or `Numeric` contains the
  search term. It may return several countries, or none.

## Data

The full list of countries is in [country-code.json](country-code.json), as an
array of country DTOs (see [DTO](../../global-notes/dto.md)) sorted by `Alpha2`:

```json
{
    "alpha-2": "US",
    "alpha-3": "USA",
    "numeric": "840",
    "name": "United States of America",
    "independent": true,
    "reference": "https://en.wikipedia.org/wiki/United_States"
}
```

Across the data set, `Alpha2`, `Alpha3`, `Numeric`, and `Name` are each unique.

## Notes

- https://www.iso.org/iso-3166-country-codes.html
- https://en.wikipedia.org/wiki/ISO_3166-1
- The data was extracted from the legacy implementation in
  `idiomatic/specification/iso/iso3166`, which was generated from the Wikipedia
  ISO 3166-1 tables. The legacy `Independent` field held the strings `Yes` and
  `No`.
