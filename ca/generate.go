package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
)

const rsaBits = 2048

// generate a RSA key and return it as a PEM block
func GenerateCAKey() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, err
	}
	derKey := x509.MarshalPKCS1PrivateKey(key)
	return derToPem(derKey, "RSA PRIVATE KEY")
}

// generate a self-signed certificate for key suitable for CA use.
// If these defaults aren't suitable, build your own using openssl or similar.
func GenerateCACert(pemKey []byte) ([]byte, error) {
	derKey, _ := pem.Decode([]byte(pemKey))
	if derKey == nil {
		return nil, fmt.Errorf("unable to decode private key in PEM format: %s", pemKey)
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(derKey.Bytes)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"Dime-A-Tap"},
			OrganizationalUnit: []string{"Fake CA"},
		},
		NotBefore:             time.Now().Add(-3600).UTC(),
		NotAfter:              time.Now().AddDate(10, 0, 0).UTC(),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	template.SubjectKeyId, err = generateSubjectKeyId(rsaKey)
	if err != nil {
		return nil, err
	}

	derCert, err := x509.CreateCertificate(rand.Reader, &template, &template, rsaKey.Public(), rsaKey)
	if err != nil {
		return nil, err
	}

	return derToPem(derCert, "CERTIFICATE")
}

func derToPem(bytes []byte, pemType string) ([]byte, error) {
	sout := strings.Builder{}
	err := pem.Encode(&sout, &pem.Block{Type: pemType, Bytes: bytes})
	return []byte(sout.String()), err
}

// required for CA certs, we generate ours by SHA1(pubkey)
func generateSubjectKeyId(key *rsa.PrivateKey) ([]byte, error) {
	// hash the public key for the subjectKeyId
	pubKeyOnly := rsa.PublicKey{N: key.PublicKey.N, E: key.PublicKey.E}
	bytes, err := asn1.Marshal(pubKeyOnly)
	if err != nil {
		return nil, err
	}
	hash := sha1.Sum(bytes)
	return hash[:], nil
}

func CreateCAStore(directory string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	key, err := GenerateCAKey()
	if err != nil {
		return err
	}
	cert, err := GenerateCACert(key)
	if err != nil {
		return err
	}

	keyfile := directory + "/ca.key"
	certfile := directory + "/ca.crt"

	err = writeFileExcl(keyfile, key, 0600)
	if err != nil {
		return err
	}
	return writeFileExcl(certfile, cert, 0644)
}

// write data to file, which must not already exist
func writeFileExcl(filename string, data []byte, perm os.FileMode) error {
	kf, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer kf.Close()
	_, err = kf.Write(data)
	return err
}
