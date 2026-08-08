package srv

import (
	_ "embed"
	"html/template"
	"io"
)

type SolutionRow struct {
	Name   string
	Valid  bool
	Misses int
	Error  string
	Link   string
}

type SizeGroup struct {
	Rows      int
	Cols      int
	Solutions []SolutionRow
}

type SummaryData struct {
	S      int
	E      int
	B      int
	Groups []SizeGroup
}

//go:embed summary.tmpl
var summaryTmplStr string

var summaryTmpl = template.Must(template.New("summary").Funcs(template.FuncMap{
	"add": func(i int) int { return i + 1 },
}).Parse(summaryTmplStr))

func RenderSummary(w io.Writer, data SummaryData) error {
	return summaryTmpl.Execute(w, data)
}
