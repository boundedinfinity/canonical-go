# Time Zone

**Entity**

**Natural Key**: `Identifier`

A time zone, identified by its [IANA Time Zone Database](https://www.iana.org/time-zones) identifier (e.g. `America/New_York`). The other fields describe the zone's standard time, the time it keeps when daylight saving time is not in effect.

## Programmatic Information

- Identifier
    type: string
    desc: The IANA time zone identifier (e.g. `America/New_York`). This is the model's natural key.
    validation:
        - one or more parts separated by `/`
        - each part is 1 or more of the characters `A`-`Z`, `a`-`z`, `0`-`9`, `_`, `+`, `-`
        - length >= 1 and <= 100
- Name
    type: string
    desc: The English name of the zone's standard time (e.g. `Eastern Standard Time`).
    validation:
        - length >= 1 and <= 100
- Abbreviation
    type: string
    desc: The abbreviation of the zone's standard time (e.g. `EST`).
    validation:
        - 2 to 6 uppercase letters `A`-`Z`
- UtcOffset
    type: string
    desc: The zone's standard offset from UTC, as `±hh:mm` (e.g. `-05:00`).
    validation:
        - a sign (`+` or `-`), 2-digit hours from `00` to `14`, a `:`, and minutes of `00`, `15`, `30`, or `45`

Operations:
    - Location:
        inputs: (none)
        outputs:
            - the native language time zone
            - error
        desc: >
            Returns the language's time zone for `Identifier`, which applies daylight saving time and
            historical changes. Returns an error if the identifier is not in the IANA database.

## Validation

- All required fields (see above) must be present.
- `Identifier` must be a well-formed IANA identifier.
- `Name` length must be >= 1 and <= 100.
- `Abbreviation` must be 2 to 6 uppercase letters.
- `UtcOffset` must be `±hh:mm` with hours `00`-`14` and minutes `00`, `15`, `30`, or `45`.

## Lookup

A time zone can be looked up from its fields. Every comparison ignores case.

- **Get** finds the single time zone whose `Identifier` or `Name` equals the search term, checking them in that order.
- **Find** returns the result of **Get** if there is one. Otherwise it returns every time zone whose `Abbreviation` or `UtcOffset` equals the search term. It may return several time zones, or none (e.g. `+00:00` finds both `Etc/UTC` and `Europe/London`).

## Data

The list of time zones is in [time-zone.json](time-zone.json), as an array of time zone DTOs (see [DTO](../../global-notes/dto.md)) sorted by `Identifier`:

```json
{
    "identifier": "America/New_York",
    "name": "Eastern Standard Time",
    "abbreviation": "EST",
    "utc-offset": "-05:00"
}
```

Across the data set, `Identifier` and `Name` are each unique. `Abbreviation` and `UtcOffset` are not: different zones can share an offset, and abbreviations are reused around the world (`CST` is both US Central and China Standard Time; `IST` is India, Ireland, and Israel).

## Notes

- https://www.iana.org/time-zones
- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
- Abbreviations and offsets are for standard time only. Use `Location` for the offset at a given moment, including daylight saving time (e.g. `America/New_York` is `-04:00`, `EDT`, in summer).
- `Validate` checks the format of `Identifier` only; `Location` checks that it exists in the IANA database.
- The data was extracted from the legacy implementation in `z_idiomatic/measurement/time/time-zone.go`, which had no IANA identifiers. Each legacy zone is mapped to a representative IANA zone, and the legacy "Unknown" zone is dropped. The legacy "Universal Coordinated Time" and "Indian Standard Time" are renamed to their official names, "Coordinated Universal Time" and "India Standard Time", and the legacy Greenwich Mean Time is `Europe/London`.
- For go implementation `Location` uses `time.LoadLocation`. Import `time/tzdata` so the IANA database is embedded and does not depend on the host system.
