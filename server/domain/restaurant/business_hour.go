package restaurant

type DayOfWeekEnum string

const (
	DayOfWeekSunday    DayOfWeekEnum = "DAY_OF_WEEK_SUNDAY"
	DayOfWeekMonday    DayOfWeekEnum = "DAY_OF_WEEK_MONDAY"
	DayOfWeekTuesday   DayOfWeekEnum = "DAY_OF_WEEK_TUESDAY"
	DayOfWeekWednesday DayOfWeekEnum = "DAY_OF_WEEK_WEDNESDAY"
	DayOfWeekThursday  DayOfWeekEnum = "DAY_OF_WEEK_THURSDAY"
	DayOfWeekFriday    DayOfWeekEnum = "DAY_OF_WEEK_FRIDAY"
	DayOfWeekSaturday  DayOfWeekEnum = "DAY_OF_WEEK_SATURDAY"
)

type BusinessHour struct {
	DayOfWeekEnum DayOfWeekEnum `json:"dayOfWeekEnum"`
	OpenTime      string        `json:"openTime"`
	ClosingTime   string        `json:"closingTime"`
	IsClosedDay   bool          `json:"isClosedDay"`
}
