package matrix

import (
	"fmt"
	"io"

	"github.com/HT4w5/csim/pkg/cache"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	accentHit      = "#859900"
	accentMiss     = "#dc322f"
	accentEviction = "#b58900"
	accentBase     = "#fdf6e3"
)

func (s *Solution) WriteReport(w io.Writer) error {
	page := components.NewPage()
	page.SetPageTitle(fmt.Sprintf("%s:%dx%d", s.cfg.Name, s.cfg.Rows, s.cfg.Cols))
	page.AddCharts(
		heatMap(s.a, "Matrix A", "hits"),
		heatMap(s.a, "Matrix A", "misses"),
		heatMap(s.a, "Matrix A", "evictions"),
		heatMap(s.b, "Matrix B", "hits"),
		heatMap(s.b, "Matrix B", "misses"),
		heatMap(s.b, "Matrix B", "evictions"),
		cacheHMEPerSet(s),
		cacheEfficiencyPerLine(s.c),
	)

	return page.Render(w)
}

func heatMap(mx *matrix, name, metric string) *charts.HeatMap {
	hm := charts.NewHeatMap()
	data := mx.genData(metric)

	var accent string

	switch metric {
	case "hits":
		accent = accentHit
	case "misses":
		accent = accentMiss
	case "evictions":
		accent = accentEviction
	}

	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: name + ": " + metric,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "X",
			Type: "category",
			Data: axis(mx.rows),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Y",
			Type: "category",
			Data: axis(mx.cols),
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Type: "continuous",
			Min:  0,
			Max:  dataMax(data),
			InRange: &opts.VisualMapInRange{
				Color: []string{accentBase, accent},
			},
		}),
	)

	hm.AddSeries(metric, data, charts.WithItemStyleOpts(opts.ItemStyle{Color: accent}))

	return hm
}

func axis(count int) []int {
	out := make([]int, 0, count)
	for i := range count {
		out = append(out, i)
	}
	return out
}

func dataMax(items []opts.HeatMapData) float32 {
	max := float32(0)
	for _, it := range items {
		t, ok := it.Value.([3]any)
		if !ok {
			continue
		}
		if v, ok := t[2].(int); ok && float32(v) > max {
			max = float32(v)
		}
	}
	if max < 1 {
		return 1
	}
	return max
}

func (mx *matrix) genData(metric string) []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0, len(mx.data))
	for i := range len(mx.stats) {
		var v int
		switch metric {
		case "hits":
			v = mx.stats[i].Hits
		case "misses":
			v = mx.stats[i].Misses
		case "evictions":
			v = mx.stats[i].Evictions
		}
		items = append(items, opts.HeatMapData{
			Value: [3]any{i / mx.cols, i % mx.cols, v},
		})
	}
	return items
}

func cacheHMEPerSet(s *Solution) *charts.Bar {
	bar := charts.NewBar()

	numSets := s.c.NumSets()
	hits := make([]opts.BarData, 0, numSets)
	misses := make([]opts.BarData, 0, numSets)
	evictions := make([]opts.BarData, 0, numSets)

	for _, st := range s.c.HMEStats() {
		hits = append(hits, opts.BarData{Value: st.Hits})
		misses = append(misses, opts.BarData{Value: st.Misses})
		evictions = append(evictions, opts.BarData{Value: st.Evictions})
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Cache HME per Set",
			Subtitle: fmt.Sprintf("Total misses: %d", s.TotalMisses()),
		}),
		charts.WithXAxisOpts(opts.XAxis{Type: "category"}),
		charts.WithYAxisOpts(opts.YAxis{Type: "value", Min: 0}),
	)
	bar.SetXAxis(axis(numSets))
	bar.AddSeries("Hits", hits, charts.WithItemStyleOpts(opts.ItemStyle{Color: accentHit}))
	bar.AddSeries("Misses", misses, charts.WithItemStyleOpts(opts.ItemStyle{Color: accentMiss}))
	bar.AddSeries("Evictions", evictions, charts.WithItemStyleOpts(opts.ItemStyle{Color: accentEviction}))

	return bar
}

func cacheEfficiencyPerLine(sim *cache.Simulator) *charts.Bar {
	bar := charts.NewBar()

	numSets := sim.NumSets()
	linesPerSet := sim.LinesPerSet()

	series := make([][]opts.BarData, linesPerSet)
	for j := range linesPerSet {
		series[j] = make([]opts.BarData, numSets)
		for i := range numSets {
			series[j][i] = opts.BarData{Value: 0}
		}
	}

	for lineIdx, eff := range sim.EfficiencyStats() {
		set := lineIdx / linesPerSet
		line := lineIdx % linesPerSet
		series[line][set] = opts.BarData{Value: eff}
	}

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Cache Efficiency per Line (by Set)"}),
		charts.WithXAxisOpts(opts.XAxis{Type: "category"}),
		charts.WithYAxisOpts(opts.YAxis{Type: "value", Min: 0, Max: 1}),
	)
	bar.SetXAxis(axis(numSets))

	palette := []string{"#268bd2", "#2aa198", "#b58900", "#cb4b16", "#d33682", "#6c71c4"}
	for j := range linesPerSet {
		bar.AddSeries(fmt.Sprintf("Line %d", j), series[j], charts.WithItemStyleOpts(opts.ItemStyle{Color: palette[j%len(palette)]}))
	}

	return bar
}
