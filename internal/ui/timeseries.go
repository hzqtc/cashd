package ui

import (
	"cashd/internal/data"
	"cashd/internal/date"
	"fmt"
	"math"
	"time"

	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart"
	tschart "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
)

type TsChartEntry struct {
	Date    time.Time
	Inc     date.Increment
	Income  float64
	Expense float64
}

type TimeSeriesChartModel struct {
	width  int
	height int

	name    string
	inc     date.Increment
	entries []*TsChartEntry

	chart tschart.Model

	// TODO: support key bindings for horizontal scrolling
}

func NewTimeSeriesChartModel() TimeSeriesChartModel {
	return TimeSeriesChartModel{}
}

func (m *TimeSeriesChartModel) SetDimension(width, height int) {
	m.width = width
	m.height = height
	m.chart.Resize(width, height)
	m.chart.DrawAll()
}

func (m *TimeSeriesChartModel) SetEntries(name string, entries []*TsChartEntry, inc date.Increment) {
	m.name = name
	m.entries = entries
	m.inc = inc
	m.redraw()
}

// Draw a timeseries chart with 2 lines for incomes and expenses each
func (m *TimeSeriesChartModel) redraw() {
	if len(m.entries) == 0 || m.width == 0 || m.height == 0 {
		return
	}

	var maxValue float64
	for _, entry := range m.entries {
		maxValue = max(maxValue, entry.Income)
		maxValue = max(maxValue, entry.Expense)
	}

	// Create a new chart on data set change, not worth reusing the model
	m.chart = tschart.New(
		m.width,
		m.height,
		tschart.WithYRange(0, maxValue),
		tschart.WithAxesStyles(tsChartAxisStyle, tsChartLabelStyle),
		tschart.WithDataSetStyle(string(data.Income), tsChartIncomeLineStyle),   // Income style
		tschart.WithDataSetLineStyle(string(data.Income), runes.ThinLineStyle),  // Income line style
		tschart.WithDataSetStyle(string(data.Expense), tsChartExpenseLineStyle), // Expense style
		tschart.WithDataSetLineStyle(string(data.Expense), runes.ArcLineStyle),  // Expense line style
		tschart.WithXLabelFormatter(dateLabelFormatter(m.inc)),
		tschart.WithYLabelFormatter(moneyAmountFormatter()),
	)

	// Push data to the respective datasets
	// TODO: add line legend
	for _, entry := range m.entries {
		m.chart.PushDataSet(string(data.Income), tschart.TimePoint{Time: entry.Date, Value: entry.Income})
		m.chart.PushDataSet(string(data.Expense), tschart.TimePoint{Time: entry.Date, Value: entry.Expense})
	}
	// Limit the X range
	m.chart.SetViewTimeRange(m.entries[0].Date, m.entries[len(m.entries)-1].Date)

	m.chart.DrawAll()
}

func (m TimeSeriesChartModel) View() string {
	return lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(m.name, m.width+hPadding*2)).
		BorderForeground(borderColor).
		Padding(vPadding, hPadding).
		Render(m.chart.View())
}

func moneyAmountFormatter() linechart.LabelFormatter {
	return func(i int, v float64) string {
		return fmt.Sprintf("$%.0f", math.Round(v/10)*10)
	}
}

func dateLabelFormatter(inc date.Increment) linechart.LabelFormatter {
	return func(i int, v float64) string {
		date := time.Unix(int64(v), 0).Local()
		return inc.FormatDateShorter(date)
	}
}
