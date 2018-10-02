package main

import (
	"flag"
	"fmt"
	"github.com/syncsynchalt/dime-a-tap/listen"
	"os"
	"strconv"
)

func dieUsage(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Fprintf(os.Stderr, "Usage: %s [port]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	handshakes := flag.String("handshakes", "",
		"directory to write client/server handshake information for each connection")
	flag.Parse()
	if flag.NArg() != 1 {
		dieUsage(fmt.Errorf("No listen port specified"))
	}

	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		dieUsage(err)
	}
	fmt.Println("port is", port)

	err = listen.Listen(listen.Opts{
		Port:       port,
		Handshakes: *handshakes,
	})
	if err != nil {
		panic(err)
	}
}
