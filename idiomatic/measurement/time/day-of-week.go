package time

type DayOfWeek string

type daysOfWeek struct {
	Unknown   DayOfWeek
	Monday    DayOfWeek
	Tuesday   DayOfWeek
	Wednesday DayOfWeek
	Thursday  DayOfWeek
	Friday    DayOfWeek
	Saturday  DayOfWeek
	Sunday    DayOfWeek
}

var DaysOfWeek = daysOfWeek{
	Unknown:   "measurement.time.day-of-week.unknown",
	Monday:    "measurement.time.day-of-week.monday",
	Tuesday:   "measurement.time.day-of-week.tuesday",
	Wednesday: "measurement.time.day-of-week.wednesday",
	Thursday:  "measurement.time.day-of-week.thursday",
	Friday:    "measurement.time.day-of-week.friday",
	Saturday:  "measurement.time.day-of-week.saturday",
	Sunday:    "measurement.time.day-of-week.sunday",
}

type DayOfWeekRange struct {
	Min DayOfWeek `json:"min"`
	Max DayOfWeek `json:"max"`
}
