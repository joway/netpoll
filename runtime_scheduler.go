package netpoll

import (
	"unsafe"

	"github.com/bytedance/gopkg/collection/lscq"
)

type Task func()

func newRunQ(size int) *runq {
	rq := new(runq)
	rq.gqsize = size
	rq.gq = lscq.NewPointer()
	rq.tq = lscq.NewPointer()
	for i := 0; i < rq.gqsize; i++ {
		go reusedWorker(rq)
	}
	return rq
}

type runq struct {
	gqsize int
	gq     *lscq.PointerQueue // goroutine queue
	tq     *lscq.PointerQueue // task queue
}

func reusedWorker(rq *runq) {
	for {
		gopark(parkg, unsafe.Pointer(rq), waitReasonChanReceive, traceBlockChanRecv, 1)
		task, ok := rq.tq.Dequeue()
		if task == nil || !ok {
			continue
		}
		tp := (*Task)(task)
		if tp == nil || *tp == nil {
			continue
		}
		(*tp)()
		wakep()
	}
}

func (rq *runq) Go(task Task) {
	if task == nil {
		return
	}

	// get worker
	gp, ok := rq.gq.Dequeue()
	if !ok {
		// no idle worker
		go task()
		return
	}
	rq.tq.Enqueue(unsafe.Pointer(&task))
	goready(uintptr(gp), 1)
	wakep()
}

func parkg(gp uintptr, q unsafe.Pointer) bool {
	rq := (*runq)(q)
	rq.gq.Enqueue(unsafe.Pointer(gp))
	return true
}
