# TLS Usage

Netpoll use the `tls.Conn(crypto/tls)` as the underlying implementation for TLS connection.

We provide `netpoll/tls` package to help you use the same `netpoll.Connection` interface for both TLS and non-TLS use cases. 

## Server

### Create TLS Config

```go
import (
    "crypto/tls"
)

// Load server cert
serverCert, _ := tls.X509KeyPair([]byte(serverPEM), []byte(serverKEY))

// Optional
clientCertPool := x509.NewCertPool()
clientCertPool.AppendCertsFromPEM([]byte(clientRootREM))

config := &tls.Config{
    Certificates: []tls.Certificate{serverCert},
    
    // Ignore follows if you don't need mTLS
    ClientAuth:   tls.RequireAndVerifyClientCert,
    ClientCAs:    clientCertPool,
}
```

### Create TLS Listener

```go
import (
    "crypto/tls"
    ntls "github.com/cloudwego/netpoll/tls"
)

ln, err := ntls.CreateListener("tcp", ":443", config)

for {
    conn, err := ln.Accept()
    
    // Get netpoll.Connection for nocopy interface
    tlsConn, err := ntls.GetConnection(conn)
    r, w := tlsConn.Reader(), tlsConn.Writer()

    msg := "hello"
    go func() {
        defer tlsConn.Close()
        for {
        	// reading
            resp, err := r.ReadString(len(msg))
            if errors.Is(err, io.EOF) {
                return
            }
            r.Release()
            
            // writing
            _, err = w.WriteString(msg)
            w.Flush()
        }
    }()
}
```

## Client

### Create TLS Config

```go
import (
    "crypto/tls"
)

// Optional
clientCert, _ := tls.X509KeyPair([]byte(clientPEM), []byte(clientKEY))
clientCertPool := x509.NewCertPool()
clientCertPool.AppendCertsFromPEM([]byte(clientRootREM))

config := &tls.Config{
    ServerName:         "p.s.m",
    
    // Ignore follows if you don't need mTLS
    RootCAs:            clientCertPool,
    Certificates:       []tls.Certificate{clientCert},
    InsecureSkipVerify: true,
}
```

### Create TLS Connection

```go
import (
    ntls "github.com/cloudwego/netpoll/tls"
)

conn, err := ntls.DialConnection("tcp", ":443", time.Second, config)
r, w := conn.Reader(), conn.Writer()

msg := "hello"
for i := 0; i < 1024; i++ {
	// writing
    _, err := w.WriteString(msg)
    err = w.Flush()
    
    // reading
    resp, err := r.ReadString(len(msg))
    err = r.Release()
}
```
