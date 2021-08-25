package tls

import (
	"crypto/tls"
	"net"

	"github.com/cloudwego/netpoll"
)

// CreateListener return a new Listener.
func CreateListener(network, addr string, config *tls.Config) (l netpoll.Listener, err error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		ln, err := net.Listen(network, addr)
		if err != nil {
			return nil, err
		}
		return newListener(ln, config)
	default:
		return nil, netpoll.Exception(netpoll.ErrUnsupported, network)
	}
}

type listener struct {
	fd     int
	ln     net.Listener
	config *tls.Config
}

func newListener(ln net.Listener, config *tls.Config) (netpoll.Listener, error) {
	l := &listener{
		ln:     ln,
		config: config,
	}
	f, err := netpoll.ParseListener(ln)
	if err != nil {
		return nil, err
	}
	l.fd = int(f.Fd())
	return l, nil
}

// Accept implements Listener.
func (ln *listener) Accept() (net.Conn, error) {
	conn, err := ln.ln.Accept()
	if err != nil {
		return nil, err
	}
	return tls.Server(conn, ln.config), nil
}

// Close implements Listener.
func (ln *listener) Close() error {
	return ln.ln.Close()
}

// Addr implements Listener.
func (ln *listener) Addr() net.Addr {
	return ln.ln.Addr()
}

// Fd implements Listener.
func (ln *listener) Fd() (fd int) {
	return ln.fd
}
