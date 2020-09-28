package astikit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Stater is an object that can compute and handle stats
type Stater struct {
	cancel  context.CancelFunc
	ctx     context.Context
	h       StatsHandleFunc
	m       *sync.Mutex // Locks ss
	period  time.Duration
	running uint32
	ss      []stat
}

// StatsHandleFunc is a method that can handle stats
type StatsHandleFunc func(stats []Stat)

// StatMetadata represents a stat metadata
type StatMetadata struct {
	Description string
	Label       string
	Name        string
	Unit        string
}

// StatHandler represents a stat handler
type StatHandler interface {
	Start()
	Stop()
	Value(delta time.Duration) interface{}
}

// Stat represents a stat
type Stat struct {
	StatMetadata
	Value interface{}
}

type stat struct {
	h StatHandler
	m StatMetadata
}

// StaterOptions represents stater options
type StaterOptions struct {
	HandleFunc StatsHandleFunc
	Period     time.Duration
}

// NewStater creates a new stater
func NewStater(o StaterOptions) *Stater {
	return &Stater{
		h:      o.HandleFunc,
		m:      &sync.Mutex{},
		period: o.Period,
	}
}

// Start starts the stater
func (s *Stater) Start(ctx context.Context) {
	// Check context
	if ctx.Err() != nil {
		return
	}

	// Make sure to start only once
	if atomic.CompareAndSwapUint32(&s.running, 0, 1) {
		// Update status
		defer atomic.StoreUint32(&s.running, 0)

		// Reset context
		s.ctx, s.cancel = context.WithCancel(ctx)

		// Start stats
		s.m.Lock()
		for _, v := range s.ss {
			v.h.Start()
		}
		s.m.Unlock()

		// Create ticker
		t := time.NewTicker(s.period)
		defer t.Stop()

		// Loop
		lastStatAt := now()
		for {
			select {
			case <-t.C:
				// Get delta
				n := now()
				delta := n.Sub(lastStatAt)
				lastStatAt = n

				// Loop through stats
				var stats []Stat
				s.m.Lock()
				for _, v := range s.ss {
					stats = append(stats, Stat{
						StatMetadata: v.m,
						Value:        v.h.Value(delta),
					})
				}
				s.m.Unlock()

				// Handle stats
				go s.h(stats)
			case <-s.ctx.Done():
				// Stop stats
				s.m.Lock()
				for _, v := range s.ss {
					v.h.Stop()
				}
				s.m.Unlock()
				return
			}
		}
	}
}

// AddStat adds a stat
func (s *Stater) AddStat(m StatMetadata, h StatHandler) {
	s.m.Lock()
	defer s.m.Unlock()
	s.ss = append(s.ss, stat{
		h: h,
		m: m,
	})
}

// Stop stops the stater
func (s *Stater) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// StatsMetadata returns the stats metadata
func (s *Stater) StatsMetadata() (ms []StatMetadata) {
	s.m.Lock()
	defer s.m.Unlock()
	ms = []StatMetadata{}
	for _, v := range s.ss {
		ms = append(ms, v.m)
	}
	return
}

type durationStat struct {
	d         time.Duration
	fn        func(d, delta time.Duration) interface{}
	isStarted bool
	m         *sync.Mutex // Locks isStarted
	startedAt time.Time
}

func newDurationStat(fn func(d, delta time.Duration) interface{}) *durationStat {
	return &durationStat{
		fn: fn,
		m:  &sync.Mutex{},
	}
}

func (s *durationStat) Begin() {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.startedAt = now()
}

func (s *durationStat) End() {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.d += now().Sub(s.startedAt)
	s.startedAt = time.Time{}
}

func (s *durationStat) Value(delta time.Duration) (o interface{}) {
	// Lock
	s.m.Lock()
	defer s.m.Unlock()

	// Get current values
	n := now()
	d := s.d

	// Recording is still in process
	if !s.startedAt.IsZero() {
		d += n.Sub(s.startedAt)
		s.startedAt = n
	}

	// Compute stat
	o = s.fn(d, delta)
	s.d = 0
	return
}

func (s *durationStat) Start() {
	s.m.Lock()
	defer s.m.Unlock()
	s.d = 0
	s.isStarted = true
}

func (s *durationStat) Stop() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isStarted = false
}

// DurationPercentageStat is an object capable of computing the percentage of time some work is taking per second
type DurationPercentageStat struct {
	*durationStat
}

// NewDurationPercentageStat creates a new duration percentage stat
func NewDurationPercentageStat() *DurationPercentageStat {
	return &DurationPercentageStat{durationStat: newDurationStat(func(d, delta time.Duration) interface{} {
		if delta == 0 {
			return 0
		}
		return float64(d) / float64(delta) * 100
	})}
}

type counterStat struct {
	c         float64
	fn        func(c, t float64, delta time.Duration) interface{}
	isStarted bool
	m         *sync.Mutex // Locks isStarted
	t         float64
}

func newCounterStat(fn func(c, t float64, delta time.Duration) interface{}) *counterStat {
	return &counterStat{
		fn: fn,
		m:  &sync.Mutex{},
	}
}

func (s *counterStat) Add(delta float64) {
	s.m.Lock()
	defer s.m.Unlock()
	if !s.isStarted {
		return
	}
	s.c += delta
	s.t++
}

func (s *counterStat) Start() {
	s.m.Lock()
	defer s.m.Unlock()
	s.c = 0
	s.isStarted = true
	s.t = 0
}

func (s *counterStat) Stop() {
	s.m.Lock()
	defer s.m.Unlock()
	s.isStarted = true
}

func (s *counterStat) Value(delta time.Duration) interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	c := s.c
	t := s.t
	s.c = 0
	s.t = 0
	return s.fn(c, t, delta)
}

// CounterAvgStat is an object capable of computing the average value of a counter
type CounterAvgStat struct {
	*counterStat
}

// NewCounterAvgStat creates a new counter avg stat
func NewCounterAvgStat() *CounterAvgStat {
	return &CounterAvgStat{counterStat: newCounterStat(func(c, t float64, delta time.Duration) interface{} {
		if t == 0 {
			return 0
		}
		return c / t
	})}
}

// CounterRateStat is an object capable of computing the average value of a counter per second
type CounterRateStat struct {
	*counterStat
}

// NewCounterRateStat creates a new counter rate stat
func NewCounterRateStat() *CounterRateStat {
	return &CounterRateStat{counterStat: newCounterStat(func(c, t float64, delta time.Duration) interface{} {
		if delta.Seconds() == 0 {
			return 0
		}
		return c / delta.Seconds()
	})}
}
