package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type DatePickerModel struct {
	currentMonth time.Time
}

func NewDatePickerModel() DatePickerModel {
	return DatePickerModel{
		currentMonth: time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1), // Start of current month
	}
}

func (m DatePickerModel) Init() tea.Cmd {
	return nil
}

func (m DatePickerModel) Update(msg tea.Msg) (DatePickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			m.currentMonth = m.currentMonth.AddDate(0, -1, 0)
			return m, DateRangeChanged(m.GetSelectedDateRange())
		case "l", "right":
			m.currentMonth = m.currentMonth.AddDate(0, 1, 0)
			return m, DateRangeChanged(m.GetSelectedDateRange())
		}
	}
	return m, nil
}

func (m DatePickerModel) View() string {
	monthStr := m.currentMonth.Format("January 2006")
	return datepickerStyle.
		Render(fmt.Sprintf(" < %s > ", monthStr))
}

func (m DatePickerModel) GetSelectedDateRange() (time.Time, time.Time) {
	startOfMonth := m.currentMonth
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	return startOfMonth, endOfMonth
}

type DateRangeChangedMsg struct {
	StartDate time.Time
	EndDate   time.Time
}

func DateRangeChanged(startDate, endDate time.Time) tea.Cmd {
	return func() tea.Msg {
		return DateRangeChangedMsg{
			StartDate: startDate,
			EndDate:   endDate,
		}
	}
}
