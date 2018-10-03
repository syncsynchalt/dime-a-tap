package ca_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/syncsynchalt/dime-a-tap/ca"
	"github.com/syncsynchalt/dime-a-tap/test"
)

func TestGenKey(t *testing.T) {
	key, err := ca.GenerateCAKey()
	test.Ok(t, err)
	t.Log(key)
	prefix := "-----BEGIN RSA PRIVATE KEY-----\nM"
	test.Equals(t, prefix, string(key[:len(prefix)]))
	// fmt.Print(key)
}

func TestGenCert(t *testing.T) {
	key, err := ca.GenerateCAKey()
	test.Ok(t, err)
	cert, err := ca.GenerateCACert(key)
	test.Ok(t, err)
	prefix := "-----BEGIN CERTIFICATE-----\nM"
	test.Equals(t, prefix, string(cert[:len(prefix)]))
	// fmt.Print(cert)
}

func TestGenCa(t *testing.T) {
	mydir := fmt.Sprintf("/tmp/golang.test.%d", time.Now().UnixNano())
	defer os.RemoveAll(mydir)
	err := ca.CreateCAStore(mydir)
	test.Ok(t, err)

	certdata, err := ioutil.ReadFile(mydir + "/ca.crt")
	test.Ok(t, err)
	prefix := "-----BEGIN CERTIFICATE-----\nM"
	test.Equals(t, prefix, string(certdata)[:len(prefix)])

	keydata, err := ioutil.ReadFile(mydir + "/ca.key")
	test.Ok(t, err)
	prefix = "-----BEGIN RSA PRIVATE KEY-----\nM"
	test.Equals(t, prefix, string(keydata)[:len(prefix)])
}

func TestGenCaNotDir(t *testing.T) {
	mydir := fmt.Sprintf("/tmp/golang.test.%d", time.Now().UnixNano())
	defer os.RemoveAll(mydir)
	ioutil.WriteFile(mydir, []byte("abc\n"), 0644)

	err := ca.CreateCAStore(mydir)
	test.Assert(t, strings.HasSuffix(err.Error(), "not a directory"), "error [%s] not expected format", err)
}

func TestGenCaNoTwice(t *testing.T) {
	mydir := fmt.Sprintf("/tmp/golang.test.%d", time.Now().UnixNano())
	defer os.RemoveAll(mydir)

	err := ca.CreateCAStore(mydir)
	test.Ok(t, err)
	err = ca.CreateCAStore(mydir)
	test.Assert(t, strings.HasSuffix(err.Error(), "file exists"), "error [%s] not expected format", err)
}
