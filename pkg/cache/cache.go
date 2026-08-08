/*
Package cache implements a cache simulator compatible with
the CMU15-213 cachelab simulator.
*/
package cache

import (
	"errors"
	"fmt"
	"iter"
)

type cacheLine struct {
	tag        uint64
	lastAccess int
	valid      bool

	loadedBytes   int
	accessedBytes int
}

type cacheSet struct {
	cfg           *Config
	lines         []cacheLine
	accessCounter int

	hits      int
	misses    int
	evictions int
}

func (cs *cacheSet) init(linesPerSet int, cfg *Config) {
	cs.cfg = cfg
	cs.lines = make([]cacheLine, linesPerSet)
	cs.accessCounter = 0
	cs.hits = 0
	cs.misses = 0
	cs.evictions = 0
}

func (cs *cacheSet) access(tag uint64, size int) (hits int, misses int, evictions int) {
	cs.accessCounter++

	found := false
	available := (*cacheLine)(nil)
	oldest := &cs.lines[0]

	for i := range len(cs.lines) {
		line := &cs.lines[i]

		if !line.valid {
			available = line
			continue
		}

		if line.lastAccess < oldest.lastAccess {
			oldest = line
		}

		if line.tag != tag {
			continue
		}

		found = true
		line.lastAccess = cs.accessCounter
		line.accessedBytes += size
		break
	}

	if found {
		cs.hits++
		hits = 1
		return
	}

	if available != nil {
		available.lastAccess = cs.accessCounter
		available.tag = tag
		available.valid = true
		available.loadedBytes += 1 << cs.cfg.BlockOffsetBits
		available.accessedBytes += size
		cs.misses++
		misses = 1
		return
	}

	oldest.lastAccess = cs.accessCounter
	oldest.tag = tag
	oldest.valid = true
	oldest.loadedBytes += 1 << cs.cfg.BlockOffsetBits
	oldest.accessedBytes += size
	cs.misses++
	cs.evictions++
	misses = 1
	evictions = 1
	return
}

func (cs *cacheSet) efficiencyStats() iter.Seq2[int, float64] {
	return func(yield func(int, float64) bool) {
		for i := range len(cs.lines) {
			line := &cs.lines[i]
			if line.loadedBytes == 0 {
				if !yield(i, 0) {
					break
				}
				continue
			}

			if !yield(i, float64(line.accessedBytes)/float64(line.loadedBytes)) {
				break
			}
		}
	}
}

type Config struct {
	SetIndexBits    int
	BlockOffsetBits int
	LinesPerSet     int
}

type Simulator struct {
	cfg *Config

	sets []cacheSet
}

const maxBits = 30

func New(cfg *Config) (*Simulator, error) {
	if cfg == nil {
		return nil, errors.New("cache.New: nil config")
	}
	if cfg.SetIndexBits < 0 || cfg.BlockOffsetBits < 0 {
		return nil, errors.New("cache.New: set index and block offset bits must be non-negative")
	}
	if cfg.LinesPerSet <= 0 {
		return nil, errors.New("cache.New: lines per set must be positive")
	}
	if cfg.SetIndexBits > maxBits {
		return nil, fmt.Errorf("cache.New: set index bits %d exceed maximum %d", cfg.SetIndexBits, maxBits)
	}
	if cfg.BlockOffsetBits > maxBits {
		return nil, fmt.Errorf("cache.New: block offset bits %d exceed maximum %d", cfg.BlockOffsetBits, maxBits)
	}

	sim := &Simulator{
		cfg: cfg,
	}

	numSets := 1 << cfg.SetIndexBits
	sim.sets = make([]cacheSet, numSets)
	for i := range numSets {
		sim.sets[i].init(cfg.LinesPerSet, cfg)
	}

	return sim, nil
}

func (sim *Simulator) Access(address uint64, size int) (hits int, misses int, evictions int) {
	if size <= 0 {
		return 0, 0, 0
	}

	record := func(h, m, e int) {
		hits += h
		misses += m
		evictions += e
	}

	blkSize := 1 << sim.cfg.BlockOffsetBits
	{
		setIdx := int((address >> uint64(sim.cfg.BlockOffsetBits)) & ((1 << uint64(sim.cfg.SetIndexBits)) - 1))
		tag := address >> uint64(sim.cfg.BlockOffsetBits+sim.cfg.SetIndexBits)
		blkOffset := int(address & ((1 << uint64(sim.cfg.BlockOffsetBits)) - 1))
		record(sim.sets[setIdx].access(tag, min(size, blkSize-blkOffset)))

		size -= (blkSize - blkOffset)
		if size <= 0 {
			return
		}
	}

	for size > blkSize {
		address += uint64(blkSize)
		setIdx := int((address >> uint64(sim.cfg.BlockOffsetBits)) & ((1 << sim.cfg.SetIndexBits) - 1))
		tag := address >> uint64(sim.cfg.BlockOffsetBits+sim.cfg.SetIndexBits)
		record(sim.sets[setIdx].access(tag, blkSize))
		size -= blkSize
	}

	{
		address += uint64(blkSize)
		setIdx := int((address >> uint64(sim.cfg.BlockOffsetBits)) & ((1 << sim.cfg.SetIndexBits) - 1))
		tag := address >> uint64(sim.cfg.BlockOffsetBits+sim.cfg.SetIndexBits)
		record(sim.sets[setIdx].access(tag, size))
	}

	return
}

func (sim *Simulator) Size() int {
	return (1 << sim.cfg.SetIndexBits) * (1 << sim.cfg.BlockOffsetBits) * sim.cfg.LinesPerSet
}

// NumSets returns the number of sets in the cache.
func (sim *Simulator) NumSets() int {
	return 1 << sim.cfg.SetIndexBits
}

// LinesPerSet returns the associativity (number of lines per set) of the cache.
func (sim *Simulator) LinesPerSet() int {
	return sim.cfg.LinesPerSet
}

type HMEStats struct {
	Hits      int
	Misses    int
	Evictions int
}

func (sim *Simulator) HMEStats() iter.Seq2[int, HMEStats] {
	return func(yield func(int, HMEStats) bool) {
		for i := range len(sim.sets) {
			set := &sim.sets[i]

			if !yield(i, HMEStats{
				Hits:      set.hits,
				Misses:    set.misses,
				Evictions: set.evictions,
			}) {
				break
			}
		}
	}
}

func (sim *Simulator) EfficiencyStats() iter.Seq2[int, float64] {
	return func(yield func(int, float64) bool) {
		for i := range len(sim.sets) {
			for j, e := range sim.sets[i].efficiencyStats() {
				if !yield(i*sim.cfg.LinesPerSet+j, e) {
					break
				}
			}
		}
	}
}
