package netpoll

import (
	"sync"
	"time"
)

var globalTimer = newTimer(time.Millisecond*10, time.Second*30)

type Timer struct {
	sync.RWMutex
	precision time.Duration
	ticker    *time.Ticker
	slots     []chan struct{}
	cap       int
	pos       int
	quit      chan struct{}
}

func newTimer(precision, maxInterval time.Duration) *Timer {
	t := &Timer{
		quit: make(chan struct{}),
	}
	t.Reset(precision, maxInterval)
	go t.run()
	return t
}

func (t *Timer) Reset(precision, maxInterval time.Duration) {
	t.Lock()
	defer t.Unlock()

	t.precision = precision
	t.ticker = time.NewTicker(precision)
	t.pos = 0
	t.cap = int(maxInterval / precision)
	t.slots = make([]chan struct{}, t.cap)
	for i := 0; i < t.cap; i++ {
		t.slots[i] = make(chan struct{})
	}
}

func (t *Timer) After(interval time.Duration) <-chan struct{} {
	target := int(interval / t.precision)
	if target > t.cap {
		target = t.cap
	}
	t.RLock()
	defer t.RUnlock()

	return t.slots[(target+t.pos)%t.cap]
}

func (t *Timer) Stop() {
	select {
	case <-t.quit:
		return
	default:
		close(t.quit)
	}
}

func (t *Timer) check() {
	t.Lock()
	signal := t.slots[t.pos]
	t.slots[t.pos] = make(chan struct{})
	t.pos = (t.pos + 1) % t.cap
	t.Unlock()

	close(signal)
}

func (t *Timer) run() {
	for {
		select {
		case <-t.quit:
			t.ticker.Stop()
			return
		case <-t.ticker.C:
			t.check()
		}
	}
}
