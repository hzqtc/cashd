package model

import (
	"cashd/internal/data"
	"cashd/internal/date"
	"cashd/internal/ui"
	"fmt"
	"sort"
	"time"
)

// Return if the Transaction matches the aggregation requirements
type matchFunc func(*data.Transaction) bool

func aggregateByAccount(transactions []*data.Transaction, aggLevel date.Increment, accountName string) []*ui.TsChartEntry {
	return aggregate(
		transactions,
		aggLevel,
		func(t *data.Transaction) bool {
			return accountName == "Total" || t.Account == accountName
		})
}

func aggregateByCategory(transactions []*data.Transaction, aggLevel date.Increment, categoryName string) []*ui.TsChartEntry {
	return aggregate(
		transactions,
		aggLevel,
		func(t *data.Transaction) bool {
			return t.Category == categoryName
		})
}

func aggregate(transactions []*data.Transaction, aggLevel date.Increment, matches matchFunc) []*ui.TsChartEntry {
	// Store aggregated results in a map for easier access by date
	// It's critical to use pointers to update entries
	entryMap := make(map[time.Time]*ui.TsChartEntry)
	for _, t := range transactions {
		if !matches(t) {
			continue
		}
		// This is the key: aggregate results by date
		date := getAggLevelDate(t.Date, aggLevel)
		entry, exist := entryMap[date]
		if !exist {
			entry = &ui.TsChartEntry{Date: date}
			entryMap[date] = entry
		}
		if t.Type == data.Income {
			entry.Income += t.Amount
		} else {
			entry.Expense += t.Amount
		}
	}
	// Convert aggregated results to an array and sort by date
	entries := []*ui.TsChartEntry{}
	for _, entry := range entryMap {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})
	return entries
}

// Convert a date to the first date in its aggregate level
// For example, 2023-04-15 in month aggregation -> 2023-04-01
func getAggLevelDate(d time.Time, aggLevel date.Increment) time.Time {
	switch aggLevel {
	case date.Weekly:
		return date.FirstDayOfWeek(d)
	case date.Monthly:
		return date.FirstDayOfMonth(d)
	case date.Quarterly:
		return date.FirstDayOfQuarter(d)
	case date.Annually:
		return date.FirstDayOfYear(d)
	default:
		panic(fmt.Sprintf("Unexpected date aggregate level: %s", aggLevel))
	}
}

