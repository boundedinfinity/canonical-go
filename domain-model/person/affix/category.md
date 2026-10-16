
# Affix Category

The "base" struct for prefix and suffix categories.

## Programmatic Information

- Id
    type: [Id](../../id/id.md)
    desc: The unique identifier of the affix category.
    validation:
        - call the validation on the ID
- Name
    type: string
    desc: The name of the affix category.
    validation:
        - length >= 1 and <= 30
- Description
    type: string
    desc: The description of the affix category.
    optional: true
    validation:
        - length >= 2 and <= 500

## Validation

- All required fields (see above) must be present.
- `ID` must be a valid [Id](../../id/id.md).
- `Name` length must be >= 1 and <= 30.
- `Description`, if present, length must be >= 2 and <= 500.
