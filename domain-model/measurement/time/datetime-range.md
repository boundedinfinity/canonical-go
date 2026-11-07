# DateTime Range

A datetime entity stores a range of DateTime entities

## Programmatic Information

Stores a range of [DateTime](datetime.md) entities.

Properties: 
    - Start
        type: [DateTime](datetime.md)
        desc: The start DateTime
        validation:
            - call the validation on the DateTime
    - End
        type: [DateTime](datetime.md)
        desc: The end DateTime
        validation:
            - call the validation on the DateTime
Operations:
    - Between:
        inputs: 
            - [DateTime](datetime.md)
        outputs:
            - boolean
        desc: Returns true if DateTime is between the start and end properties, inclusive of both, false otherwise.
    - Before:
        inputs: 
            - [DateTime](datetime.md)
        outputs:
            - boolean
        desc: Returns true if DateTime is strictly before the start property, false otherwise.
    - After:
        inputs: 
            - [DateTime](datetime.md)
        outputs:
            - boolean
        desc: Returns true if DateTime is strictly after the end property, false otherwise.
    - Overlaps:
        inputs:
            - DateTime Range
        outputs:
            - boolean
        desc: >
            Returns true if the input shares at least one instant with this range, false otherwise.
            Bounds are inclusive, so ranges that only touch at an endpoint overlap.

## Validation

- All required fields (see above) must be present.
- `Start` must be <= `End`. A range where `Start` equals `End` is a single instant.

## Notes

- For any DateTime, exactly one of `Before`, `Between`, and `After` is true.
- Comparisons are by instant, so DateTimes with different time offsets compare correctly.
