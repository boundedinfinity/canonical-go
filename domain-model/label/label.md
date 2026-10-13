# Label

**Top Level Model**

This is a label for a document. These are user defined labels.

## Programmatic Information

- Id
    type: [Id](id.md)
    desc: The unique identifier of the label.
    validation:
        - call the validation on the ID
- Name
    type: string
    desc: The name of the label.
    validation:
        - length >= 1 and <= 20
- Acronym
    type: string
    optional: true
    desc: A short acronym for the label.
    validation:
        - length >= 2 and <= 10
- Description
    type: string
    optional: true
    desc: The description of the label.
    validation:
        - length >= 1 and <= 250

## Validation

- All required fields (see above) must be present.
- `Name` length must be >= 1 and <= 20.
- `Acronym`, if present, length must be >= 2 and <= 10.
- `Desc`, if present, length must be >= 1 and <= 250.
