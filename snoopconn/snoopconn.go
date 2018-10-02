package snoopconn

// intercepts the net.Conn interface and records the first 10kiB of data.
// Optionally writes all reads/writes to RawDir

import (
	"log"
	"net"
	"os"

	"github.com/syncsynchalt/dime-a-tap/disklog"
)

const SnoopBytes = 10240

type TattleConn struct {
	net.Conn
	rawDir     string
	log        *log.Logger
	remoteName string
	ReadData   []byte
	WriteData  []byte
}

func New(conn net.Conn, rawDir string) *TattleConn {
	return &TattleConn{
		Conn:       conn,
		rawDir:     rawDir,
		log:        log.New(os.Stdout, conn.RemoteAddr().String()+" ", log.Ldate|log.Ltime),
		remoteName: conn.RemoteAddr().String(),
		ReadData:   make([]byte, 0, SnoopBytes),
		WriteData:  make([]byte, 0, SnoopBytes),
	}
}

func (c *TattleConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		return n, err
	}

	err = disklog.DumpPacket(c.rawDir, c.remoteName, "c", b[:n])
	if err != nil {
		log.Println("unable to dump raw:", err)
		// ignore error
	}

	if len(c.ReadData) < SnoopBytes {
		na := n
		if len(c.ReadData)+na > SnoopBytes {
			na = SnoopBytes - len(c.ReadData)
		}
		c.ReadData = append(c.ReadData, b[:na]...)
	}
	return n, nil
}

func (c *TattleConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		return n, err
	}

	err = disklog.DumpPacket(c.rawDir, c.remoteName, "s", b[:n])
	if err != nil {
		log.Println("unable to dump raw:", err)
		// ignore error
	}

	if len(c.WriteData) < 10240 {
		na := n
		if len(c.WriteData)+na > SnoopBytes {
			na = SnoopBytes - len(c.WriteData)
		}
		c.WriteData = append(c.WriteData, b[:na]...)
	}
	return n, nil
}
