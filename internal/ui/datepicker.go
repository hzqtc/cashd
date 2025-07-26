package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type DatePickerModel struct {
	currentMonth time.Time
}

type DateRangeChangedMsg struct {
	Start time.Time // Inclusive
	End   time.Time // Exclusive
}

func NewDatePickerModel() DatePickerModel {
	now := time.Now()
	return DatePickerModel{
		// First day of the month
		currentMonth: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
	}
}

func (m DatePickerModel) Update(msg tea.Msg) (DatePickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.currentMonth = m.currentMonth.AddDate(0, -1, 0)
			return m, onDateRangeChange(m.GetSelectedDateRange())
		case "l", "right":
			m.currentMonth = m.currentMonth.AddDate(0, 1, 0)
			return m, onDateRangeChange(m.GetSelectedDateRange())
		}
	}
	return m, nil
}

func (m DatePickerModel) View() string {
	return datepickerStyle.Render(fmt.Sprintf(" < %s > ", m.currentMonth.Format("January 2006")))
}

func (m DatePickerModel) GetSelectedDateRange() (time.Time, time.Time) {
	startOfMonth := m.currentMonth
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	return startOfMonth, endOfMonth
}

func onDateRangeChange(startDate, endDate time.Time) tea.Cmd {
	return func() tea.Msg {
		return DateRangeChangedMsg{
			Start: startDate,
			End:   endDate,
		}
	}
}
