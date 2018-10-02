package snoopconn_test

import (
	"io"
	"net"
	"strings"
	"testing"

	"github.com/syncsynchalt/dime-a-tap/snoopconn"
	"github.com/syncsynchalt/dime-a-tap/test"
)

func TestTattleConn(t *testing.T) {
	l, err := net.Listen("tcp", ":5436")
	test.Ok(t, err)
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Log("gofunc error accepting:", err)
		}
		defer conn.Close()
		b := make([]byte, 1000)
		_, err = conn.Read(b)
		if err != nil {
			t.Log("gofunc error reading:", err)
		}
		_, err = conn.Write([]byte("RESPONSE\r\n"))
		if err != nil {
			t.Log("gofunc error writing:", err)
		}
	}()

	conn, err := net.Dial("tcp", ":5436")
	test.Ok(t, err)
	tc := snoopconn.New(conn, "")
	defer tc.Close()
	n, err := tc.Write([]byte("QUERY\r\n"))
	test.Ok(t, err)
	test.Equals(t, 7, n)

	b := make([]byte, 1000)
	n, err = tc.Read(b)
	test.Ok(t, err)
	test.Equals(t, 10, n)
	test.Equals(t, []byte("RESPONSE\r\n"), b[:n])
	test.Equals(t, []byte("QUERY\r\n"), tc.WriteData)
	test.Equals(t, []byte("RESPONSE\r\n"), tc.ReadData)
}

func TestTattleConnMulti(t *testing.T) {
	l, err := net.Listen("tcp", ":5436")
	test.Ok(t, err)
	defer l.Close()
	const testSize = 100000

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Log("gofunc error accepting:", err)
		}
		defer conn.Close()
		bb := make([]byte, 805)
		i := 0
		for i < testSize {
			n1, err := conn.Read(bb)
			if err != nil {
				t.Log("Stopping reads from socket at i=", i, err)
				break
			}
			n2, err := conn.Write([]byte(strings.Repeat("y", n1)))
			if err != nil {
				t.Log("Stopping writes from socket at i=", i, err)
				break
			}
			if n1 != n2 {
				t.Log("Mismatched reads and writes at i=", i)
				break
			}
			i += n1
		}
	}()

	conn, err := net.Dial("tcp", ":5436")
	test.Ok(t, err)
	tc := snoopconn.New(conn, "")
	defer tc.Close()

	send := []byte(strings.Repeat("x", testSize))
	receive := []byte("")

	b := make([]byte, 1024)
	for {
		toSend := send[:]
		if len(toSend) > 2048 {
			toSend = toSend[:2048]
		}
		n, err := tc.Write(toSend)
		test.Ok(t, err)
		// t.Logf("Sent %d bytes", n)
		send = send[n:]

		n, err = tc.Read(b)
		if err == io.EOF {
			break
		}
		test.Ok(t, err)
		// t.Logf("Read %d bytes", n)
		receive = append(receive, b[:n]...)
	}

	test.Equals(t, testSize, len(receive))
	test.Equals(t, []byte(strings.Repeat("y", testSize)), receive)

	test.Equals(t, 10240, len(tc.WriteData))
	test.Equals(t, []byte(strings.Repeat("x", 10240)), tc.WriteData)
	test.Equals(t, 10240, len(tc.ReadData))
	test.Equals(t, []byte(strings.Repeat("y", 10240)), tc.ReadData)
}
