package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/syncsynchalt/dime-a-tap/ca"
	"github.com/syncsynchalt/dime-a-tap/server"
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
	if len(os.Args) == 3 && os.Args[1] == "ca-init" {
		err := ca.CreateCAStore(os.Args[2])
		if err != nil {
			dieUsage(err)
		} else {
			fmt.Println("success")
			os.Exit(0)
		}
	}

	caDir := flag.String("cadir", "", "optional path to CA key store (use 'dime-a-tap ca-init {dir}' to create)")
	tapPort := flag.Int("tapport", 4430, "localhost port to send unencrypted data over")
	rawDir := flag.String("rawdir", "", "optional directory to write raw data written to/from client")
	capDir := flag.String("capturedir", "", "optional directory to capture unencrypted data written to/from client")
	flag.Parse()
	if flag.NArg() != 1 {
		dieUsage(fmt.Errorf("No listen port specified"))
	}

	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		dieUsage(err)
	}

	opts := server.Opts{
		Port:       port,
		RawDir:     *rawDir,
		CaptureDir: *capDir,
		CADir:      *caDir,
		TapPort:    *tapPort,
	}

	err = server.Listen(opts)
	if err != nil {
		panic(err)
	}
}
