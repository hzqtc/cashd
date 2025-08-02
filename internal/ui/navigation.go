package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMode string

const (
	TransactionView ViewMode = "Transactions"
	AccountView     ViewMode = "Accounts"
	CategoryView    ViewMode = "Categories"
)

type NavBarModel struct {
	width    int
	viewMode ViewMode

	navTransactionView key.Binding
	navAccountView     key.Binding
	navCategoryView    key.Binding
}

type NavigationMsg struct {
	viewMode ViewMode
}

func NewNavBarModel() NavBarModel {
	return NavBarModel{
		viewMode:           TransactionView,
		navTransactionView: key.NewBinding(key.WithKeys("1")),
		navAccountView:     key.NewBinding(key.WithKeys("2")),
		navCategoryView:    key.NewBinding(key.WithKeys("3")),
	}
}

func (m *NavBarModel) SetWidth(w int) {
	m.width = w
}

func (m *NavBarModel) ViewMode() ViewMode {
	return m.viewMode
}

func (m NavBarModel) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s %s", keyStyle.Render(m.navTransactionView.Keys()[0]), TransactionView))
	s.WriteString(" ")
	s.WriteString(fmt.Sprintf("%s %s", keyStyle.Render(m.navAccountView.Keys()[0]), AccountView))
	s.WriteString(" ")
	s.WriteString(fmt.Sprintf("%s %s", keyStyle.Render(m.navCategoryView.Keys()[0]), CategoryView))

	style := lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(fmt.Sprintf("View: %s", m.viewMode), m.width)).
		BorderForeground(borderColor).
		Padding(0, 1).
		Margin(1, 0, 0).
		Width(m.width)
	return style.Render(s.String())
}

func (m NavBarModel) Update(msg tea.Msg) (NavBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.navTransactionView):
			m.viewMode = TransactionView
			return m, m.getNavMsg()
		case key.Matches(msg, m.navAccountView):
			m.viewMode = AccountView
			return m, m.getNavMsg()
		case key.Matches(msg, m.navCategoryView):
			m.viewMode = CategoryView
			return m, m.getNavMsg()
		}
	}
	return m, nil
}

func (m *NavBarModel) getNavMsg() tea.Cmd {
	return func() tea.Msg {
		return NavigationMsg{
			viewMode: m.viewMode,
		}
	}
}
