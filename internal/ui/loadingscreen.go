package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

var showTimer bool

func init() {
	pflag.BoolVar(&showTimer, "show-timer", false, "Show a stop watch on the loading screen")
}

const logo = `
 _______  _______  _______           ______
(  ____ \(  ___  )(  ____ \|\     /|(  __  \
| (    \/| (   ) || (    \/| )   ( || (  \  )
| |      | (___) || (_____ | (___) || |   ) |
| |      |  ___  |(_____  )|  ___  || |   | |
| |      | (   ) |      ) || (   ) || |   ) |
| (____/\| )   ( |/\____) || )   ( || (__/  )
(_______/|/     \|\_______)|/     \|(______/
`

type LoadingScreenModel struct {
	loading   bool
	spinner   spinner.Model
	stopwatch stopwatch.Model
}

func NewLoadingScreenModel() LoadingScreenModel {
	s := spinner.New()
	s.Spinner = spinner.Meter

	var sw stopwatch.Model
	if showTimer {
		sw = stopwatch.NewWithInterval(time.Millisecond)
	}

	return LoadingScreenModel{
		loading:   true,
		spinner:   s,
		stopwatch: sw,
	}
}

func (m *LoadingScreenModel) Handles(msg tea.Msg) bool {
	switch msg.(type) {
	case spinner.TickMsg, stopwatch.TickMsg, stopwatch.StartStopMsg, stopwatch.ResetMsg:
		return true
	default:
		return false
	}
}

func (m LoadingScreenModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}
	if showTimer {
		cmds = append(cmds, m.stopwatch.Start())
	}
	return tea.Batch(cmds...)
}

func (m LoadingScreenModel) Update(msg tea.Msg) (LoadingScreenModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case stopwatch.TickMsg:
		if m.loading {
			m.stopwatch, cmd = m.stopwatch.Update(msg)
			return m, cmd
		}
	case stopwatch.StartStopMsg, stopwatch.ResetMsg:
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *LoadingScreenModel) Stop() tea.Cmd {
	m.loading = false
	return tea.Batch(m.stopwatch.Stop(), m.stopwatch.Reset())
}

func (m *LoadingScreenModel) IsLoading() bool {
	return m.loading
}

func (m LoadingScreenModel) View() string {
	var b strings.Builder
	m.spinner.Style = keyStyle
	b.WriteString(fmt.Sprintf("%s\n%s Loading...", keyStyle.Render(logo), m.spinner.View()))
	if showTimer {
		b.WriteString(m.stopwatch.View())
	}
	return b.String()
}
