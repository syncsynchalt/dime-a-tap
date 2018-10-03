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
	if len(os.Args) == 3 && os.Args[1] == "ca-init" {
		err := ca.CreateCAStore(os.Args[2])
		if err != nil {
			dieUsage(err)
		} else {
			fmt.Println("success")
			os.Exit(0)
		}
	}

	rawDir := flag.String("rawdir", "", "optional directory to write raw data written to/from client")
	caCert := flag.String("cacert", "", "optional path to CA certificate file")
	caKey := flag.String("cakey", "", "optional path to CA private key file")
	caDir := flag.String("cadir", "", "optional path to CA key store (use 'dime-a-tap ca-init {dir}' to create)")
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
		CADir:  *caDir,
	}

	opts.CAKey, err = loadCAKey(*caKey, *caDir)
	if err != nil {
		dieUsage(err)
	}

	opts.CACert, err = loadCACert(*caCert, *caDir)
	if err != nil {
		dieUsage(err)
	}

	err = server.Listen(opts)
	if err != nil {
		panic(err)
	}
}

func loadCAKey(caKey, caDir string) ([]byte, error) {
	if caKey != "" {
		return ioutil.ReadFile(caKey)
	} else if caDir != "" {
		return ioutil.ReadFile(caDir + "/ca.key")
	} else {
		return ca.GenerateCAKey()
	}
}

func loadCACert(caKey, caDir string) ([]byte, error) {
	if caKey != "" {
		return ioutil.ReadFile(caKey)
	} else if caDir != "" {
		return ioutil.ReadFile(caDir + "/ca.crt")
	} else {
		return ca.GenerateCAKey()
	}
}
