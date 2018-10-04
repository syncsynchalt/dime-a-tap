package rwpipe

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/syncsynchalt/dime-a-tap/disklog"
)

// pipe the inputs and outputs of each conn to each other
func PipeConns(client, server net.Conn, captureDir string) error {
	clientName := client.RemoteAddr().String()
	l := log.New(os.Stdout, clientName+" ", log.Ldate|log.Ltime)

	readFromClient := make(chan []byte)
	go func() {
		// reads from client and writes to readFromClient
		b := make([]byte, 4096)
		for {
			n, err := client.Read(b)
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

	readFromServer := make(chan []byte)
	go func() {
		// reads from server and writes to readFromServer
		b := make([]byte, 4096)
		for {
			n, err := server.Read(b)
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

loop:
	for {
		select {
		case b, more := <-readFromClient:
			err := disklog.DumpPacket(captureDir, clientName, "c", b)
			if err != nil {
				l.Println("unable to dump raw:", err)
				// ignore error
			}

			for len(b) > 0 {
				n, err := server.Write(b)
				if err != nil {
					return fmt.Errorf("unable to write data to server: %s", err)
				}
				b = b[n:]
			}

			if !more {
				l.Println("client conn closed")
				break loop
			}
		case b, more := <-readFromServer:
			err := disklog.DumpPacket(captureDir, clientName, "s", b)
			if err != nil {
				l.Println("unable to dump raw:", err)
				// ignore error
			}

			for len(b) > 0 {
				n, err := client.Write(b)
				if err != nil {
					return fmt.Errorf("unable to write data to client: %s", err)
				}
				b = b[n:]
			}
			if !more {
				l.Println("server conn closed")
				break loop
			}
		}
	}
	return nil
}
