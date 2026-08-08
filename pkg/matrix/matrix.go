/*
Package matrix provides the framework to writing and evaluating matrix transpose functions
*/
package matrix

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/HT4w5/csim/pkg/cache"
)

type TransposeFunc func(n, m int, a, b *Matrix)

type HMEStats struct {
	Hits      int
	Misses    int
	Evictions int
}

type Matrix struct {
	mx *matrix
}

func (mx *Matrix) Read(x, y int) int32 {
	return mx.mx.read2D(x, y)
}

func (mx *Matrix) Write(x, y int, v int32) {
	mx.mx.write2D(x, y, v)
}

type matrix struct {
	rows     int
	cols     int
	baseAddr uint64
	cache    *cache.Simulator
	data     []int32
	stats    []HMEStats
}

func newMatrix(cache *cache.Simulator, rows, cols int, baseAddr uint64) *matrix {
	length := rows * cols

	return &matrix{
		rows:     rows,
		cols:     cols,
		baseAddr: baseAddr,
		cache:    cache,
		data:     make([]int32, length),
		stats:    make([]HMEStats, length),
	}
}

// Equivelent to int[x][y] in C.
func (mx *matrix) read2D(x, y int) int32 {
	return mx.read1D(x*mx.cols + y)
}

func (mx *matrix) read1D(idx int) int32 {
	h, m, e := mx.cache.Access(mx.baseAddr+uint64(idx*4), 4)
	stat := &mx.stats[idx]
	stat.Hits += h
	stat.Misses += m
	stat.Evictions += e
	return mx.data[idx]
}

func (mx *matrix) write2D(x, y int, v int32) {
	mx.write1D(x*mx.cols+y, v)
}

func (mx *matrix) write1D(idx int, v int32) {
	h, m, e := mx.cache.Access(mx.baseAddr+uint64(idx*4), 4)
	stat := &mx.stats[idx]
	stat.Hits += h
	stat.Misses += m
	stat.Evictions += e
	mx.data[idx] = v
}

func (mx *matrix) populate() {
	for i := range len(mx.data) {
		mx.data[i] = rand.Int32()
	}
}

func (mx *matrix) isTransposeOf(other *matrix) error {
	if mx.cols != other.rows || mx.rows != other.cols {
		return fmt.Errorf("matrix size mismatch: %dx%d, %dx%d", mx.rows, mx.cols, other.rows, other.cols)
	}

	for i := range len(mx.data) {
		if mx.data[i] != other.data[(i%mx.cols)*mx.rows+i/mx.cols] {
			return fmt.Errorf("matrix cell mismatch at %d", i)
		}
	}

	return nil
}

type Solution struct {
	cfg *SolutionConfig
	c   *cache.Simulator
	a   *matrix
	b   *matrix
	f   TransposeFunc
}

type SolutionConfig struct {
	Name  string
	Cache *cache.Config
	Rows  int
	Cols  int
	TFunc TransposeFunc
}

func NewSolution(cfg *SolutionConfig) (*Solution, error) {
	if cfg.Rows <= 0 || cfg.Cols <= 0 {
		return nil, fmt.Errorf("matrix.NewSolution: invalid matrix size: %dx%d", cfg.Rows, cfg.Cols)
	}

	if cfg.TFunc == nil {
		return nil, errors.New("matrix.NewSolution: nil TFunc")
	}

	c, err := cache.New(cfg.Cache)
	if err != nil {
		return nil, fmt.Errorf("matrix.NewSolution: create cache simulator: %w", err)
	}

	a := newMatrix(c, cfg.Rows, cfg.Cols, 0)
	b := newMatrix(c, cfg.Cols, cfg.Rows, 256*256*4)

	return &Solution{
		cfg: cfg,
		c:   c,
		a:   a,
		b:   b,
		f:   cfg.TFunc,
	}, nil
}

func (s *Solution) Evaluate() error {
	s.a.populate()
	s.f(s.cfg.Rows, s.cfg.Cols, &Matrix{mx: s.a}, &Matrix{mx: s.b})
	if err := s.a.isTransposeOf(s.b); err != nil {
		return fmt.Errorf("Evaluate: validation failed: %w", err)
	}
	return nil
}

func (s *Solution) TotalMisses() int {
	total := 0
	for _, st := range s.c.HMEStats() {
		total += st.Misses
	}
	return total
}
