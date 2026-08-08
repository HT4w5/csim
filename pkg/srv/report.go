package srv

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
)

type Report struct {
	Path string
	HTML []byte
}

func NewMux(reports []Report, data SummaryData) *http.ServeMux {
	mux := http.NewServeMux()

	for _, r := range reports {
		r := r
		mux.HandleFunc(r.Path, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(r.HTML)
		})
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		RenderSummary(w, data)
	})

	return mux
}

//go:embed error.tmpl
var errorTmplStr string

var errorTmpl = template.Must(template.New("error").Parse(errorTmplStr))

func ErrorPage(title, msg string) []byte {
	var buf bytes.Buffer
	errorTmpl.Execute(&buf, struct {
		Title   string
		Message string
	}{title, msg})
	return buf.Bytes()
}
