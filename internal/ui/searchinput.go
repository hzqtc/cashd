package ui

import (
	"log"
	"strings"

	"cashd/internal/data"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchInputModel struct {
	width         int
	showNameInput bool
	showTable     bool

	input     textinput.Model
	nameInput textinput.Model
	table     table.Model

	cancel       key.Binding
	enter        key.Binding
	saveSearch   key.Binding
	loadSearch   key.Binding
	deleteSearch key.Binding
}

type SearchMsg struct {
	query string
}

const nameLength = 20

func NewSearchInputModel() SearchInputModel {
	input := textinput.New()
	input.Placeholder = "Search..."
	input.Prompt = " / "

	nameInput := textinput.New()
	nameInput.Prompt = "Name this search: "
	nameInput.CharLimit = nameLength

	return SearchInputModel{
		input:     input,
		nameInput: nameInput,
		table: table.New(
			table.WithFocused(true),
			table.WithHeight(5),
			table.WithStyles(getTableStyle()),
		),

		showNameInput: false,
		showTable:     false,

		cancel:       key.NewBinding(key.WithKeys("esc")),
		enter:        key.NewBinding(key.WithKeys("enter")),
		saveSearch:   key.NewBinding(key.WithKeys("ctrl+s")),
		loadSearch:   key.NewBinding(key.WithKeys("ctrl+l")),
		deleteSearch: key.NewBinding(key.WithKeys("ctrl+d")),
	}
}

func (m *SearchInputModel) SetWidth(w int) {
	m.width = w
	columns := []table.Column{
		{Title: "Name", Width: nameLength},
		{Title: "Query", Width: w - nameLength - 4},
	}
	m.table.SetColumns(columns)
	m.table.SetWidth(w)
	m.updateLayout()
}

func (m *SearchInputModel) updateLayout() {
	if m.showNameInput {
		nameInputMaxWidth := len(m.nameInput.Prompt) + nameLength + 1 // 1 for width of the cursor
		m.input.Width = m.width - 4 - nameInputMaxWidth
		m.input.SetCursor(m.input.Position()) // Trigger textinput.handleOverflow
	} else {
		m.input.Width = m.width - 4 - lipgloss.Width(m.renderHelp())
	}
}

func (m *SearchInputModel) Focused() bool {
	return m.input.Focused() || m.showNameInput || m.showTable
}

func (m *SearchInputModel) Focus() {
	m.input.Focus()
	m.updateLayout()
}

func (m *SearchInputModel) Blur() {
	m.input.Blur()
	m.showTable = false
	m.showNameInput = false
	m.updateLayout()
}

func (m *SearchInputModel) Value() string {
	return m.input.Value()
}

func (m *SearchInputModel) Clear() tea.Cmd {
	m.input.SetValue("")
	return m.sendSearchMsg()
}

func (m SearchInputModel) View() string {
	var s string
	if m.showNameInput {
		s = lipgloss.JoinHorizontal(lipgloss.Top, m.input.View(), m.nameInput.View())
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Top, m.input.View(), m.renderHelp())
	}
	s = baseStyle.Width(m.width).Render(s)

	if m.showTable {
		s = lipgloss.JoinVertical(lipgloss.Left, s, baseStyle.Render(m.table.View()))
	}

	return s
}

func (m *SearchInputModel) renderHelp() string {
	var s strings.Builder
	if m.showTable {
		s.WriteString(keyStyle.Render("⏎") + " select")
		s.WriteString(" | ")
		s.WriteString(keyStyle.Render("⎋") + " hide")
		s.WriteString(" | ")
		s.WriteString(keyStyle.Render("^d") + " delete")
	} else if m.input.Focused() {
		s.WriteString(keyStyle.Render("⏎") + " search")
		s.WriteString(" | ")
		if strings.TrimSpace(m.input.Value()) != "" {
			s.WriteString(keyStyle.Render("^s") + " save")
			s.WriteString(" | ")
		}
		s.WriteString(keyStyle.Render("^l") + " load saved")
	}
	return s.String()
}

func (m *SearchInputModel) refreshSavedSearch() {
	if searches, err := data.LoadSavedSearches(); err != nil {
		log.Printf("Error loading saved searches: %v", err)
	} else {
		rows := make([]table.Row, len(searches))
		for i, s := range searches {
			rows[i] = table.Row{s.Name, s.Query}
		}
		m.table.SetRows(rows)
	}
}

func (m SearchInputModel) Update(msg tea.Msg) (SearchInputModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.showTable:
			return m, m.handleTableKeys(msg)
		case m.nameInput.Focused():
			return m, m.handleNameInputKeys(msg)
		case m.input.Focused():
			return m, m.handleSearchInputKeys(msg)
		}
	}
	return m, nil
}

func (m *SearchInputModel) handleSearchInputKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, m.cancel):
		m.input.SetValue("")
		fallthrough
	case key.Matches(msg, m.enter):
		m.Blur()
		cmd = m.sendSearchMsg()
	case key.Matches(msg, m.saveSearch):
		m.input.Blur()
		m.showNameInput = true
		m.updateLayout()
		m.nameInput.Focus()
	case key.Matches(msg, m.loadSearch):
		m.input.Blur()
		m.refreshSavedSearch()
		m.showTable = true
	default:
		m.input, cmd = m.input.Update(msg)
	}
	return cmd
}

func (m *SearchInputModel) handleNameInputKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, m.enter):
		searchName := strings.TrimSpace(m.nameInput.Value())
		if searchName == "" {
			break
		}
		err := data.AddOrUpdateSavedSearch(searchName, m.input.Value())
		if err != nil {
			log.Printf("Error saving search: %v", err)
			break
		}
		m.refreshSavedSearch()
		fallthrough
	case key.Matches(msg, m.cancel):
		m.showNameInput = false
		m.updateLayout()
		m.nameInput.Blur()
		m.nameInput.SetValue("")
		m.input.Focus()
	default:
		m.nameInput, cmd = m.nameInput.Update(msg)
	}
	return cmd
}

func (m *SearchInputModel) handleTableKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case key.Matches(msg, m.enter):
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 1 {
			name := selectedRow[0]
			query := selectedRow[1]
			m.input.SetValue(query)
			m.Blur()
			data.AddOrUpdateSavedSearch(name, query) // Update timestamp
			m.refreshSavedSearch()
			cmd = m.sendSearchMsg()
		}
	case key.Matches(msg, m.cancel):
		m.showTable = false
		m.input.Focus()
	case key.Matches(msg, m.deleteSearch):
		selectedRow := m.table.SelectedRow()
		if len(selectedRow) > 1 {
			name := selectedRow[0]
			data.DeleteSavedSearch(name)
			m.refreshSavedSearch()
		}
	default:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m *SearchInputModel) sendSearchMsg() tea.Cmd {
	return func() tea.Msg {
		return SearchMsg{
			query: m.input.Value(),
		}
	}
}
