# Country

**Top Level Model**

This is a country definition.

## Programmatic Information

- Id
    type: [Id](id.md)
    desc: The unique identifier of the country.
    validation:
        - call the validation on the ID
- Name
    type: string
    desc: The name of the country.
    validation:
        - length >= 1 and <= 100
- Iso3166
    type: [Iso]

## Validation

- All required fields (see above) must be present.
- `Name` length must be >= 1 and <= 20.
- `Acronyms`, if present, must contain at least 1 element, and each acronym length must be >= 2 and <= 10.
- `Description`, if present, length must be >= 1 and <= 250.