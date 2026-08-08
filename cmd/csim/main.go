package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/HT4w5/csim/pkg/cache"
	"github.com/HT4w5/csim/pkg/matrix"
	"github.com/HT4w5/csim/pkg/srv"
	"github.com/HT4w5/csim/pkg/transpose"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8081", "listen address")
	s := flag.Int("s", 5, "set index bits")
	E := flag.Int("E", 1, "lines per set")
	b := flag.Int("b", 5, "block offset bits")
	build := flag.String("build", "", "build report pages and exit if set")
	flag.Parse()

	cfg := &cache.Config{
		SetIndexBits:    *s,
		LinesPerSet:     *E,
		BlockOffsetBits: *b,
	}

	var reports []srv.Report
	var groups []srv.SizeGroup

	for _, size := range transpose.Sizes {
		rows, cols := size[0], size[1]
		var sols []srv.SolutionRow

		for _, fn := range transpose.Functions {
			path := fmt.Sprintf("/%s-%dx%d", fn.Name, rows, cols)
			row := srv.SolutionRow{Name: fn.Name, Link: path}

			sol, err := matrix.NewSolution(&matrix.SolutionConfig{
				Name:  fn.Name,
				Cache: cfg,
				Rows:  rows,
				Cols:  cols,
				TFunc: fn.F,
			})
			if err != nil {
				row.Error = err.Error()
			} else if err := evaluateRecover(sol); err != nil {
				row.Error = err.Error()
			} else {
				row.Valid = true
				row.Misses = sol.TotalMisses()
				var buf bytes.Buffer
				if err := sol.WriteReport(&buf); err != nil {
					row.Valid = false
					row.Error = err.Error()
				} else {
					reports = append(reports, srv.Report{Path: path, HTML: buf.Bytes()})
				}
			}

			if !row.Valid {
				reports = append(reports, srv.Report{Path: path, HTML: srv.ErrorPage(path, row.Error)})
			}

			fmt.Println(path)
			sols = append(sols, row)
		}

		sort.Slice(sols, func(i, j int) bool {
			if sols[i].Valid != sols[j].Valid {
				return sols[i].Valid
			}
			return sols[i].Misses < sols[j].Misses
		})

		groups = append(groups, srv.SizeGroup{Rows: rows, Cols: cols, Solutions: sols})
	}

	summaryData := srv.SummaryData{
		S:      *s,
		E:      *E,
		B:      *b,
		Groups: groups,
	}

	if *build != "" {
		if err := srv.WriteStatic(*build, reports, summaryData); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("wrote %d report pages + summary to %s\n", len(reports), *build)
		return
	}

	mux := srv.NewMux(reports, summaryData)

	fmt.Printf("listening on http://%s\n", *listen)
	if err := srv.Serve(*listen, mux); err != nil {
		fmt.Println(err)
	}
}

func evaluateRecover(sol *matrix.Solution) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return sol.Evaluate()
}
