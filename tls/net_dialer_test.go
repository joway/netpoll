package tls

import (
	"testing"
	"time"
)

func TestDialer(t *testing.T) {
	go testNativeTLSServer(t)

	time.Sleep(time.Millisecond * 100)

	testTLSClient(t)
}
