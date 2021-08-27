package tls

import (
	"crypto/tls"
	"net"
	"strings"
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
	conn, err := netpoll.DialConnection(network, address, timeout)
	if err != nil {
		return nil, err
	}

	colonPos := strings.LastIndex(address, ":")
	if colonPos == -1 {
		colonPos = len(address)
	}
	hostname := address[:colonPos]

	// If no ServerName is set, infer the ServerName
	// from the hostname we're connecting to.
	if config.ServerName == "" {
		// Make a copy to avoid polluting argument or default.
		c := config.Clone()
		c.ServerName = hostname
		config = c
	}

	client := tls.Client(conn, config)
	if err = client.Handshake(); err != nil {
		return nil, err
	}

	return GetConnection(client)
}
