package netpoll

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"testing"
	"unsafe"

	"github.com/bytedance/gopkg/collection/lscq"
)

func TestRuntimeScheduler(t *testing.T) {
	rq := newRunQ(100)
	Equal(t, rq.gqsize, 100)
	var wg sync.WaitGroup
	for n := 0; n < 10; n++ {
		for i := 0; i < rq.gqsize*2; i++ {
			wg.Add(1)
			rq.Go(func() {
				defer wg.Done()
				x := factorial(1000000000000000000)
				_ = x
			})
		}
		wg.Wait()
	}
}

func TestPointerQueue(t *testing.T) {
	pq := lscq.NewPointer()
	for c := 0; c < 100; c++ {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			var task Task = func() {
				defer wg.Done()
				x := factorial(100000000000000000)
				_ = x
			}
			pq.Enqueue(unsafe.Pointer(&task))
		}
		for {
			tp, ok := pq.Dequeue()
			if !ok {
				break
			}
			(*(*Task)(tp))()
		}
		wg.Wait()
	}
}

//go:noinline
func factorial(n int) int {
	if n <= 1 {
		return n
	}
	return n * (n - 1)
}

//go:noinline
func cpuWork(n int) (sum int) {
	for i := 0; i < n; i++ {
		sum += factorial(n - i)
	}
	return sum
}

func BenchmarkScheduler(b *testing.B) {
	fn := fmt.Sprintf("scheduler.trace")
	f, _ := os.Create(fn)
	defer f.Close()
	_ = trace.Start(f)
	defer trace.Stop()

	concurrency := 100
	rq := newRunQ(concurrency)
	var benchcases = []struct {
		Name   string
		GoFunc func(t func())
	}{
		{Name: "Native", GoFunc: func(t func()) {
			go t()
		}},
		{Name: "NetpollScheduler", GoFunc: func(t func()) {
			rq.Go(t)
		}},
	}
	for _, bc := range benchcases {
		b.Run(bc.Name, func(b *testing.B) {
			old := runtime.GOMAXPROCS(8)
			defer runtime.GOMAXPROCS(old)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				for id := 0; id < concurrency; id++ {
					wg.Add(1)
					bc.GoFunc(func() {
						defer wg.Done()
						_ = cpuWork(10000)
					})
				}
				wg.Wait()
			}
		})
	}
}
