package srv

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func WriteStatic(outdir string, reports []Report, data SummaryData) error {
	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return err
	}

	for _, r := range reports {
		name := strings.TrimPrefix(r.Path, "/") + ".html"
		path := filepath.Join(outdir, name)
		if dir := filepath.Dir(path); dir != outdir {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
		}
		if err := os.WriteFile(path, r.HTML, 0o644); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if err := RenderSummary(&buf, staticSummary(data)); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outdir, "index.html"), buf.Bytes(), 0o644)
}

func staticSummary(data SummaryData) SummaryData {
	out := data
	out.Groups = make([]SizeGroup, len(data.Groups))
	for i, g := range data.Groups {
		out.Groups[i] = g
		out.Groups[i].Solutions = make([]SolutionRow, len(g.Solutions))
		for j, s := range g.Solutions {
			out.Groups[i].Solutions[j] = s
			out.Groups[i].Solutions[j].Link = strings.TrimPrefix(s.Link, "/") + ".html"
		}
	}
	return out
}
