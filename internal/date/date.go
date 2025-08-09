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
	AllTime   Increment = "All time"
)

func (inc Increment) String() string {
	return string(inc)
}

// Give a date, return the first day of the increment. For example,
// Monthly.FirstDayInIncrement(2025-04-15) => 2025-04-01
// Annually.FirstDayInIncrement(2025-04-15) => 2025-01-01
func (inc Increment) FirstDayInIncrement(date time.Time) time.Time {
	switch inc {
	case Weekly:
		return firstDayOfWeek(date)
	case Monthly:
		return firstDayOfMonth(date)
	case Quarterly:
		return firstDayOfQuarter(date)
	case Annually, AllTime:
		return firstDayOfYear(date)
	default:
		panic(fmt.Sprintf("unexpected date increment: %s", inc))
	}
}

func (inc Increment) AddIncrement(date time.Time) time.Time {
	switch inc {
	case Weekly:
		return date.AddDate(0, 0, 7)
	case Monthly:
		return date.AddDate(0, 1, 0)
	case Quarterly:
		return date.AddDate(0, 3, 0)
	case Annually:
		return date.AddDate(1, 0, 0)
	default:
		panic(fmt.Sprintf("unexpected date increment: %s", inc))
	}
}

func (inc Increment) SubtractIncrement(date time.Time) time.Time {
	switch inc {
	case Weekly:
		return date.AddDate(0, 0, -7)
	case Monthly:
		return date.AddDate(0, -1, 0)
	case Quarterly:
		return date.AddDate(0, -3, 0)
	case Annually:
		return date.AddDate(-1, 0, 0)
	default:
		panic(fmt.Sprintf("unexpected date increment: %s", inc))
	}
}

func QuarterOfYear(date time.Time) int {
	return (int(date.Month())-1)/3 + 1
}

func firstDayOfWeek(date time.Time) time.Time {
	year, week := date.ISOWeek()

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

func firstDayOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func firstDayOfQuarter(date time.Time) time.Time {
	quarter := QuarterOfYear(date)
	firstMonthOfQuarter := quarter*3 - 2
	return time.Date(date.Year(), time.Month(firstMonthOfQuarter), 1, 0, 0, 0, 0, date.Location())
}

func firstDayOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
}
