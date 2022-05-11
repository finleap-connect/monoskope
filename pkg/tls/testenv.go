// Copyright 2022 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/fs"
	"io/ioutil"
	"math/big"
	"net"
	"time"

	"github.com/finleap-connect/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv

	ca         *x509.Certificate
	caPrivKey  *rsa.PrivateKey
	caCertFile string

	cert        tls.Certificate
	certFile    string
	certKeyFile string
}

func NewTestEnv(testEnv *test.TestEnv) (*TestEnv, error) {
	env := &TestEnv{
		TestEnv: testEnv,
	}

	if err := env.CreateCACertificate(); err != nil {
		return nil, err
	}

	return env, nil
}

func (t *TestEnv) CreateCACertificate() error {
	t.ca = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	t.caPrivKey = caPrivKey

	caBytes, err := x509.CreateCertificate(rand.Reader, t.ca, t.ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return err
	}

	caCertFile, err := ioutil.TempFile("", "ca.crt")
	if err != nil {
		return err
	}
	defer caCertFile.Close()
	t.caCertFile = caCertFile.Name()

	err = ioutil.WriteFile(t.caCertFile, caPEM.Bytes(), fs.ModeAppend)
	if err != nil {
		return err
	}

	return nil
}

func (env *TestEnv) Shutdown() error {
	return nil
}

func (t *TestEnv) CreateCertificate() error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, t.ca, &certPrivKey.PublicKey, t.caPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return err
	}

	certFile, err := ioutil.TempFile("", "cert.crt")
	if err != nil {
		return err
	}
	defer certFile.Close()
	t.certFile = certFile.Name()

	err = ioutil.WriteFile(t.certFile, certPEM.Bytes(), fs.ModeAppend)
	if err != nil {
		return err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return err
	}

	certKeyFile, err := ioutil.TempFile("", "cert.key")
	if err != nil {
		return err
	}
	defer certKeyFile.Close()
	t.certKeyFile = certKeyFile.Name()

	err = ioutil.WriteFile(t.certKeyFile, certPrivKeyPEM.Bytes(), fs.ModeAppend)
	if err != nil {
		return err
	}

	certPair, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return err
	}
	t.cert = certPair

	return nil
}
