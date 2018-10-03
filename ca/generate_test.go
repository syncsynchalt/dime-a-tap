package ca_test

import (
	"testing"

	"github.com/syncsynchalt/dime-a-tap/ca"
	"github.com/syncsynchalt/dime-a-tap/test"
)

func TestGenKey(t *testing.T) {
	key, err := ca.GenerateCaKey()
	test.Ok(t, err)
	t.Log(key)
	prefix := "-----BEGIN RSA PRIVATE KEY-----\nM"
	test.Equals(t, prefix, string(key[:len(prefix)]))
	// fmt.Print(key)
}

func TestGenCert(t *testing.T) {
	key, err := ca.GenerateCaKey()
	test.Ok(t, err)
	cert, err := ca.GenerateCaCert(key)
	test.Ok(t, err)
	prefix := "-----BEGIN CERTIFICATE-----\nM"
	test.Equals(t, prefix, string(cert[:len(prefix)]))
	// fmt.Print(cert)
}
