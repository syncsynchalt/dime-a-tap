package listen

import (
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
	Port   int
	RawDir string
}

func Listen(opts Opts) error {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(opts.Port))
	mylisten := ListenWrap{Listener: ln, opts: opts}
	if err != nil {
		return err
	}
	l.Printf("started listen on port %d\n", opts.Port)
	for {
		conn, err := mylisten.Accept()
		if err != nil {
			l.Panicln("unable to accept connection:", err)
		}
		l.Printf("accepted connection from %s", conn.RemoteAddr().String())
		go handleConnection(conn, opts)
	}
}

func handleConnection(conn net.Conn, opts Opts) {
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
