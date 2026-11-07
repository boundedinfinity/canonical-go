# Day Of Week Range

A set of days of the week, e.g. the working days of a business: Monday through Friday, and Saturday.

## Programmatic Information

Stores a list of [Day Of Week](day-of-week.md) entities.

Properties: 
    - Days
        type: list of [Day Of Week](day-of-week.md)
        desc: The days in the set, in ISO 8601 order (Monday first). May be empty.
        validation:
            - call the validation on each Day Of Week
            - must not contain the same day more than once
Operations:
    - AddDay:
        inputs: 
            - [Day Of Week](day-of-week.md)
        outputs:
            - error
        desc: >
            Adds the day to `Days`, keeping ISO 8601 order. Adding a day that is already present does nothing.
            Returns an error if the day is invalid.
    - AddRange:
        inputs: 
            - [Day Of Week](day-of-week.md) (start)
            - [Day Of Week](day-of-week.md) (end)
        outputs:
            - error
        desc: >
            Adds every day from start to end, inclusive, as with `AddDay`. If start comes after end, the
            range wraps past Sunday: Friday to Monday adds Friday, Saturday, Sunday, and Monday.
            Returns an error, and adds nothing, if either day is invalid.

## Validation

- `Days` must not contain the same day more than once.

## Notes

- Example: working days of Monday through Friday, and Saturday, are `AddRange(Monday, Friday)` then `AddDay(Saturday)`, giving all days but Sunday.
- The days are a set: adding the same day twice, or adding overlapping ranges, stores each day once.
