package tls

import (
	"testing"
	"time"
)

func TestListener(t *testing.T) {
	go testTLSServer(t)

	time.Sleep(time.Millisecond * 100)

	testNativeTLSClient(t)
}
