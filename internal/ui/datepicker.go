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
	minDate   time.Time
	maxDate   time.Time

	reset     key.Binding
	next      key.Binding
	prev      key.Binding
	byWeek    key.Binding
	byMonth   key.Binding
	byQuarter key.Binding
	byYear    key.Binding
	allTime   key.Binding
}

type DateRangeChangedMsg struct {
	Start time.Time // Inclusive
	End   time.Time // Exclusive
}

type DateIncrementChangedMsg struct {
	Inc date.Increment
}

func NewDatePickerModel() DatePickerModel {
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return DatePickerModel{
		// First day of the month
		startDate: currentMonth,
		endDate:   currentMonth.AddDate(0, 1, 0),
		inc:       date.Monthly,
		reset:     key.NewBinding(key.WithKeys("0")),
		next:      key.NewBinding(key.WithKeys("l", "right")),
		prev:      key.NewBinding(key.WithKeys("h", "left")),
		byWeek:    key.NewBinding(key.WithKeys("w")),
		byMonth:   key.NewBinding(key.WithKeys("m")),
		byQuarter: key.NewBinding(key.WithKeys("q")),
		byYear:    key.NewBinding(key.WithKeys("y")),
		allTime:   key.NewBinding(key.WithKeys("a")),
	}
}

func (m *DatePickerModel) SetWidth(w int) {
	m.width = w
}

func (m *DatePickerModel) SetLimits(minDate, maxDate time.Time) {
	m.minDate = minDate
	m.maxDate = maxDate
	m.clampDateRangeToLimits()
}

func (m *DatePickerModel) minStartDate() time.Time {
	return m.inc.FirstDayInIncrement(m.minDate)
}

func (m *DatePickerModel) maxEndDate() time.Time {
	return m.inc.AddIncrement(m.inc.FirstDayInIncrement(m.maxDate))
}

func (m DatePickerModel) Update(msg tea.Msg) (DatePickerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.reset):
			cmd = m.resetDateRange()
		case key.Matches(msg, m.prev):
			cmd = m.prevDateRange()
		case key.Matches(msg, m.next):
			cmd = m.nextDateRange()
		case key.Matches(msg, m.byWeek):
			cmd = m.updateIncrement(date.Weekly)
		case key.Matches(msg, m.byMonth):
			cmd = m.updateIncrement(date.Monthly)
		case key.Matches(msg, m.byQuarter):
			cmd = m.updateIncrement(date.Quarterly)
		case key.Matches(msg, m.byYear):
			cmd = m.updateIncrement(date.Annually)
		case key.Matches(msg, m.allTime):
			cmd = m.updateIncrement(date.AllTime)
		}
	}
	return m, cmd
}

func (m *DatePickerModel) SelectedDateRange() (time.Time, time.Time) {
	return m.startDate, m.endDate
}

func (m *DatePickerModel) ViewDateRange() string {
	switch m.inc {
	case date.Weekly:
		year, week := m.startDate.ISOWeek()
		return fmt.Sprintf("%d Week %02d", year, week)
	case date.Monthly:
		return m.startDate.Format("2006 January")
	case date.Quarterly:
		return fmt.Sprintf("%d Q%d", m.startDate.Year(), date.QuarterOfYear(m.startDate))
	case date.Annually:
		return m.startDate.Format("2006")
	case date.AllTime:
		return m.inc.String()
	default:
		panic(fmt.Sprintf("Unexpected date increment: %s", m.inc))
	}
}

func (m *DatePickerModel) resetDateRange() tea.Cmd {
	if m.inc == date.AllTime {
		return nil
	}
	// Reset to current date while keeping increment
	m.startDate = m.inc.FirstDayInIncrement(time.Now())
	m.endDate = m.inc.AddIncrement(m.startDate)
	m.clampDateRangeToLimits()
	return m.sendDateRangeChangedMsg()
}

func (m *DatePickerModel) nextDateRange() tea.Cmd {
	if m.inc == date.AllTime {
		return nil
	}
	if nextEndDate := m.inc.AddIncrement(m.endDate); !nextEndDate.After(m.maxEndDate()) {
		m.startDate = m.inc.AddIncrement(m.startDate)
		m.endDate = nextEndDate
		return m.sendDateRangeChangedMsg()
	} else {
		return nil
	}
}

func (m *DatePickerModel) prevDateRange() tea.Cmd {
	if m.inc == date.AllTime {
		return nil
	}
	if prevStartDate := m.inc.SubtractIncrement(m.startDate); !prevStartDate.Before(m.minStartDate()) {
		m.startDate = prevStartDate
		m.endDate = m.inc.SubtractIncrement(m.endDate)
		return m.sendDateRangeChangedMsg()
	} else {
		return nil
	}
}

func (m DatePickerModel) View() string {
	var leftStr strings.Builder
	var rightStr strings.Builder

	// Current date increment and current date range selection
	if m.inc == date.AllTime {
		leftStr.WriteString(fmt.Sprintf("%s", m.ViewDateRange()))
	} else {
		leftStr.WriteString(fmt.Sprintf("%s: < %s >", m.inc, m.ViewDateRange()))
	}

	// Key bindings
	rightStr.WriteString(keyStyle.Render("h/←") + " prev")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render("l/→") + " next")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render("0") + " now")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render(m.byWeek.Keys()[0]) + "eekly")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render(m.byMonth.Keys()[0]) + "onthly")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render(m.byQuarter.Keys()[0]) + "uarterly")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render(m.byYear.Keys()[0]) + "early")
	rightStr.WriteString(" | ")
	rightStr.WriteString(keyStyle.Render(m.allTime.Keys()[0]) + "ll time")

	style := lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(
			fmt.Sprintf("Date range: %s - %s", m.startDate.Format(time.DateOnly), m.endDate.Format(time.DateOnly)),
			m.width,
		)).
		BorderForeground(borderColor).
		Padding(0, 1).
		Margin(1, 0, 0).
		Width(m.width)

	// Add spces to align rightStr to right side
	spaces := m.width - hPadding*2 - lipgloss.Width(leftStr.String()) - lipgloss.Width(rightStr.String())
	return style.
		Render(leftStr.String() + strings.Repeat(" ", max(0, spaces)) + rightStr.String())
}

func (m *DatePickerModel) Inc() date.Increment {
	return m.inc
}

func (m *DatePickerModel) updateIncrement(newInc date.Increment) tea.Cmd {
	if m.inc == newInc {
		return nil
	}

	m.inc = newInc
	if m.inc == date.AllTime {
		m.startDate = m.minDate
		m.endDate = m.maxDate
	} else {
		// Snap start and end dates to increment
		m.startDate = m.inc.FirstDayInIncrement(m.startDate)
		m.endDate = m.inc.AddIncrement(m.startDate)
		m.clampDateRangeToLimits()
	}
	return tea.Batch(m.sendIncrementChangedMsg(), m.sendDateRangeChangedMsg())
}

func (m *DatePickerModel) clampDateRangeToLimits() {
	if m.startDate.Before(m.minStartDate()) {
		m.startDate = m.minStartDate()
		m.endDate = m.inc.AddIncrement(m.startDate)
	} else if m.endDate.After(m.maxEndDate()) {
		m.endDate = m.maxEndDate()
		m.startDate = m.inc.SubtractIncrement(m.endDate)
	}
}

func (m *DatePickerModel) sendDateRangeChangedMsg() tea.Cmd {
	return func() tea.Msg {
		return DateRangeChangedMsg{
			Start: m.startDate,
			End:   m.endDate,
		}
	}
}

func (m *DatePickerModel) sendIncrementChangedMsg() tea.Cmd {
	return func() tea.Msg {
		return DateIncrementChangedMsg{
			Inc: m.inc,
		}
	}
}
