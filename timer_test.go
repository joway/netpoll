package netpoll

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	maxInterval := time.Second
	tmr := newTimer(time.Millisecond*10, maxInterval)

	flag := -1
	select {
	case <-tmr.After(time.Millisecond):
		flag = 0
	}
	Assert(t, flag == 0)

	flag = -1
	select {
	case <-tmr.After(time.Millisecond):
		flag = 0
	default:
		flag = 1
	}
	Assert(t, flag == 1)

	flag = -1
	select {
	case <-tmr.After(time.Millisecond * 10):
		flag = 0
	case <-time.After(time.Millisecond):
		flag = 1
	}
	Assert(t, flag == 1)

	flag = -1
	select {
	case <-tmr.After(time.Millisecond * 10):
		flag = 0
	case <-time.After(time.Millisecond * 100):
		flag = 1
	}
	Assert(t, flag == 0)

	beg := time.Now()
	var end time.Time
	select {
	case <-tmr.After(time.Minute):
		end = time.Now()
	}
	Assert(t, (end.Sub(beg).Milliseconds()-maxInterval.Milliseconds()) < 20)
}

func TestMultiTimer(t *testing.T) {
	tmr := newTimer(time.Millisecond*10, time.Second)

	var wg sync.WaitGroup
	wg.Add(2)
	var flag int32
	go func() {
		defer wg.Done()
		<-tmr.After(time.Millisecond * 50)
		atomic.AddInt32(&flag, 1)
	}()
	go func() {
		defer wg.Done()
		<-tmr.After(time.Millisecond * 200)
		atomic.AddInt32(&flag, 1)
	}()

	Assert(t, atomic.LoadInt32(&flag) == 0)
	time.Sleep(time.Millisecond * 10)
	Assert(t, atomic.LoadInt32(&flag) == 0)
	time.Sleep(time.Millisecond * 100)
	Assert(t, atomic.LoadInt32(&flag) == 1)
	time.Sleep(time.Millisecond * 150)
	Assert(t, flag == 2)
	wg.Wait()

	flag = 0
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-tmr.After(time.Millisecond * 10)
		flag = 1
	}()
	tmr.Stop()
	tmr.Stop() // check double stop safety
	Assert(t, flag == 0)
	time.Sleep(time.Millisecond * 50)
	Assert(t, flag == 0)
}

func benchmarkTimer(pb *testing.PB, waitFunc func()) {
	for pb.Next() {
		var wg sync.WaitGroup
		for i := 0; i < 1024; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				waitFunc()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkTimer(b *testing.B) {
	tmr := newTimer(time.Millisecond*10, time.Second*60)
	defer tmr.Stop()

	b.RunParallel(func(pb *testing.PB) {
		benchmarkTimer(pb, func() {
			<-tmr.After(time.Millisecond * 100)
		})
	})
}

func BenchmarkNativeTimer(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		benchmarkTimer(pb, func() {
			<-time.After(time.Millisecond * 100)
		})
	})
}
