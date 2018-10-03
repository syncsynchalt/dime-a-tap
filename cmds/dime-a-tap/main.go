package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	rawDir := flag.String("rawdir", "", "optional directory to write raw data written to/from client")
	caCert := flag.String("cacert", "", "optional path to CA certificate file")
	caKey := flag.String("cakey", "", "optional path to CA private key file")
	flag.Parse()
	if flag.NArg() != 1 {
		dieUsage(fmt.Errorf("No listen port specified"))
	}

	port, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		dieUsage(err)
	}

	opts := server.Opts{
		Port:   port,
		RawDir: *rawDir,
	}

	if *caKey != "" {
		opts.CaKey, err = ioutil.ReadFile(*caKey)
	} else {
		opts.CaKey, err = ca.GenerateCaKey()
	}
	if err != nil {
		dieUsage(err)
	}

	if *caCert != "" {
		opts.CaCert, err = ioutil.ReadFile(*caCert)
	} else {
		opts.CaCert, err = ca.GenerateCaCert(opts.CaKey)
	}
	if err != nil {
		dieUsage(err)
	}

	err = server.Listen(opts)
	if err != nil {
		panic(err)
	}
}
