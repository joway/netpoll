package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/netpoll"
)

func TestDialerAndListener(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		testTLSServer(t)
	}()

	time.Sleep(time.Millisecond * 100)

	testTLSClient(t)
	wg.Wait()
}

func TestDialer(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		testNativeTLSServer(t)
	}()

	time.Sleep(time.Millisecond * 100)

	testTLSClient(t)
	wg.Wait()
}

func TestListener(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		testTLSServer(t)
	}()

	time.Sleep(time.Millisecond * 100)

	testNativeTLSClient(t)
	wg.Wait()
}

var (
	clientRootREM = clientPEM
	clientPEM     = `-----BEGIN CERTIFICATE-----
MIICvjCCAaYCCQCGwzRh/G+VoTANBgkqhkiG9w0BAQsFADAhMQswCQYDVQQGEwJi
ZDESMBAGA1UECAwJYnl0ZWRhbmNlMB4XDTIxMDgyMDA5MDA0MFoXDTMxMDgxODA5
MDA0MFowITELMAkGA1UEBhMCYmQxEjAQBgNVBAgMCWJ5dGVkYW5jZTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBANUEoJl5kYg+EfdromDd7f9Ae+3hfGMr
LkPVxpxB+9a3rYACyhz/0x3uj4/QF1Ldc3i9KVOrHrFwc8mGuBcm2P+1Ij/903go
t6T2iwRi313kEzDqHUfG3ACaMwoI9loIu9F8nd24H4fqHNl3kOekNSd4GMheo5Cf
c05q+5fnXxDBq6TC+7FftvFtf6dVpTaDblR0lScoq+ejEoho1Yze/Kn3NCW9sr+z
8Pui8ycSet4xwZLZNXVgkXyNIpoEU2vc7++3ndjef82HabHD6OydsWjJllvSl/cd
jiGYZpjgwSjMbvWefANoi3SgB/TaJqgR+5HnWodl7bocblysacGVg8sCAwEAATAN
BgkqhkiG9w0BAQsFAAOCAQEA0bvIXb1Sq+StxPL9MI1qVPUoxzA4x0ugHgt+8qx8
L7MiFMcVuLR022PcWsmrRPrShSeHtkA/e4fXd7mMQdbm/udi9u7UhS+nSrKUzLO4
yxHot9kEKwynghNhTr4FYLfBwWI9CCS51Ju3FKe1o4TdYHDSPUR99h0xWVhIsUq2
kDIoIvWZb/ixsT+oMGVsCVZDchlDsf2UDq6XQdGIZWlaCMW9RS8rjDvCU4ibZ635
JeoCU4fPuPkNxD6Qixbv65By1J2IPst7pQBBg4tfZ9LnYIsx/FLViK99TiwY/SZB
+MBqU9vo37G5SdtkrLBlKGNJvozghg+DK1d/sF0YPZBQuQ==
-----END CERTIFICATE-----`
	clientKEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1QSgmXmRiD4R92uiYN3t/0B77eF8YysuQ9XGnEH71retgALK
HP/THe6Pj9AXUt1zeL0pU6sesXBzyYa4FybY/7UiP/3TeCi3pPaLBGLfXeQTMOod
R8bcAJozCgj2Wgi70Xyd3bgfh+oc2XeQ56Q1J3gYyF6jkJ9zTmr7l+dfEMGrpML7
sV+28W1/p1WlNoNuVHSVJyir56MSiGjVjN78qfc0Jb2yv7Pw+6LzJxJ63jHBktk1
dWCRfI0imgRTa9zv77ed2N5/zYdpscPo7J2xaMmWW9KX9x2OIZhmmODBKMxu9Z58
A2iLdKAH9NomqBH7kedah2XtuhxuXKxpwZWDywIDAQABAoIBAQCjI2MXhvoE3HyH
/1+hTfXqaWk/+lN342cQPpVXsFT90TfH9NrzMJ1vq/c4jYZ/SOgZTM1srX3ZKKUU
bt9no7HGy+MKHQuqV4ylgsdeuZYVYwmriXCZOtvcfjuecNSyiUPGIOkKF+vV/F7R
0XchYCnxK1HXils79FGredVrNaAhxLsoVtxk5FGhuc7KzWTAwM/YvCwMhBbn8dEI
9Xd6Nkt4dzhuU8tfUTAPVa9tdZD93SXN1ZZB4Yx7atjfBTmJwJjPP+m5/nXYqQbB
EbfjCrzrPq6qBInFSBYUtFG5bONjU103UrOp4wja6HVLniPaVaPcgz5wakToSIN/
PaU6xvRRAoGBAPBaQ8h2Bui5zsyXdcaqkyCG4V9XGLKxGd8x+ptPza4Qc6N/E6UH
qgf0QQXw4x0PlHi/WIAOwaRvpKZ92fboTUgJbSYcKRRD4F4qe3Du4oHkq/ww5rXg
xdw0JZpLjrep768VVW4eOxSq94WYiM3aObCWsxzL5dtDf//ARw+h58rDAoGBAOLi
zaRvFobZwgsb7mChI36nKbejRNp+ESZZBjCmte2EhZM+kRKinmjtPRObjDA8pFsB
J85e38WnekQrEZXYBRj0CxVvdI7ARP2X5e1N4fMhk2DyTf4e4a1swqeW/omZWb1w
jctDoWe4mbtjm/mpZ3/EBrGbkUOqSY9c0oVmQ4JZAoGBAK4riwEB3mHY+62wd+1e
AD3K4BhZohEjWy8tJYBUpSRk0ZeB57doRWN8MX5foASYHKwfZL9vcg15xaPMgR12
G0J/ajP6ph4ETduPB2LACS29io+21AiqeFbpBvK3nBUltQV3/S9OAtwoRDPwY/pg
D/wSSHsumkN4t5GaQSRn2/NHAoGAdpUN4AyDnJWBmqbNj0mJMLAT2LwHx56uPfm3
h4QKgAqMeenwjunZm4OrMW1R9wAq8rmG4ZCqqjafa7OK7GNMPr+Gb3yiUd3h8R0L
+lyDZLy+t6PM6a2gTDEVB9yeSrKQubdzFLLTUE+mYc9s/S7yPk/pI7joUpJVAg4E
pd5OGHkCgYEA52Bckm0PAGX/ZtHLZVYRviVIAGdhDcXXwG0MSAI9kwm01a5MnReB
Mo34DNlSDEwOVx6AlD9p0NqdRZU3hxzdrhgrGse7egi/83jLWODOaF5joRohxJxu
kJ3MMtBCAUihO0ecYvscqdh2kR/Lud2NomWgyGTMPSlBcWcYbBSUzlY=
-----END RSA PRIVATE KEY-----`
	serverPEM = `-----BEGIN CERTIFICATE-----
MIIDQDCCAigCCQDpbv3hEXFyMzANBgkqhkiG9w0BAQsFADBiMQswCQYDVQQGEwJi
ZDESMBAGA1UECAwJYnl0ZWRhbmNlMREwDwYDVQQHDAhzaGFuZ2hhaTELMAkGA1UE
CgwCYmQxCzAJBgNVBAsMAmJkMRIwEAYDVQQDDAlieXRlZGFuY2UwHhcNMjEwODIw
MDkwMDA4WhcNMzEwODE4MDkwMDA4WjBiMQswCQYDVQQGEwJiZDESMBAGA1UECAwJ
Ynl0ZWRhbmNlMREwDwYDVQQHDAhzaGFuZ2hhaTELMAkGA1UECgwCYmQxCzAJBgNV
BAsMAmJkMRIwEAYDVQQDDAlieXRlZGFuY2UwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQCv9f0u0wM3vX7mG5zzmYVLA1mfznX9SuRDaY8ZVe5uU3+MylZx
eURSytlOROYnYkJ515aJjRiTzcCYIskpO8BvpRfbMJP3NP2fseVUPBl9Csg0r7SW
l26SYnZlwswoBJVwVgbz++KHKziNEApAoT/WY5qlAxw6JjCyn/7fh3498ZZIZVHr
e6NfKbNKFpO53UVrXVWXOO4yy9cYJAMr9udluF14OCgzO2977UYBrs44O2+MMf6o
s8t651YGXI/76c8VfAzYJzNRj42ksgayusadGi1MXU4lS4y2gAor5luzgEY/a6Dw
EYj6Ishw3GlOogWxgEG0jWxOkf9a3K5PUMqzAgMBAAEwDQYJKoZIhvcNAQELBQAD
ggEBAIiMOiE/JjjEw2A9qOQkdMPtDSTA2i1bNse1NesLyhY/kQab3H1jeoy4vn35
hlHHf1jGAQVtcLIedVnfiSDZF6S0YtqaF3RklfE3T+KirwFKv0X+8aadURkDdW9Q
Kc0aZdNKk7I1FUwMeGy54LHLsqdfhil8TNFHZUzz7huubemWxG+VPE5bbHuRutX6
eUdBk+Vd+t6fcY7EPEqB5bqFoH0i/Np2mC0ddTLTe9HC7Dt/Z0lEAoJKlCvqIwbN
6Z0SOR+gdRbUkrsfc4h46V6Pl+k9H/G0B58AvxCDZmLdo37oVbC8P/f+96EUdLDC
9SiVX5aU3v+CB9/sw1u9xMqM3BY=
-----END CERTIFICATE-----`
	serverKEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAr/X9LtMDN71+5huc85mFSwNZn851/UrkQ2mPGVXublN/jMpW
cXlEUsrZTkTmJ2JCedeWiY0Yk83AmCLJKTvAb6UX2zCT9zT9n7HlVDwZfQrINK+0
lpdukmJ2ZcLMKASVcFYG8/vihys4jRAKQKE/1mOapQMcOiYwsp/+34d+PfGWSGVR
63ujXymzShaTud1Fa11VlzjuMsvXGCQDK/bnZbhdeDgoMztve+1GAa7OODtvjDH+
qLPLeudWBlyP++nPFXwM2CczUY+NpLIGsrrGnRotTF1OJUuMtoAKK+Zbs4BGP2ug
8BGI+iLIcNxpTqIFsYBBtI1sTpH/WtyuT1DKswIDAQABAoIBAFYDhTPyakxBEVsi
fGKH1SSKCrGXlw2uIk7httXHr7m2j08fSYpGoSNnjTo5a9gsrBorTKcIlc8KcO7t
ci/5yWPJ7BN1d58MBD6gE15N0pfRQXSfv0Mt2xsrjnzH8ExPkU1MgDGcG2u/sPEj
uh9Fs5L1NY8cKhwytvNaSpix/v//3NqxOPRXYNZHOG5AZYxooSAoTwrIzUv+Q1VV
2eIZDHE4C1IE04jKobMnS20an6E1PNllGNyHjpuhP8wlq0dtMiWZ4kNHZrnaWIKr
FMl1j1vMAQFVFfus1EERATyyJPsPNaiZygUhLMByVuBabXXOzOT5sLs7COT+mi/m
RI12e4ECgYEA1LfwUQ+e24AoLLA9+1du7bHZ9wolik0kMdIfpBgaCadKw8FvdR57
WET/F5EfmbWti6qDLLzC5A9A1PTva2M7vvcsN4hM2jj4ci95zpo52f9OMawsPPiL
GCxDzxImwZJ1FqislBdulXRsQxs6e5km/cjJJsc7HTvp+H/q2Tpxz+MCgYEA08Ns
T9yH8FMI86aHeNyMWDrUBWDvRSxc4xhc6wQ9qjtOyGWgSNoD3dMppVe9Kd4HRe0c
R2Bh97qJ7tJJH8LWMn1G2jmf5ti3sRLOOxBl4EylLBKTAqg01j0o8CnMoUeG4YjA
HbMuXQXM1k9TQCRs+iJjwuQFmBxIE8fbxgdTcvECgYEAwSFXsWFLS+MplplVTRPv
BSRKzz6JIu4zBIDQdX3kdgtLmDvR5DYOqu/J4y3B0R7gOERR4JZpJAKtTCwuDAQG
xLVJkgnQLPk4qQNtxiTEjaZ86iB18c1/DC10S1chlPJSGIaAWdyEbHFNsgfQq1M7
0YMxDmIoy7wQC6yoHx4vIx0CgYBGRAgCj9iDg+nqfw1gqz3eYNbNWhbKyyefKwxZ
5zRW8gr4L9B5m+3AgzrEZFKeO9AKAd1qSa74NmtiVWByK4JLiooiCxDl1m9NUEIi
ExTa0lPURe2F6i5uECkvV49QzJ0S5P9qW3Q85ZnMWtHy7KNEdHjJyEOa73dzKNPh
57hm4QKBgQCpv9y5Z6QVko9NUN+lzcQffaFbP1LhKzZfNGBe3ip2WNCHoq+ZMquR
2/W/fcOly/i9PCBYsfes/zETP4Ah53sp1up4C1Z9n/O809LA3XWgluwxw7O68Y1h
JDqUZDGDMvz1AWPCYe+U/BW+RdffuOIMvIJwRrOZeFEiUXiQnlSmAg==
-----END RSA PRIVATE KEY-----`
)

func testNativeTLSServer(t *testing.T) {
	ln, err := tls.Listen("tcp", ":443", getServerTLSConfig())
	MustNil(t, err)
	defer func() {
		err := ln.Close()
		MustNil(t, err)
		t.Logf("native server closed")
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	conn, err := ln.Accept()
	MustNil(t, err)

	buf := make([]byte, 1024)
	go func() {
		defer func() {
			err := conn.Close()
			MustNil(t, err)
			wg.Done()
		}()

		for {
			n, err := conn.Read(buf)
			if errors.Is(err, io.EOF) {
				return
			}
			MustNil(t, err)

			n, err = conn.Write(buf[:n])
			MustNil(t, err)
		}
	}()

	wg.Wait()
}

func testTLSServer(t *testing.T) {
	ln, err := CreateListener("tcp", ":443", getServerTLSConfig())
	MustNil(t, err)
	defer func() {
		err := ln.Close()
		MustNil(t, err)
		t.Logf("server closed")
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	conn, err := ln.Accept()
	MustNil(t, err)
	tlsConn, err := GetConnection(conn)
	MustNil(t, err)
	r, w := tlsConn.Reader(), tlsConn.Writer()

	msg := "hello"
	go func() {
		defer func() {
			err := tlsConn.Close()
			MustNil(t, err)
			wg.Done()
		}()

		for {
			resp, err := r.ReadString(len(msg))
			if errors.Is(err, netpoll.ErrEOF) {
				return
			}
			if errors.Is(err, io.EOF) {
				return
			}
			MustNil(t, err)
			MustTrue(t, msg == resp)
			err = r.Release()
			MustNil(t, err)

			_, err = w.WriteString(msg)
			MustNil(t, err)
			err = w.Flush()
			MustNil(t, err)
		}
	}()

	wg.Wait()
}

func testTLSClient(t *testing.T) {
	conn, err := DialConnection("tcp", ":443", time.Second, getClientTLSConfig())
	MustNil(t, err)
	r, w := conn.Reader(), conn.Writer()
	defer func() {
		err := conn.Close()
		MustNil(t, err)
		t.Logf("client closed")
	}()

	msg := "hello"
	for i := 0; i < 1024; i++ {
		_, err := w.WriteString(msg)
		MustNil(t, err)
		err = w.Flush()
		MustNil(t, err)

		resp, err := r.ReadString(len(msg))
		MustNil(t, err)
		MustTrue(t, resp == msg)
		err = r.Release()
		MustNil(t, err)
	}
}

func testNativeTLSClient(t *testing.T) {
	conn, err := tls.Dial("tcp", ":443", getClientTLSConfig())
	MustNil(t, err)
	defer func() {
		err := conn.Close()
		MustNil(t, err)
		t.Logf("native client closed")
	}()

	buf := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		_, err := conn.Write([]byte("hello"))
		MustNil(t, err)

		n, err := conn.Read(buf)
		MustNil(t, err)
		MustTrue(t, n > 0)
	}
}

func getServerTLSConfig() *tls.Config {
	serverCert, _ := tls.X509KeyPair([]byte(serverPEM), []byte(serverKEY))
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM([]byte(clientRootREM))
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	return config
}

func getClientTLSConfig() *tls.Config {
	clientCert, _ := tls.X509KeyPair([]byte(clientPEM), []byte(clientKEY))
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM([]byte(clientRootREM))
	config := &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{clientCert},
		ServerName:         "p.s.m",
		InsecureSkipVerify: true,
	}
	return config
}

func MustNil(t *testing.T, val interface{}) {
	if t == nil {
		return
	}
	t.Helper()
	Assert(t, val == nil, val)
	if val != nil {
		t.Fatal("assertion nil failed, val=", val)
	}
}

func MustTrue(t *testing.T, cond bool) {
	if t == nil {
		return
	}
	t.Helper()
	if !cond {
		t.Fatal("assertion true failed.")
	}
}

func Equal(t *testing.T, got, expect interface{}) {
	if t == nil {
		return
	}
	t.Helper()
	if got != expect {
		t.Fatalf("assertion equal failed, got=[%v], expect=[%v]", got, expect)
	}
}

func Assert(t *testing.T, cond bool, val ...interface{}) {
	if t == nil {
		return
	}
	t.Helper()
	if !cond {
		if len(val) > 0 {
			val = append([]interface{}{"assertion failed:"}, val...)
			t.Fatal(val...)
		} else {
			t.Fatal("assertion failed")
		}
	}
}
