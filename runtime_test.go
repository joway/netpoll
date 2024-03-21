package netpoll

import "testing"

func TestGetG(t *testing.T) {
	gp := getg()
	var gp2 uintptr
	done := make(chan struct{})
	go func() {
		gp2 = getg()
		done <- struct{}{}
	}()
	<-done
	t.Logf("gp=0x%X, gp2=0x%X", gp, gp2)
}

func TestWaitG(t *testing.T) {
	gpch := make(chan uintptr, 1)
	done := make(chan struct{})
	go func() {
		gp := getg()
		gpch <- gp
		park()
		done <- struct{}{}
	}()
	gp := <-gpch
	t.Logf("gp=0x%X parking", gp)
}
