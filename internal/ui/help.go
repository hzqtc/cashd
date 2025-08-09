package ui

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

var hideHelp bool

func init() {
	pflag.BoolVar(&hideHelp, "hide-help", false, "Hide the help panel in the TUI")
}

type HelpModel struct {
	visible bool
	width   int
}

func NewHelpModel() HelpModel {
	return HelpModel{
		visible: !hideHelp,
	}
}

func (m *HelpModel) ToggleVisibility() {
	m.visible = !m.visible
}

func (m *HelpModel) Visible() bool {
	return m.visible
}

func (m *HelpModel) SetWidth(width int) {
	m.width = width
}

func (m HelpModel) View() string {
	if !m.visible {
		return ""
	}
	var s strings.Builder
	s.WriteString(fmt.Sprintf(
		"General: %s quit | %s toggle help | %s search transactiona | %s clear search\n",
		keyStyle.Render("ctrl + c"),
		keyStyle.Render("?"),
		keyStyle.Render("/"),
		keyStyle.Render("esc"),
	))
	s.WriteString(fmt.Sprintf(
		"Table: %s down | %s up | %s pgDown | %s pgUp | %s top | %s bottom\n",
		keyStyle.Render("j/↓"),
		keyStyle.Render("k/↑"),
		keyStyle.Render("PgDn"),
		keyStyle.Render("PgUp"),
		keyStyle.Render("g"),
		keyStyle.Render("G"),
	))
	s.WriteString(fmt.Sprintf(
		"Sorting: %s sort by next column | %s sort by prev column | %s reverse sort direction",
		keyStyle.Render("s"),
		keyStyle.Render("S"),
		keyStyle.Render("r"),
	))
	return baseStyle.Width(m.width).Render(s.String())
}
