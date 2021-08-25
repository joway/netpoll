package tls

import (
	"crypto/tls"
	"net"
	"sync/atomic"
	"time"

	"github.com/cloudwego/netpoll"
)

var _ netpoll.Connection = (*connection)(nil)

type connection struct {
	*tls.Conn
	r           netpoll.Reader
	w           netpoll.Writer
	readTimeout time.Duration
	idleTimeout time.Duration
	closed      int32
}

func GetConnection(conn net.Conn) (*connection, error) {
	c, ok := conn.(*tls.Conn)
	if !ok {
		return nil, netpoll.Exception(netpoll.ErrUnsupported, "GetConnection")
	}
	return &connection{
		Conn: c,
		r:    netpoll.NewReader(conn),
		w:    netpoll.NewWriter(conn),
	}, nil
}

func (c *connection) Reader() netpoll.Reader {
	return c.r
}

func (c *connection) Writer() netpoll.Writer {
	return c.w
}

func (c *connection) IsActive() bool {
	return atomic.LoadInt32(&c.closed) == 0
}

func (c *connection) SetReadTimeout(timeout time.Duration) error {
	return netpoll.Exception(netpoll.ErrUnsupported, "SetReadTimeout")
}

func (c *connection) SetIdleTimeout(timeout time.Duration) error {
	return netpoll.Exception(netpoll.ErrUnsupported, "SetIdleTimeout")
}

// SetOnRequest is not implement in TLS connection
func (c *connection) SetOnRequest(on netpoll.OnRequest) error {
	return netpoll.Exception(netpoll.ErrUnsupported, "SetOnRequest")
}

func (c *connection) AddCloseCallback(callback netpoll.CloseCallback) error {
	return netpoll.Exception(netpoll.ErrUnsupported, "AddCloseCallback")
}

func (c *connection) Close() error {
	atomic.StoreInt32(&c.closed, 1)
	return c.Conn.Close()
}

func (c *connection) CloseWrite() error {
	atomic.StoreInt32(&c.closed, 1)
	return c.Conn.CloseWrite()
}
