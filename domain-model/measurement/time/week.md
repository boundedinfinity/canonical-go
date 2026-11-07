# Week

All the days of a week, in order, starting on a configurable first day. The default first day is Monday.

## Programmatic Information

- Start
    type: [Day Of Week](day-of-week.md)
    optional: true
    desc: The first day of the week. When not present, the week starts on Monday.
    validation:
        - if present, call the validation on the Day Of Week

Operations:
    - First:
        inputs: (none)
        outputs:
            - [Day Of Week](day-of-week.md)
        desc: Returns the first day of the week, `Start` or Monday when not present.
    - Last:
        inputs: (none)
        outputs:
            - [Day Of Week](day-of-week.md)
        desc: Returns the last day of the week, the day before `First`.
    - Days:
        inputs: (none)
        outputs:
            - list of [Day Of Week](day-of-week.md)
        desc: Returns all 7 days in order, starting with `First`.
    - Position:
        inputs:
            - [Day Of Week](day-of-week.md)
        outputs:
            - integer
        desc: Returns the day's place in the week, 1 (`First`) through 7 (`Last`), or 0 if the day is invalid.
    - Compare:
        inputs:
            - [Day Of Week](day-of-week.md)
            - [Day Of Week](day-of-week.md)
        outputs:
            - integer
        desc: Returns -1, 0, or 1 as the first day comes before, is the same as, or comes after the second, by `Position`.

## Validation

- `Start`, if present, must be a valid [Day Of Week](day-of-week.md).

## Notes

- The days are derived from `Start` rather than stored, so the order can never disagree with the first day.
- The first day changes only the order of days (`Days`, `Position`, `Compare`). A Day Of Week's `Number` is always the ISO 8601 number, whatever the week's first day.
