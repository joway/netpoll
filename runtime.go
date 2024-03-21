// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netpoll

import (
	"sync/atomic"
	"unsafe"
	_ "unsafe"
)

type waitReason uint8
type traceBlockReason uint8

const (
	waitReasonChanReceive waitReason = 14
	waitReasonChanSend    waitReason = 15

	traceBlockChanSend traceBlockReason = 22
	traceBlockChanRecv traceBlockReason = 23
)

// -- asm ---
func getg() uintptr

// --- link ---

//go:linkname runtime_pollOpen internal/poll.runtime_pollOpen
func runtime_pollOpen(fd uintptr) (pd uintptr, errno int)

//go:linkname runtime_pollWait internal/poll.runtime_pollWait
func runtime_pollWait(pd uintptr, mode int) (errno int)

//go:linkname runtime_pollReset internal/poll.runtime_pollReset
func runtime_pollReset(pd uintptr, mode int) (errno int)

//go:linkname runtime_pollClose internal/poll.runtime_pollClose
func runtime_pollClose(pd uintptr)

//go:linkname gopark runtime.gopark
func gopark(unlockf func(gp uintptr, _ unsafe.Pointer) bool, lock unsafe.Pointer, reason waitReason, traceReason traceBlockReason, traceskip int)

// goready: runqput and wakep
//
//go:linkname goready runtime.goready
func goready(gp uintptr, traceskip int)

//go:linkname wakep runtime.wakep
func wakep()

var golist = make([]atomic.Value, 1)
var taskRing = make([]atomic.Value, 1)

func notify(gp uintptr, _ unsafe.Pointer) bool {
	gv := golist[0]
	gv.Store(gp)
	return true
}

func park() {
	gopark(notify, nil, waitReasonChanReceive, traceBlockChanRecv, 1)
}

func ready(gp uintptr) {
	goready(gp, 1)
}
