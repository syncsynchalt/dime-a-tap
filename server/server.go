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

// intercepts Accept() and introduces a SnoopConn
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

type Opts struct {
	// port to listen on
	Port int
	// optional dir to store raw read/write info
	RawDir string
	// optional dir to store unencrypted read/write info
	CaptureDir string
	// if not set then creates memory-only version
	CADir string
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

func handleConnection(clientConn net.Conn, opts *Opts) (err error) {
	defer clientConn.Close()

	clientName := clientConn.RemoteAddr().String()
	l := log.New(os.Stdout, clientName+" ", log.Ldate|log.Ltime)
	defer func() {
		if err != nil {
			l.Println(err)
		}
	}()

	tlsConn, ok := clientConn.(*tls.Conn)
	if !ok {
		return fmt.Errorf("unable to convert connection to tls.Conn")
	}
	// needed to set up ServerName
	err = tlsConn.Handshake()
	if err != nil {
		return fmt.Errorf("error performing handshake: %s", err)
	}
	serverName := tlsConn.ConnectionState().ServerName
	if serverName == "" {
		return fmt.Errorf("client did not send hostname (SNI), unable to proceed, closing")
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

	return rwpipe.PipeConns(clientConn, serverConn, opts.CaptureDir)
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
