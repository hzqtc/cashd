package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type increment string

const (
	weekly    increment = "Week"
	monthly   increment = "Month"
	quarterly increment = "Quarter"
	annually  increment = "Year"
)

type DatePickerModel struct {
	width int

	startDate time.Time // Inclusive
	endDate   time.Time // Exclusive
	inc       increment

	next      key.Binding
	prev      key.Binding
	byWeek    key.Binding
	byMonth   key.Binding
	byQuarter key.Binding
	byYear    key.Binding
}

type DateRangeChangedMsg struct {
	Start time.Time // Inclusive
	End   time.Time // Exclusive
}

func NewDatePickerModel() DatePickerModel {
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return DatePickerModel{
		// First day of the month
		startDate: currentMonth,
		endDate:   currentMonth.AddDate(0, 1, 0),
		inc:       monthly,
		next:      key.NewBinding(key.WithKeys("l", "right")),
		prev:      key.NewBinding(key.WithKeys("h", "left")),
		byWeek:    key.NewBinding(key.WithKeys("w")),
		byMonth:   key.NewBinding(key.WithKeys("m")),
		byQuarter: key.NewBinding(key.WithKeys("t")),
		byYear:    key.NewBinding(key.WithKeys("y")),
	}
}

func (m *DatePickerModel) SetWidth(w int) {
	m.width = w
}

func (m DatePickerModel) Update(msg tea.Msg) (DatePickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.prev):
			m.prevDateRange()
		case key.Matches(msg, m.next):
			m.nextDateRange()
		case key.Matches(msg, m.byWeek):
			m.inc = weekly
			m.updateIncrement()
		case key.Matches(msg, m.byMonth):
			m.inc = monthly
			m.updateIncrement()
		case key.Matches(msg, m.byQuarter):
			m.inc = quarterly
			m.updateIncrement()
		case key.Matches(msg, m.byYear):
			m.inc = annually
			m.updateIncrement()
		default:
			// Important to return nil cmd if nothing is changed
			return m, nil
		}
	default:
		// Important to return nil cmd if nothing is changed
		return m, nil
	}
	return m, func() tea.Msg {
		return DateRangeChangedMsg{
			Start: m.startDate,
			End:   m.endDate,
		}
	}
}

func (m *DatePickerModel) SelectedDateRange() (time.Time, time.Time) {
	return m.startDate, m.endDate
}

func (m *DatePickerModel) nextDateRange() {
	switch m.inc {
	case weekly:
		m.startDate = m.startDate.AddDate(0, 0, 7)
		m.endDate = m.endDate.AddDate(0, 0, 7)
	case monthly:
		m.startDate = m.startDate.AddDate(0, 1, 0)
		m.endDate = m.endDate.AddDate(0, 1, 0)
	case quarterly:
		m.startDate = m.startDate.AddDate(0, 3, 0)
		m.endDate = m.endDate.AddDate(0, 3, 0)
	case annually:
		m.startDate = m.startDate.AddDate(1, 0, 0)
		m.endDate = m.endDate.AddDate(1, 0, 0)
	}
}

func (m *DatePickerModel) prevDateRange() {
	switch m.inc {
	case weekly:
		m.startDate = m.startDate.AddDate(0, 0, -7)
		m.endDate = m.endDate.AddDate(0, 0, -7)
	case monthly:
		m.startDate = m.startDate.AddDate(0, -1, 0)
		m.endDate = m.endDate.AddDate(0, -1, 0)
	case quarterly:
		m.startDate = m.startDate.AddDate(0, -3, 0)
		m.endDate = m.endDate.AddDate(0, -3, 0)
	case annually:
		m.startDate = m.startDate.AddDate(-1, 0, 0)
		m.endDate = m.endDate.AddDate(-1, 0, 0)
	}
}

func (m DatePickerModel) View() string {
	var leftStr strings.Builder
	var rightStr strings.Builder

	switch m.inc {
	case weekly:
		year, week := m.startDate.ISOWeek()
		leftStr.WriteString(fmt.Sprintf("< %d week %02d >", year, week))
	case monthly:
		leftStr.WriteString(fmt.Sprintf("< %s >", m.startDate.Format("January 2006")))
	case quarterly:
		leftStr.WriteString(fmt.Sprintf("< %d Q%d >", m.startDate.Year(), quarterOfYear(m.startDate)))
	case annually:
		leftStr.WriteString(fmt.Sprintf("< %d >", m.startDate.Year()))
	}

	// Key bindings
	rightStr.WriteString("Navigate: ")
	rightStr.WriteString(keyStyle.Render("h/←") + " prev")
	rightStr.WriteString(" ")
	rightStr.WriteString(keyStyle.Render("l/→") + " next")
	rightStr.WriteString(" ")
	rightStr.WriteString("Switch to: ")
	rightStr.WriteString(keyStyle.Render(m.byWeek.Keys()[0]) + "eekly")
	rightStr.WriteString(" ")
	rightStr.WriteString(keyStyle.Render(m.byMonth.Keys()[0]) + "onthly")
	rightStr.WriteString(" ")
	rightStr.WriteString("quar" + keyStyle.Render(m.byQuarter.Keys()[0]) + "erly")
	rightStr.WriteString(" ")
	rightStr.WriteString(keyStyle.Render(m.byYear.Keys()[0]) + "early")

	style := lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(fmt.Sprintf("View by: %s", m.inc), m.width)).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(m.width)
	spaces := m.width - hPadding*2 - lipgloss.Width(leftStr.String()) - lipgloss.Width(rightStr.String())

	return style.
		Render(leftStr.String() + strings.Repeat(" ", max(0, spaces)) + rightStr.String())
}

func (m *DatePickerModel) updateIncrement() {
	// Snap start and end dates to increment
	switch m.inc {
	case weekly:
		m.startDate = firstDayOfWeek(m.startDate)
		m.endDate = m.startDate.AddDate(0, 0, 7)
	case monthly:
		m.startDate = firstDayOfMonth(m.startDate)
		m.endDate = m.startDate.AddDate(0, 1, 0)
	case quarterly:
		m.startDate = firstDayOfQuarter(m.startDate)
		m.endDate = m.startDate.AddDate(0, 3, 0)
	case annually:
		m.startDate = firstDayOfYear(m.startDate)
		m.endDate = m.startDate.AddDate(1, 0, 0)
	}
}

func firstDayOfWeek(refDate time.Time) time.Time {
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

func firstDayOfMonth(refDate time.Time) time.Time {
	return time.Date(refDate.Year(), refDate.Month(), 1, 0, 0, 0, 0, refDate.Location())
}

func firstDayOfQuarter(refDate time.Time) time.Time {
	quarter := quarterOfYear(refDate)
	firstMonthOfQuarter := quarter*3 - 2
	return time.Date(refDate.Year(), time.Month(firstMonthOfQuarter), 1, 0, 0, 0, 0, refDate.Location())
}

func quarterOfYear(date time.Time) int {
	return (int(date.Month())-1)/3 + 1
}

func firstDayOfYear(refDate time.Time) time.Time {
	return time.Date(refDate.Year(), 1, 1, 0, 0, 0, 0, refDate.Location())
}
