package tls

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/cloudwego/netpoll"
)

var defaultDialer = NewDialer()

func DialConnection(network, address string, timeout time.Duration, config *tls.Config) (connection netpoll.Connection, err error) {
	return defaultDialer.DialConnection(network, address, timeout, config)
}

type Dialer interface {
	DialConnection(network, address string, timeout time.Duration, config *tls.Config) (connection netpoll.Connection, err error)
	DialTimeout(network, address string, timeout time.Duration, config *tls.Config) (conn net.Conn, err error)
}

func NewDialer() Dialer {
	return &dialer{}
}

type dialer struct{}

// DialTimeout implements Dialer.
func (d *dialer) DialTimeout(network, address string, timeout time.Duration, config *tls.Config) (net.Conn, error) {
	return d.DialConnection(network, address, timeout, config)
}

func (d *dialer) DialConnection(
	network, address string, timeout time.Duration, config *tls.Config,
) (connection netpoll.Connection, err error) {
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(dialer, network, address, config)
	if err != nil {
		return nil, err
	}
	return GetConnection(conn)
}
