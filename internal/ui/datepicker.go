package ui

import (
	"cashd/internal/date"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DatePickerModel struct {
	width int

	startDate time.Time // Inclusive
	endDate   time.Time // Exclusive
	inc       date.Increment

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
		inc:       date.Monthly,
		next:      key.NewBinding(key.WithKeys("l", "right")),
		prev:      key.NewBinding(key.WithKeys("h", "left")),
		byWeek:    key.NewBinding(key.WithKeys("w")),
		byMonth:   key.NewBinding(key.WithKeys("m")),
		byQuarter: key.NewBinding(key.WithKeys("q")),
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
			if m.inc == date.Weekly {
				return m, nil
			}
			m.inc = date.Weekly
			m.updateIncrement()
		case key.Matches(msg, m.byMonth):
			if m.inc == date.Monthly {
				return m, nil
			}
			m.inc = date.Monthly
			m.updateIncrement()
		case key.Matches(msg, m.byQuarter):
			if m.inc == date.Quarterly {
				return m, nil
			}
			m.inc = date.Quarterly
			m.updateIncrement()
		case key.Matches(msg, m.byYear):
			if m.inc == date.Annually {
				return m, nil
			}
			m.inc = date.Annually
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
	case date.Weekly:
		m.startDate = m.startDate.AddDate(0, 0, 7)
		m.endDate = m.endDate.AddDate(0, 0, 7)
	case date.Monthly:
		m.startDate = m.startDate.AddDate(0, 1, 0)
		m.endDate = m.endDate.AddDate(0, 1, 0)
	case date.Quarterly:
		m.startDate = m.startDate.AddDate(0, 3, 0)
		m.endDate = m.endDate.AddDate(0, 3, 0)
	case date.Annually:
		m.startDate = m.startDate.AddDate(1, 0, 0)
		m.endDate = m.endDate.AddDate(1, 0, 0)
	}
}

func (m *DatePickerModel) prevDateRange() {
	switch m.inc {
	case date.Weekly:
		m.startDate = m.startDate.AddDate(0, 0, -7)
		m.endDate = m.endDate.AddDate(0, 0, -7)
	case date.Monthly:
		m.startDate = m.startDate.AddDate(0, -1, 0)
		m.endDate = m.endDate.AddDate(0, -1, 0)
	case date.Quarterly:
		m.startDate = m.startDate.AddDate(0, -3, 0)
		m.endDate = m.endDate.AddDate(0, -3, 0)
	case date.Annually:
		m.startDate = m.startDate.AddDate(-1, 0, 0)
		m.endDate = m.endDate.AddDate(-1, 0, 0)
	}
}

func (m DatePickerModel) View() string {
	var leftStr strings.Builder
	var rightStr strings.Builder

	leftStr.WriteString(fmt.Sprintf("%s: < %s >", m.inc, date.FormatDateToIncrement(m.startDate, m.inc)))

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
	rightStr.WriteString(keyStyle.Render(m.byQuarter.Keys()[0]) + "uarterly")
	rightStr.WriteString(" ")
	rightStr.WriteString(keyStyle.Render(m.byYear.Keys()[0]) + "early")

	style := lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(
			fmt.Sprintf("Date range: %s - %s", m.startDate.Format("2006-01-02"), m.endDate.Format("2006-01-02")),
			m.width,
		)).
		BorderForeground(borderColor).
		Padding(0, 1).
		Margin(1, 0, 0).
		Width(m.width)
	spaces := m.width - hPadding*2 - lipgloss.Width(leftStr.String()) - lipgloss.Width(rightStr.String())

	return style.
		Render(leftStr.String() + strings.Repeat(" ", max(0, spaces)) + rightStr.String())
}

func (m *DatePickerModel) updateIncrement() {
	// Snap start and end dates to increment
	switch m.inc {
	case date.Weekly:
		m.startDate = date.FirstDayOfWeek(m.startDate)
		m.endDate = m.startDate.AddDate(0, 0, 7)
	case date.Monthly:
		m.startDate = date.FirstDayOfMonth(m.startDate)
		m.endDate = m.startDate.AddDate(0, 1, 0)
	case date.Quarterly:
		m.startDate = date.FirstDayOfQuarter(m.startDate)
		m.endDate = m.startDate.AddDate(0, 3, 0)
	case date.Annually:
		m.startDate = date.FirstDayOfYear(m.startDate)
		m.endDate = m.startDate.AddDate(1, 0, 0)
	}
}
