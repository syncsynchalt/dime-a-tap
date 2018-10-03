package ca_test

import (
	"fmt"
	"github.com/syncsynchalt/dime-a-tap/test"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/syncsynchalt/dime-a-tap/ca"
)

var testCACert string = `-----BEGIN CERTIFICATE-----
MIIDJTCCAg2gAwIBAgIBATANBgkqhkiG9w0BAQsFADA0MQswCQYDVQQGEwJVUzET
MBEGA1UEChMKRGltZS1BLVRhcDEQMA4GA1UECxMHRmFrZSBDQTAeFw0xODEwMDMx
OTU2MjVaFw0yODEwMDMxOTU2MjVaMDQxCzAJBgNVBAYTAlVTMRMwEQYDVQQKEwpE
aW1lLUEtVGFwMRAwDgYDVQQLEwdGYWtlIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEA1RyUGMGSuoXqzFkK9SDTFjiBB6aglWpF2h6nlWXZeSxG1eGk
RmnXR4VNCefsscKlav2eAn/pWMTGQ7PeQOUqPBq3sQy94dr7VNWt34uQ58Q2ffWc
4ay3znDeFCuPf/PiabFmk+ktm/AEGRLJHP9A15gUJ9bjEomNdyJI0EKvOXu2rkZT
+qBDxSwTVKdJ7IaahuYAB7om8lyEqqAevUhCMNRsTG4I9Ea4k4Dg8n+gY81ABJIL
fSqSWLWhlFUMvmc6PiX0SdF/mwcs+KZsOFoua2rYGrkyjMyPhfYjaTKXnbhGC9mU
6a1B7wm3vl+Bu6a/Lht7sI71VdnJ4I0IajSEjQIDAQABo0IwQDAOBgNVHQ8BAf8E
BAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUNVJTMZdEd9AHuWum+tvx
8ZgAFyowDQYJKoZIhvcNAQELBQADggEBAAvKkMyESOC41P7Bqb6af1xdiNfuA/+Y
XSRbHaaqBjz/3cqKogwpMTqUROdkOIw8QqF+BarjbfU2os2m5saCu4tW7VEk2IUX
zUY/BgTYOa4b+NvA4jdr37MtCZ3r6E+fE1i6S3cxy4mX9ZMY8Tm94FoOKc1qhqn1
yNNzGLEahQ8Qf2GqfNFcIg4MERzohPnP4X/JJD4PHQ2mXxPKoXSgP/eVIZt1h2vg
JjY2Ky2gFYTDwSCN0wctLRCPwHp7YMVFB/u79VeCg7MyHss4palpkh+YLPhgEFOi
0Sa1nCRk4/av/sn5HwzzRkjaZDNcPpHhcTvpsi1eYc05KFK25U0viuk=
-----END CERTIFICATE-----
`

var testCAKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA1RyUGMGSuoXqzFkK9SDTFjiBB6aglWpF2h6nlWXZeSxG1eGk
RmnXR4VNCefsscKlav2eAn/pWMTGQ7PeQOUqPBq3sQy94dr7VNWt34uQ58Q2ffWc
4ay3znDeFCuPf/PiabFmk+ktm/AEGRLJHP9A15gUJ9bjEomNdyJI0EKvOXu2rkZT
+qBDxSwTVKdJ7IaahuYAB7om8lyEqqAevUhCMNRsTG4I9Ea4k4Dg8n+gY81ABJIL
fSqSWLWhlFUMvmc6PiX0SdF/mwcs+KZsOFoua2rYGrkyjMyPhfYjaTKXnbhGC9mU
6a1B7wm3vl+Bu6a/Lht7sI71VdnJ4I0IajSEjQIDAQABAoIBAQC6fhrfmy4TCjQS
BW4AW2w9ys6nalqmxmxAV4khxRJN5sBKVP6UG/UngnCLVakdWg+2FCEdYOBMLU6v
Wo0JT0HpfRv41QSpzB8a+y8ALDtvhpaFHdXe622iO8Ur837NYxhkk7kHgQvHpX+A
jZ7vQDR3Nn+U6Yim5Tal5ZvAnEqIyrA6CwKW/8wsAl2E8CVvugjzQ2UKhfOUbarM
iQEMw4sBe1OL8beZ6EH9hWxxgvB3GLubVIAZlYUtjkVsh6m18pBl5XvTFXiJSgd4
Wdmzj75vHfa3b0G/Lppem5FHPnfgi6kWb7KM+IdMAZoTCoEbNzxqc7LGAexS3u6n
48R30bFlAoGBAPQ5XLaPitdBcWi03uGblO9DypJa5vPKkTUm5Sbv02gZvapsPPrG
9uru5yPVDJNSYTm7RKLf5YidxPiR61lDNGFRGpvuiRLjxYledmFhNzIXJujVpuYU
E+5Gbq7RLaFrmanGzQ7TL9UdIKra6ZRmwKYO//gyVtB3sIWdSgH5xl/XAoGBAN9j
LC8vqR+0EIrD3DvXDPB+US8Tckl+qXTiVM7l22lvJw3nys1SDaxl5F5AfPJ/7OqY
8/CiEMA1Jp0PlcL437+zBUsmebzAdUbVy4cQf4qVsW93d5yG3FwSl8ntkW9hvssc
nghpp06iPhWRT9U1+PWGSnHE0tsVRhvnAyN0kUI7AoGAfTMO6XQKzDD7b58Rh3zX
zBTnu0GoliApcqMe5Ggb64kOp1hXpoPrPyL8EW19xeR8fTkYhZrcM74VpQxBJ4CB
UMZgKsINOUbVFIf9jgxlXGNsCf7FUbvHP+aRhUMs7kyX+OY2ZzwykEEfZxdUmURX
zIlyBY3g3XwOXWD1+K9QV/8CgYBeZ1rU1h9y9nXHLt5zq34cZEWKz30M8ipK6xtM
FHeVJxQqHDroajS9FpJcAoTLNqS4v8rXdqX9lHitB1kS/HoSWWVzTN9FlU/6j39j
pOVBe+FwadxymcumXXUoMO21VGl9DKr8gynhYU87bh1+zUBZAleTnMo/K85lHEuH
QEvi4QKBgH5i2mftXBFhuTvR9X2+E0vO8cLbp1OGXVEfHAfVe4mxYblTj+VGjaAm
DL+c+tuNGLeYoe3RF0C3gZQqwiBZVK7mnjJU4nfJg8bEifOuzTeYe90wK1hKztxq
r05hVbaT7+DBl6NLVjl4iCBQnS5CTc0jggs0U0z6McyMB+p0ErwU
-----END RSA PRIVATE KEY-----
`

func TestStoreBadDomains(t *testing.T) {
	store, err := ca.NewStore("")
	test.Ok(t, err)
	_, _, err = store.GetCertificate("")
	test.Equals(t, "domain  is not valid", err.Error())

	_, _, err = store.GetCertificate(".a")
	test.Equals(t, "domain .a is not valid", err.Error())

	_, _, err = store.GetCertificate("a_z")
	test.Equals(t, "domain a_z is not valid", err.Error())
}

func looksLikeRsaKey(t *testing.T, pem []byte) {
	test.CallerDepth++
	defer func() { test.CallerDepth-- }()
	prefix := "-----BEGIN RSA PRIVATE KEY-----\nM"
	test.Equals(t, prefix, string(pem)[:len(prefix)])
}

func looksLikeCertificate(t *testing.T, pem []byte) {
	test.CallerDepth++
	defer func() { test.CallerDepth-- }()
	prefix := "-----BEGIN CERTIFICATE-----\nM"
	test.Equals(t, prefix, string(pem)[:len(prefix)])
}

func TestStoreMemCache(t *testing.T) {
	store, err := ca.NewStore("")
	test.Ok(t, err)
	key1, cert1, err := store.GetCertificate("a")
	test.Ok(t, err)
	looksLikeRsaKey(t, key1)
	looksLikeCertificate(t, cert1)

	key2, cert2, err := store.GetCertificate("a")
	test.Ok(t, err)
	test.Equals(t, key1, key2)
	test.Equals(t, cert1, cert2)
}

func TestStoreSmashCase(t *testing.T) {
	store, err := ca.NewStore("")
	test.Ok(t, err)
	key1, cert1, err := store.GetCertificate("a")
	test.Ok(t, err)

	key2, cert2, err := store.GetCertificate("A")
	test.Ok(t, err)

	test.Equals(t, key1, key2)
	test.Equals(t, cert1, cert2)
}

func TestStoreFileStore(t *testing.T) {
	mydir := fmt.Sprintf("/tmp/golang.test.%d", time.Now().UnixNano())
	defer os.RemoveAll(mydir)
	ca.CreateCAStore(mydir)

	store, err := ca.NewStore(mydir)
	test.Ok(t, err)
	key1, cert1, err := store.GetCertificate("a.com")
	test.Ok(t, err)

	key2, err := ioutil.ReadFile(mydir + "/domain-a.com.key")
	test.Ok(t, err)
	cert2, err := ioutil.ReadFile(mydir + "/domain-a.com.crt")
	test.Ok(t, err)
	test.Equals(t, key1, key2)
	test.Equals(t, cert1, cert2)

	key3, cert3, err := store.GetCertificate("a.com")
	test.Ok(t, err)
	test.Equals(t, key1, key3)
	test.Equals(t, cert1, cert3)

	key4, err := ioutil.ReadFile(mydir + "/domain-a.com.key")
	test.Ok(t, err)
	cert4, err := ioutil.ReadFile(mydir + "/domain-a.com.crt")
	test.Ok(t, err)
	test.Equals(t, key1, key4)
	test.Equals(t, cert1, cert4)
}
