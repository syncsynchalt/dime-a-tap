package listen

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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
	CaKey []byte
	// in PEM format
	CaCert []byte
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
		l.Printf("accepted connection from %s", conn.RemoteAddr().String())
		go handleConnection(conn, &opts)
	}
}

func handleConnection(conn net.Conn, opts *Opts) {
	remoteName := conn.RemoteAddr().String()
	l := log.New(os.Stdout, remoteName+" ", log.Ldate|log.Ltime)
	buf := make([]byte, 10240)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			l.Println("unable to read:", err)
			break
		}

		data := buf[:n]
		for len(data) > 0 && (data[len(data)-1] == '\r' || data[len(data)-1] == '\n') {
			data = data[:len(data)-1]
		}
		l.Printf("Read [%s]\n", data)
		if strings.ToUpper(string(data)) == "QUIT" {
			break
		}
		conn.Write([]byte(fmt.Sprintf("echoing:%s\r\n", data)))
	}
	conn.Close()
}

func getCertificate(hello *tls.ClientHelloInfo, l *log.Logger, opts *Opts) (*tls.Certificate, error) {
	if hello.ServerName == "" {
		l.Println("returning generic cert because client did not provide hostname in SNI")
	} else {
		l.Printf("returning certificate for %s\n", hello.ServerName)
	}

	// xxx todo
	cert, err := tls.X509KeyPair(opts.CaCert, opts.CaKey)
	return &cert, err
}
