package ddos

import (
	"testing"
	"time"
)

func TestTakeWithinWindow(t *testing.T) {
	g := &Guard{
		window:   time.Second,
		lastTime: time.Now(),
		rate:     0,
	}
	g.Take()
	if g.rate != 1 {
		t.Errorf("Expected rate to be 1, got %d", g.rate)
	}
}

func TestTakeOutsideWindow(t *testing.T) {
	g := &Guard{
		window:   time.Second,
		lastTime: time.Now().Add(-2 * time.Second),
		rate:     5,
	}
	g.Take()
	if g.rate != 1 {
		t.Errorf("Expected rate to be 1, got %d", g.rate)
	}
	if time.Since(g.lastTime) >= g.window {
		t.Errorf("lastTime should be within the window")
	}
}

func TestReset(t *testing.T) {
	g := &Guard{
		window:   time.Second,
		lastTime: time.Now().Add(-2 * time.Second),
		rate:     5,
	}
	g.Take()
	g.Reset()
	if g.rate != 0 {
		t.Errorf("Expected rate to be 0, got %d", g.rate)
	}
}

func TestRate(t *testing.T) {
	g := &Guard{
		window:   time.Second,
		lastTime: time.Now().Add(-2 * time.Second),
		rate:     5,
	}
	g.Take()
	rate := g.Rate()
	if rate != 1 {
		t.Errorf("Expected rate to be 5, got %d", rate)
	}
}
