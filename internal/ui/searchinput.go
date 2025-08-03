package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchInputModel struct {
	input textinput.Model

	cancelSearch key.Binding
	search       key.Binding
}

type SearchMsg struct {
	query string
}

func NewSearchInputModel() SearchInputModel {
	input := textinput.New()
	input.Placeholder = "Search..."
	input.Prompt = " / "

	return SearchInputModel{
		input:        input,
		cancelSearch: key.NewBinding(key.WithKeys("esc")),
		search:       key.NewBinding(key.WithKeys("enter")),
	}
}

func (m *SearchInputModel) SetWidth(w int) {
	m.input.Width = w
}

func (m *SearchInputModel) Focused() bool {
	return m.input.Focused()
}

func (m *SearchInputModel) Focus() {
	m.input.Focus()
}

func (m *SearchInputModel) Blur() {
	m.input.Blur()
}

func (m *SearchInputModel) Value() string {
	return m.input.Value()
}

func (m *SearchInputModel) Clear() tea.Cmd {
	m.input.SetValue("")
	return m.sendSearchMsg()
}

func (m SearchInputModel) View() string {
	return searchStyle.Render(m.input.View())
}

func (m SearchInputModel) Update(msg tea.Msg) (SearchInputModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.cancelSearch):
			m.input.SetValue("")
			m.input.Blur()
			cmd = m.sendSearchMsg()
		case key.Matches(msg, m.search):
			m.input.Blur()
			cmd = m.sendSearchMsg()
		default:
			m.input, cmd = m.input.Update(msg)
		}
	default:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m *SearchInputModel) sendSearchMsg() tea.Cmd {
	return func() tea.Msg {
		return SearchMsg{
			query: m.input.Value(),
		}
	}
}
