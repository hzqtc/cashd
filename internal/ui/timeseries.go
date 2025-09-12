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
}

func NewTimeSeriesChartModel() TimeSeriesChartModel {
	return TimeSeriesChartModel{}
}

func (m *TimeSeriesChartModel) SetDimension(width, height int) {
	m.width = width
	m.height = height
	if len(m.entries) > 0 {
		m.chart.Resize(width, height)
		m.chart.DrawBrailleAll()
	}
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
		tschart.WithDataSetStyle(string(data.Expense), tsChartExpenseLineStyle), // Expense style
		tschart.WithXLabelFormatter(dateLabelFormatter(m.inc)),
		tschart.WithYLabelFormatter(moneyAmountFormatter),
	)

	// Push data to the respective datasets
	for _, entry := range m.entries {
		m.chart.PushDataSet(string(data.Income), tschart.TimePoint{Time: entry.Date, Value: entry.Income})
		m.chart.PushDataSet(string(data.Expense), tschart.TimePoint{Time: entry.Date, Value: entry.Expense})
	}
	// Limit the X range
	m.chart.SetViewTimeRange(m.entries[0].Date, m.entries[len(m.entries)-1].Date)

	m.chart.DrawBrailleAll()
}

func (m TimeSeriesChartModel) View() string {
	return lipgloss.NewStyle().
		Border(getRoundedBorderWithTitle(m.name, m.width+hPadding*2)).
		BorderForeground(borderColor).
		Padding(vPadding, hPadding).
		Render(m.renderLegend() + m.chart.View())
}

func (m TimeSeriesChartModel) renderLegend() string {
	return fmt.Sprintf(
		"\n%s %s    %s %s\n\n",
		tsChartIncomeLineStyle.Render(string(runes.FullBlock)),
		data.Income,
		tsChartExpenseLineStyle.Render(string(runes.FullBlock)),
		data.Expense,
	)
}

func moneyAmountFormatter(i int, v float64) string {
	return data.FormatMoneyInteger(math.Round(v/10) * 10)
}

func dateLabelFormatter(inc date.Increment) linechart.LabelFormatter {
	return func(i int, v float64) string {
		d := time.Unix(int64(v), 0).Local()
		switch inc {
		case date.Weekly:
			_, week := d.ISOWeek()
			return fmt.Sprintf("%s'W%02d", d.Format("06"), week)
		case date.Monthly:
			return d.Format("06'Jan")
		case date.Quarterly:
			return fmt.Sprintf("%s'Q%d", d.Format("06"), date.QuarterOfYear(d))
		case date.Annually, date.AllTime:
			return d.Format("2006")
		default:
			panic(fmt.Sprintf("Unexpected date increment: %s", inc))
		}
	}
}
