# Operating Hours

The hours a business or service is open on a set of days of the week, e.g. Monday through Friday, 9 AM to 5 PM Eastern Time.

## Programmatic Information

- Days
    type: [Day Of Week Range](day-of-week-range.md)
    desc: The days these hours apply to.
    validation:
        - call the validation on the Day Of Week Range
        - at least 1 day
- Hours
    type: [Hour Of Day Range](hour-of-day-range.md)
    desc: >
        The hours it is open, on the local clock of `TimeZone`. Both ends are open hours:
        Start 9 and End 16 is open from 9:00 until 17:00.
    validation:
        - call the validation on the Hour Of Day Range
- TimeZone
    type: [Time Zone](time-zone.md)
    desc: The time zone the hours are kept in.
    validation:
        - call the validation on the Time Zone
- Clock
    type: [Clock Format](clock-format.md)
    optional: true
    desc: How to display the hours. When not present, the display format is left to the caller.
    validation:
        - if present, call the validation on the Clock Format

## Validation

- All required fields (see above) must be present.
- `Days` must contain at least 1 day.

## Notes

- `Hours` are local clock hours, not fixed offsets: 9 AM in `America/New_York` is 9 AM in both winter (EST) and summer (EDT). Use the Time Zone's `Location` to find the actual instant.
- Because `End` is an open hour, a business that closes at 5 PM has End 16, one that is open until midnight has End 23, and one that is open 24 hours has Start 0 and End 23.
- Hours do not wrap past midnight. Hours that cross midnight are two Operating Hours: open Friday 6 PM to 2 AM is Friday 18 to 23, and Saturday 0 to 1.
- Hours are whole hours; opening or closing on the half hour is not modeled.
- A business with different hours on different days holds a list of Operating Hours, e.g. Monday through Friday 9 to 16, and Saturday 10 to 13.
- The data was based on the legacy implementation in `z_idiomatic/measurement/time/operating-hours.go`, which held a single Day Of Week, a [Time Range](time-range.md) (whose times carry fixed offsets), a Time Zone, and a Clock Format.
