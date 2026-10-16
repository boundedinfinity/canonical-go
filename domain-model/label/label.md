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
- Acronyms
    type: list of strings
    optional: true
    desc: >
        A short acronym for the label. 
        
        First item in the list is the "Primary" acronym.
    validation:
        - if present, at least 1 element
        - each string length >= 2 and <= 10
- Description
    type: string
    optional: true
    desc: The description of the label.
    validation:
        - length >= 1 and <= 250

## Validation

- All required fields (see above) must be present.
- `Name` length must be >= 1 and <= 20.
- `Acronyms`, if present, must contain at least 1 element, and each acronym length must be >= 2 and <= 10.
- `Description`, if present, length must be >= 1 and <= 250.
