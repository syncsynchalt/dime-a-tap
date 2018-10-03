package ca

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Store struct {
	lock   sync.Mutex
	keys   map[string][]byte
	certs  map[string][]byte
	dir    string
	cakey  []byte
	cacert []byte
}

func NewStore(directory string) (*Store, error) {
	var cakey, cacert []byte
	var err1, err2 error
	if directory == "" {
		cakey, err1 = GenerateCAKey()
		cacert, err2 = GenerateCACert(cakey)
	} else {
		cakey, err1 = ioutil.ReadFile(directory + "/ca.key")
		cacert, err2 = ioutil.ReadFile(directory + "/ca.crt")
	}
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	return &Store{
		keys:   make(map[string][]byte),
		certs:  make(map[string][]byte),
		dir:    directory,
		cakey:  cakey,
		cacert: cacert,
	}, nil
}

func (store *Store) keyfile(domain string) string {
	if store.dir != "" {
		return fmt.Sprintf("%s/domain-%s.key", store.dir, domain)
	} else {
		return ""
	}
}

func (store *Store) certfile(domain string) string {
	if store.dir != "" {
		return fmt.Sprintf("%s/domain-%s.crt", store.dir, domain)
	} else {
		return ""
	}
}

func (store *Store) writeFile(filename string, data []byte, perm os.FileMode) error {
	if store.dir != "" {
		return writeFileExcl(filename, data, perm)
	} else {
		return nil
	}
}

func (store *Store) exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func (store *Store) GetCertificate(domain string) (key, cert []byte, err error) {
	domain = strings.ToLower(domain)
	if !isSafeDomain(domain) {
		return nil, nil, fmt.Errorf("domain %s is not valid", domain)
	}

	store.lock.Lock()
	defer store.lock.Unlock()

	if store.keys[domain] == nil {
		filename := store.keyfile(domain)
		if store.exists(filename) {
			store.keys[domain], err = ioutil.ReadFile(filename)
			if err != nil {
				return nil, nil, err
			}
		} else {
			store.keys[domain], err = generateKey()
			if err != nil {
				return nil, nil, err
			}
			err = store.writeFile(filename, store.keys[domain], 0600)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	if store.certs[domain] == nil {
		filename := store.certfile(domain)
		if store.exists(filename) {
			store.certs[domain], err = ioutil.ReadFile(filename)
			if err != nil {
				return nil, nil, err
			}
		} else {
			store.certs[domain], err = generateCert(domain, store.keys[domain], store.cakey, store.cacert)
			if err != nil {
				return nil, nil, err
			}
			err = store.writeFile(filename, store.certs[domain], 0644)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return store.keys[domain], store.certs[domain], nil
}

func isSafeDomain(domain string) bool {
	if len(domain) == 0 {
		return false
	}
	if domain[0] == '.' {
		return false
	}
	return regexp.MustCompile(`^[-\.a-z0-9]+$`).MatchString(domain)
}
