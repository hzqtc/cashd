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
	case Annually:
		return firstDayOfYear(date)
	default:
		panic(fmt.Sprintf("unexpected date increment: %s", inc))
	}
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
	quarter := quarterOfYear(date)
	firstMonthOfQuarter := quarter*3 - 2
	return time.Date(date.Year(), time.Month(firstMonthOfQuarter), 1, 0, 0, 0, 0, date.Location())
}

func quarterOfYear(date time.Time) int {
	return (int(date.Month())-1)/3 + 1
}

func firstDayOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
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

func (inc Increment) FormatDate(date time.Time) string {
	switch inc {
	case Weekly:
		year, week := date.ISOWeek()
		return fmt.Sprintf("%d week %02d", year, week)
	case Monthly:
		return fmt.Sprintf("%s", date.Format("2006 January"))
	case Quarterly:
		return fmt.Sprintf("%d Q%d", date.Year(), quarterOfYear(date))
	case Annually:
		return fmt.Sprintf("%d", date.Year())
	}
	panic(fmt.Sprintf("Unexpected date increment: %s", inc))
}

// Shorter format with less space and 2 digits year number
func (inc Increment) FormatDateShorter(date time.Time) string {
	year := date.Format("06")
	switch inc {
	case Weekly:
		_, week := date.ISOWeek()
		return fmt.Sprintf("%s'W%02d", year, week)
	case Monthly:
		return fmt.Sprintf("%s'%s", year, date.Format("Jan"))
	case Quarterly:
		return fmt.Sprintf("%s'Q%d", year, quarterOfYear(date))
	case Annually:
		return fmt.Sprintf("%d", date.Year())
	}
	panic(fmt.Sprintf("Unexpected date increment: %s", inc))
}
