package rwpipe

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/syncsynchalt/dime-a-tap/disklog"
)

// pipe the inputs and outputs of each conn to each other
func PipeConns(conn1 net.Conn, name1 string, conn2 net.Conn, name2 string, captureDir string) error {
	logName := conn1.RemoteAddr().String()
	l := log.New(os.Stdout, logName+" ", log.Ldate|log.Ltime)

	readFrom1 := make(chan []byte)
	go func() {
		// reads from conn1 and writes to readFrom1
		b := make([]byte, 4096)
		for {
			n, err := conn1.Read(b)
			if n != 0 {
				bcopy := make([]byte, n)
				copy(bcopy, b[:n])
				readFrom1 <- bcopy
			}
			if err != nil {
				break
			}
		}
		close(readFrom1)
	}()

	readFrom2 := make(chan []byte)
	go func() {
		// reads from conn2 and writes to readFrom2
		b := make([]byte, 4096)
		for {
			n, err := conn2.Read(b)
			if n != 0 {
				bcopy := make([]byte, n)
				copy(bcopy, b[:n])
				readFrom2 <- bcopy
			}
			if err != nil {
				break
			}
		}
		close(readFrom2)
	}()

loop:
	for {
		select {
		case b, more := <-readFrom1:
			err := disklog.DumpPacket(captureDir, logName, "c", b)
			if err != nil {
				l.Println("unable to dump clean:", err)
				// ignore error
			}

			for len(b) > 0 {
				n, err := conn2.Write(b)
				if err != nil {
					return fmt.Errorf("unable to write data to %s: %s", name2, err)
				}
				b = b[n:]
			}

			if !more {
				l.Printf("%s conn closed\n", name1)
				break loop
			}
		case b, more := <-readFrom2:
			err := disklog.DumpPacket(captureDir, logName, "s", b)
			if err != nil {
				l.Println("unable to dump clean:", err)
				// ignore error
			}

			for len(b) > 0 {
				n, err := conn1.Write(b)
				if err != nil {
					return fmt.Errorf("unable to write data to %s: %s", name1, err)
				}
				b = b[n:]
			}
			if !more {
				l.Printf("%s conn closed\n", name2)
				break loop
			}
		}
	}
	return nil
}
