# Month

A month of the year, with no year or day

**Natural Key**: Abbreviation

## Programmatic Information

Follows [ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html) month numbering, with the year starting in January.

- Month
    type: string
    desc: The English name of the month, e.g. "January".
    validation:
        - must be one of the names or abbreviations below, ignoring case

Values:

| Name      | Abbreviation | Number |
| --------- | ------------ | ------ |
| January   | Jan          | 1      |
| February  | Feb          | 2      |
| March     | Mar          | 3      |
| April     | Apr          | 4      |
| May       | May          | 5      |
| June      | Jun          | 6      |
| July      | Jul          | 7      |
| August    | Aug          | 8      |
| September | Sep          | 9      |
| October   | Oct          | 10     |
| November  | Nov          | 11     |
| December  | Dec          | 12     |

Operations:
    - Number:
        inputs: (none)
        outputs:
            - integer
        desc: Returns the ISO 8601 month number, 1 (January) through 12 (December).
    - Month:
        inputs: (none)
        outputs:
            - the native language month
        desc: Returns the month as the language's native month type.

## Notes

- Serialize as the full, capitalized name (e.g. "January"). Parsing also accepts the abbreviation and any letter case ("jan", "JANUARY"), but nothing else ("Sept", "1", " January").
- There is no zero month: a value outside 1–12 is invalid, so an unset Month fails validation.
- A Month can be built from a native date or date-time; the month is taken in that value's own time offset.
- The language's string representation (e.g. Go's `String()`) returns the Name, e.g. "January".
- For go implementation store the ISO 8601 number in an integer type. Go's `time.Month` uses the same numbering, so no conversion is needed.
