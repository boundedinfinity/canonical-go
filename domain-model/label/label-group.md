# Label Group

**Top Level Model**

A user defined group of [Labels](label.md).

## Programmatic Information

- Id
    type: [Id](../id/id.md)
    desc: The unique identifier of the label group.
    validation:
        - call the validation on the ID
- Name
    type: string
    desc: The name of the label group.
    validation:
        - length >= 1 and <= 20
- Acronym
    type: string
    optional: true
    desc: A short acronym for the label group.
    validation:
        - length >= 2 and <= 10
- Description
    type: string
    optional: true
    desc: The description of the label group.
    validation:
        - length >= 1 and <= 250
- Labels
    type: List [Label](label.md)
    optional: true
    desc: The labels for this group.
    validation:
        - if present, at least 1 element
        - call the validation for each [Label](label.md)

## Validation

- All required fields (see above) must be present.
- `Id` must be a valid [Id](../id/id.md).
- `Name` length must be >= 1 and <= 20.
- `Acronym`, if present, length must be >= 2 and <= 10.
- `Description`, if present, length must be >= 1 and <= 250.
- `Labels`, if present, must contain at least 1 element, and each [Label](label.md) must be valid.
