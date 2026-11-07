# Month Range

A set of months of the year, e.g. the season of a business: April through September, and December.

## Programmatic Information

Stores a list of [Month](month.md) entities.

Properties: 
    - Months
        type: list of [Month](month.md)
        desc: The months in the set, in ISO 8601 order (January first). May be empty.
        validation:
            - call the validation on each Month
            - must not contain the same month more than once
Operations:
    - AddMonth:
        inputs: 
            - [Month](month.md)
        outputs:
            - error
        desc: >
            Adds the month to `Months`, keeping ISO 8601 order. Adding a month that is already present does nothing.
            Returns an error if the month is invalid.
    - AddRange:
        inputs: 
            - [Month](month.md) (start)
            - [Month](month.md) (end)
        outputs:
            - error
        desc: >
            Adds every month from start to end, inclusive, as with `AddMonth`. If start comes after end, the
            range wraps past December: November to February adds November, December, January, and February.
            Returns an error, and adds nothing, if either month is invalid.
    - Contains:
        inputs: 
            - [Month](month.md)
        outputs:
            - boolean
            - error
        desc: >
            Returns true if the month is in `Months`, false otherwise.
            Returns an error if the month is invalid.

## Validation

- `Months` must not contain the same month more than once.

## Notes

- Example: a season of April through September, and December, is `AddRange(April, September)` then `AddMonth(December)`.
- The months are a set: adding the same month twice, or adding overlapping ranges, stores each month once.
