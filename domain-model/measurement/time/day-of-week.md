# Day Of Week

A day of the week, with no date

**Natural Key**: Abbreviation

## Programmatic Information

Follows [ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html) day numbering, with the week starting on Monday.

- Day Of Week
    type: string
    desc: The English name of the day, e.g. "Monday".
    validation:
        - must be one of the names or abbreviations below, ignoring case

Values:

| Name      | Abbreviation | Number |
| --------- | ------------ | ------ |
| Monday    | Mon          | 1      |
| Tuesday   | Tue          | 2      |
| Wednesday | Wed          | 3      |
| Thursday  | Thu          | 4      |
| Friday    | Fri          | 5      |
| Saturday  | Sat          | 6      |
| Sunday    | Sun          | 7      |

Operations:
    - Number:
        inputs: (none)
        outputs:
            - integer
        desc: Returns the ISO 8601 day number, 1 (Monday) through 7 (Sunday).
    - Weekday:
        inputs: (none)
        outputs:
            - the native language day of week
        desc: Returns the day as the language's native day of week type.

## Notes

- Serialize as the full, capitalized name (e.g. "Monday"). Parsing also accepts the abbreviation and any letter case ("mon", "MONDAY"), but nothing else ("Tues", "1", " Monday").
- There is no zero day: a value outside 1–7 is invalid, so an unset Day Of Week fails validation.
- A Day Of Week can be built from a native date or date-time; the day is taken in that value's own time offset.
- The language's string representation (e.g. Go's `String()`) returns the Name, e.g. "Monday".
- For go implementation store the ISO 8601 number in an integer type. Go's `time.Weekday` numbers Sunday as 0, so convert Sunday between 0 and 7.