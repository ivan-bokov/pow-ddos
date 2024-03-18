package ddos

import (
	"sync"
	"time"
)

type Guard struct {
	mu         sync.RWMutex
	rate       uint64
	lastTime   time.Time
	window     time.Duration
	difficulty uint8
	stepRate   uint64
}

func New(window time.Duration, difficulty uint8, stepRate uint64) *Guard {
	return &Guard{
		mu:         sync.RWMutex{},
		window:     window,
		lastTime:   time.Now(),
		difficulty: difficulty,
		stepRate:   stepRate,
	}
}

func (g *Guard) Take() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if time.Since(g.lastTime) >= g.window {
		g.calculate(int(time.Since(g.lastTime) / g.window))
		g.rate = 0
		g.lastTime = time.Now()
	}
	g.rate++
}

func (g *Guard) calculate(delta int) {
	switch {
	case uint64(g.difficulty) >= g.rate/g.stepRate:
		g.difficulty -= uint8(delta)
	case uint64(g.difficulty) >= g.rate/g.stepRate:
		g.difficulty += uint8(delta)
	}
	if g.difficulty < 1 {
		g.difficulty = 1
	}
}

func (g *Guard) Rate() uint64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.rate
}

func (g *Guard) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.rate = 0
	g.lastTime = time.Now()
}

func (g *Guard) Difficulty() uint8 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.difficulty
}
