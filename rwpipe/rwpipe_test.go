package rwpipe_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/syncsynchalt/dime-a-tap/rwpipe"
	"github.com/syncsynchalt/dime-a-tap/server"
	"github.com/syncsynchalt/dime-a-tap/test"
)

func makeConnPair(t *testing.T, port int) (net.Conn, net.Conn) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	test.Ok(t, err)

	c1 := make(chan net.Conn)
	go func() {
		conn1, err := l.Accept()
		test.Ok(t, err)
		c1 <- conn1
	}()
	c2 := make(chan net.Conn)
	go func() {
		conn2, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		test.Ok(t, err)
		c2 <- conn2
	}()
	conn1 := <-c1
	conn2 := <-c2
	l.Close()
	return conn1, conn2
}

func TestRWPipe(t *testing.T) {

	// conn1 <-tcp-> conn2 <-rwpipe-> conn3 <-tcp-> conn4

	conn1, conn2, err := server.MakeConnPair(8314)
	test.Ok(t, err)
	conn3, conn4, err := server.MakeConnPair(8315)
	test.Ok(t, err)

	// the data we'll send on each end
	data1 := strings.Split("asdfasdf\nasfasdfasdf\nasdfasdfASdfasdf\nasdfasdfasDFasdf\nasdfasDFasdfasdf\n", "\n")
	data4 := strings.Split("iuewriouy\nwieoruower\nnaseroiuw\nwerqweroiqweroiwqeiuor\noweuruoweiruwer\n", "\n")

	// join conn2 to conn3
	done := make(chan bool)
	go func() {
		rwpipe.PipeConns(conn2, "conn2", conn3, "conn3", "")
		conn2.Close()
		conn3.Close()
		done <- true
	}()

	// alternate write and read on conn1
	chan1 := make(chan string)
	go func() {
		defer conn1.Close()
		from1 := make([]byte, 0)
		scan1 := bufio.NewScanner(conn1)
		for _, s := range data1 {
			_, err := conn1.Write([]byte(s + "\n"))
			test.Ok(t, err)
			scanned := scan1.Scan()
			test.Assert(t, scanned, "conn1 didn't scan")
			test.Ok(t, scan1.Err())
			from1 = append(from1, scan1.Bytes()...)
			from1 = append(from1, "\n"...)
		}
		chan1 <- string(from1)
	}()

	// alternate write and read on conn4
	chan4 := make(chan string)
	go func() {
		defer conn4.Close()
		from4 := make([]byte, 0)
		scan4 := bufio.NewScanner(conn4)
		for _, s := range data4 {
			_, err := conn4.Write([]byte(s + "\n"))
			test.Ok(t, err)
			scanned := scan4.Scan()
			test.Assert(t, scanned, "conn4 didn't scan")
			test.Ok(t, scan4.Err())
			from4 = append(from4, scan4.Bytes()...)
			from4 = append(from4, "\n"...)
		}
		chan4 <- string(from4)
	}()

	read1 := <-chan1
	read4 := <-chan4
	_ = <-done

	test.Equals(t, strings.Join(data4, "\n")+"\n", read1)
	test.Equals(t, strings.Join(data1, "\n")+"\n", read4)
}
