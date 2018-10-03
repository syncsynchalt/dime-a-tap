package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

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
	// in PEM format
	CAKey []byte
	// in PEM format
	CACert []byte
	CADir  string
}

func Listen(opts Opts) error {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(opts.Port))
	if err != nil {
		return err
	}
	snooplisten := ListenWrap{Listener: ln, opts: opts}
	l.Printf("started listen on port %d\n", opts.Port)
	tlslisten := tls.NewListener(&snooplisten, &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return getCertificate(hello, l, &opts)
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

func handleConnection(clientConn net.Conn, opts *Opts) {
	defer clientConn.Close()

	clientName := clientConn.RemoteAddr().String()
	l := log.New(os.Stdout, clientName+" ", log.Ldate|log.Ltime)

	tlsConn, ok := clientConn.(*tls.Conn)
	if !ok {
		panic("Unable to convert connection to tls.Conn, shouldn't happen")
	}
	// needed to set up ServerName
	err := tlsConn.Handshake()
	if err != nil {
		l.Println("error performing handshake", err)
		return
	}
	serverName := tlsConn.ConnectionState().ServerName
	if serverName == "" {
		l.Println("client did not send hostname (SNI), unable to proceed, closing")
		return
	}
	l.Printf("intercepted connection to %s:%d\n", serverName, opts.Port)

	serverConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", serverName, opts.Port), &tls.Config{
		ServerName: serverName,
	})
	if err != nil {
		l.Printf("unable to connect to %s: %s\n", serverName, err)
		return
	}
	defer serverConn.Close()
	l.Printf("connected to %s:%d\n", serverName, opts.Port)

	readFromClient := make(chan []byte)
	readFromServer := make(chan []byte)

	go func() {
		// reads from clientConn and writes to readFromClient
		b := make([]byte, 4096)
		for {
			n, err := clientConn.Read(b)
			if n != 0 {
				bcopy := make([]byte, n)
				copy(bcopy, b[:n])
				readFromClient <- bcopy
			}
			if err != nil {
				break
			}
		}
		close(readFromClient)
	}()

	go func() {
		// reads from serverConn and writes to readFromServer
		b := make([]byte, 4096)
		for {
			n, err := serverConn.Read(b)
			if n != 0 {
				bcopy := make([]byte, n)
				copy(bcopy, b[:n])
				readFromServer <- bcopy
			}
			if err != nil {
				break
			}
		}
		close(readFromServer)
	}()

	for {
		select {
		case b, more := <-readFromClient:
			for len(b) > 0 {
				n, err := serverConn.Write(b)
				if err != nil {
					l.Println("unable to write data to server:", err)
					return
				}
				b = b[n:]
			}
			if !more {
				l.Println("client conn closed")
				return
			}
		case b, more := <-readFromServer:
			for len(b) > 0 {
				n, err := clientConn.Write(b)
				if err != nil {
					l.Println("unable to write data to server:", err)
					return
				}
				b = b[n:]
			}
			if !more {
				l.Println("server conn closed")
				return
			}
		}
	}
}

func getCertificate(hello *tls.ClientHelloInfo, l *log.Logger, opts *Opts) (*tls.Certificate, error) {
	if hello.ServerName == "" {
		l.Println("returning generic cert because client did not provide hostname in SNI")
	} else {
		l.Printf("returning certificate for %s\n", hello.ServerName)
	}

	// xxx todo
	cert, err := tls.X509KeyPair(opts.CACert, opts.CAKey)
	return &cert, err
}
