package date

import (
	"fmt"
	"time"
)

type Increment string

const (
	Weekly    Increment = "Weekly"
	Monthly   Increment = "Monthly"
	Quarterly Increment = "Quarterly"
	Annually  Increment = "Yearly"
)

func (inc Increment) Days() int {
	switch inc {
	case Weekly:
		return 7
	case Monthly:
		return 30
	case Quarterly:
		return 91
	case Annually:
		return 365
	}
	panic(fmt.Sprintf("Unexpected date increment: %s", inc))
}

func FirstDayOfWeek(refDate time.Time) time.Time {
	year, week := refDate.ISOWeek()

	// ISO 8601: Week 1 is the week with the first Thursday of the year
	// Start with Jan 4th (guaranteed to be in week 1)
	jan4 := time.Date(year, time.January, 4, 0, 0, 0, 0, time.UTC)

	// Get the ISO weekday of Jan 4th (Monday=1, Sunday=7)
	isoWeekday := int(jan4.Weekday())
	if isoWeekday == 0 {
		isoWeekday = 7
	}

	// Go back to Monday of that week
	monday := jan4.AddDate(0, 0, -isoWeekday+1)

	// Add (week - 1) weeks
	return monday.AddDate(0, 0, (week-1)*7)
}

func FirstDayOfMonth(refDate time.Time) time.Time {
	return time.Date(refDate.Year(), refDate.Month(), 1, 0, 0, 0, 0, refDate.Location())
}

func FirstDayOfQuarter(refDate time.Time) time.Time {
	quarter := QuarterOfYear(refDate)
	firstMonthOfQuarter := quarter*3 - 2
	return time.Date(refDate.Year(), time.Month(firstMonthOfQuarter), 1, 0, 0, 0, 0, refDate.Location())
}

func QuarterOfYear(date time.Time) int {
	return (int(date.Month())-1)/3 + 1
}

func FirstDayOfYear(refDate time.Time) time.Time {
	return time.Date(refDate.Year(), 1, 1, 0, 0, 0, 0, refDate.Location())
}

func FormatDateToIncrement(date time.Time, inc Increment) string {
	switch inc {
	case Weekly:
		year, week := date.ISOWeek()
		return fmt.Sprintf("%d week %02d", year, week)
	case Monthly:
		return fmt.Sprintf("%s", date.Format("January 2006"))
	case Quarterly:
		return fmt.Sprintf("%d Q%d", date.Year(), QuarterOfYear(date))
	case Annually:
		return fmt.Sprintf("%d", date.Year())
	}
	panic(fmt.Sprintf("Unexpected date increment: %s", inc))
}
