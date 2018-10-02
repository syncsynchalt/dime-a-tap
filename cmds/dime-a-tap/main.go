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
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] [port]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	rawDir := flag.String("rawdir", "", "optional directory to write raw data written to/from client")
	flag.Parse()
	if flag.NArg() != 1 {
		dieUsage(fmt.Errorf("No listen port specified"))
	}

	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		dieUsage(err)
	}

	err = listen.Listen(listen.Opts{
		Port:   port,
		RawDir: *rawDir,
	})
	if err != nil {
		panic(err)
	}
}
