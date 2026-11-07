# Quarter

A calendar quarter of the year, with no year

**Natural Key**: Name

## Programmatic Information

Divides the calendar year into four quarters of three [Months](month.md) each, starting in January.

- Quarter
    type: string
    desc: The name of the quarter, e.g. "Q1".
    validation:
        - must be one of the names below, ignoring case

Values:

| Name | Number | Months                         |
| ---- | ------ | ------------------------------ |
| Q1   | 1      | January, February, March       |
| Q2   | 2      | April, May, June               |
| Q3   | 3      | July, August, September        |
| Q4   | 4      | October, November, December    |

Operations:
    - Number:
        inputs: (none)
        outputs:
            - integer
        desc: Returns the quarter number, 1 (Q1) through 4 (Q4).
    - Months:
        inputs: (none)
        outputs:
            - list of [Month](month.md)
        desc: Returns the quarter's three months, in order.

## Notes

- Serialize as the capitalized name (e.g. "Q1"). Parsing also accepts any letter case ("q1"), but nothing else ("1", "Q5", "first", " Q1").
- There is no zero quarter: a value outside 1–4 is invalid, so an unset Quarter fails validation.
- A Quarter can be built from a native date or date-time (the month is taken in that value's own time offset), or from a [Month](month.md).
- These are calendar quarters. Fiscal quarters, which start in another month, are not modeled.
- The language's string representation (e.g. Go's `String()`) returns the Name, e.g. "Q1".
- For go implementation store the number in an integer type. Go has no native quarter type.
