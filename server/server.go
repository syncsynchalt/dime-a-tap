package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/syncsynchalt/dime-a-tap/ca"
	"github.com/syncsynchalt/dime-a-tap/rwpipe"
	"github.com/syncsynchalt/dime-a-tap/snoopconn"
)

type Opts struct {
	// port to listen on
	Port int
	// optional dir to store raw read/write info
	RawDir string
	// optional dir to store unencrypted read/write info
	CaptureDir string
	// if not set then creates memory-only version
	CADir string
	// unencrypted data will be sent over localhost:4430 for tcpdump-ability
	TapPort int
}

// intercepts Accept() and wraps Conn in a SnoopConn
type ListenWrap struct {
	net.Listener
	opts Opts
}

// override of net.Listener.Accept()
func (l *ListenWrap) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	tc := snoopconn.New(conn, l.opts.RawDir)
	return tc, err
}

func Listen(opts Opts) error {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	caStore, err := ca.NewStore(opts.CADir)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(opts.Port))
	if err != nil {
		return err
	}
	snooplisten := ListenWrap{Listener: ln, opts: opts}
	l.Printf("started listen on port %d\n", opts.Port)
	tlslisten := tls.NewListener(&snooplisten, &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return getCertificate(hello, l, caStore)
		},
	})

	for {
		conn, err := tlslisten.Accept()
		if err != nil {
			l.Panicln("unable to accept connection:", err)
		}
		go handleConnection(conn, &opts)
	}
}

func getSNIServerName(conn net.Conn) (string, error) {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return "", fmt.Errorf("unable to convert connection to tls.Conn")
	}
	// needed to set up ServerName
	err := tlsConn.Handshake()
	if err != nil {
		return "", fmt.Errorf("error performing handshake: %s", err)
	}
	serverName := tlsConn.ConnectionState().ServerName
	if serverName == "" {
		return "", fmt.Errorf("client did not send hostname (SNI), unable to proceed")
	}
	return serverName, nil
}

func handleConnection(clientConn net.Conn, opts *Opts) (err error) {
	defer clientConn.Close()

	// create and connect the following:
	// server <-tcp-> serverConn <-rwpipe-> tapPort <-tcp-> tapAnon <-rwpipe-> clientConn <-tcp-> client
	//   tapPort: localhost:4430
	//   tapAnon: localhost:{ephemeral}

	clientName := clientConn.RemoteAddr().String()
	l := log.New(os.Stdout, clientName+" ", log.Ldate|log.Ltime)
	defer func() {
		if err != nil {
			l.Println(err)
		}
	}()

	serverName, err := getSNIServerName(clientConn)
	if err != nil {
		return err
	}
	l.Printf("intercepted connection to %s:%d\n", serverName, opts.Port)

	serverConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", serverName, opts.Port), &tls.Config{
		ServerName: serverName,
	})
	if err != nil {
		return fmt.Errorf("unable to connect to %s: %s\n", serverName, err)
	}
	defer serverConn.Close()
	l.Printf("connected to %s:%d\n", serverName, opts.Port)

	tapPort, tapAnon, err := MakeConnPair(opts.TapPort)
	if err != nil {
		return err
	}
	defer tapPort.Close()
	defer tapAnon.Close()

	go func() {
		err = rwpipe.PipeConns(serverConn, "server", tapPort, "tap", "")
		if err != nil {
			l.Println(err)
		}
	}()
	return rwpipe.PipeConns(clientConn, "client", tapAnon, "tapanon", opts.CaptureDir)
}

func getCertificate(hello *tls.ClientHelloInfo, l *log.Logger, caStore *ca.Store) (*tls.Certificate, error) {
	if hello.ServerName == "" {
		return nil, fmt.Errorf("server did not provide hostname in SNI")
	}
	l.Printf("returning certificate for %s\n", hello.ServerName)

	key, cert, err := caStore.GetCertificate(hello.ServerName)
	kp, err := tls.X509KeyPair(cert, key)
	return &kp, err
}

// given a port, creates a pair of connections in the form of
// localhost:port and localhost:ephemeralport that are connected
// to each other
func MakeConnPair(port int) (onPort, anonPort net.Conn, err error) {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, nil, err
	}
	defer l.Close()

	type connAndError struct {
		conn net.Conn
		err  error
	}

	c1 := make(chan *connAndError)
	go func() {
		conn1, err := l.Accept()
		c1 <- &connAndError{conn1, err}
	}()
	c2 := make(chan *connAndError)
	go func() {
		conn2, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		c2 <- &connAndError{conn2, err}
	}()

	ce1 := <-c1
	ce2 := <-c2
	if ce1.err != nil || ce2.err != nil {
		if ce1.conn != nil {
			ce1.conn.Close()
		}
		if ce2.conn != nil {
			ce2.conn.Close()
		}
		if ce1.err != nil {
			return nil, nil, ce1.err
		} else {
			return nil, nil, ce2.err
		}
	}
	return ce1.conn, ce2.conn, nil
}
