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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/pkg/errors"
	"gopkg.in/fsnotify.v1"
)

type TLSConfigLoader struct {
	mu                       sync.RWMutex
	serverCACertificateFile  string
	serverCertificateFile    string
	serverCertificateKeyFile string
	serverCertificate        *tls.Certificate
	clientCACertificateFile  string
	clientCertificateFile    string
	clientCertificateKeyFile string
	clientCertificate        *tls.Certificate
	serverCAs                *x509.CertPool
	clientCAs                *x509.CertPool
	watcher                  *fsnotify.Watcher
	watching                 chan struct{}
	log                      logger.Logger
}

func NewTLSConfigLoader() (*TLSConfigLoader, error) {
	return &TLSConfigLoader{
		mu:  sync.RWMutex{},
		log: logger.WithName("tls-config-loader"),
	}, nil
}

func (t *TLSConfigLoader) SetServerCACertificate(caCertificateFile string) error {
	if caCertificateFile == "" {
		return errors.New("caCertificateFile must not be empty")
	}
	var err error
	caCertificateFile, err = filepath.Abs(caCertificateFile)
	if err != nil {
		return err
	}
	t.serverCACertificateFile = caCertificateFile
	return nil
}

func (t *TLSConfigLoader) SetClientCACertificate(caCertificateFile string) error {
	if caCertificateFile == "" {
		return errors.New("caCertificateFile must not be empty")
	}
	var err error
	caCertificateFile, err = filepath.Abs(caCertificateFile)
	if err != nil {
		return err
	}
	t.clientCACertificateFile = caCertificateFile
	return nil
}

func (t *TLSConfigLoader) SetClientCertificate(certificateFile, keyFile string) error {
	if certificateFile == "" {
		return errors.New("certificateFile must not be empty")
	}
	if keyFile == "" {
		return errors.New("keyFile must not be empty")
	}

	var err error
	certificateFile, err = filepath.Abs(certificateFile)
	if err != nil {
		return err
	}
	t.clientCertificateFile = certificateFile

	keyFile, err = filepath.Abs(keyFile)
	if err != nil {
		return err
	}
	t.clientCertificateKeyFile = keyFile
	return nil
}

func (t *TLSConfigLoader) SetServerCertificate(certificateFile, keyFile string) error {
	if certificateFile == "" {
		return errors.New("certificateFile must not be empty")
	}
	if keyFile == "" {
		return errors.New("keyFile must not be empty")
	}

	var err error
	certificateFile, err = filepath.Abs(certificateFile)
	if err != nil {
		return err
	}
	t.serverCertificateFile = certificateFile

	keyFile, err = filepath.Abs(keyFile)
	if err != nil {
		return err
	}
	t.serverCertificateKeyFile = keyFile
	return nil
}

// On any change the certificate and key
// are reloaded. If there is an issue the load will fail
// and the old (if any) certificates and keys will continue
// to be used.
func (t *TLSConfigLoader) Watch() error {
	var err error
	if t.watcher, err = fsnotify.NewWatcher(); err != nil {
		return errors.Wrap(err, "can't create watcher")
	}

	if t.serverCACertificateFile != "" {
		if err = t.watcher.Add(t.serverCACertificateFile); err != nil {
			return errors.Wrap(err, "can't watch ca cert file")
		}
	}

	if t.clientCACertificateFile != "" {
		if err = t.watcher.Add(t.clientCACertificateFile); err != nil {
			return errors.Wrap(err, "can't watch ca cert file")
		}
	}

	if t.clientCertificateFile != "" {
		if err = t.watcher.Add(t.clientCertificateFile); err != nil {
			return errors.Wrap(err, "can't watch cert file")
		}
	}

	if t.serverCertificateFile != "" {
		if err = t.watcher.Add(t.serverCertificateFile); err != nil {
			return errors.Wrap(err, "can't watch cert file")
		}
	}

	if t.clientCertificateKeyFile != "" {
		if err = t.watcher.Add(t.clientCertificateKeyFile); err != nil {
			return errors.Wrap(err, "can't watch key file")
		}
	}

	if t.serverCertificateKeyFile != "" {
		if err = t.watcher.Add(t.serverCertificateKeyFile); err != nil {
			return errors.Wrap(err, "can't watch key file")
		}
	}

	if err := t.load(); err != nil {
		t.log.Error(err, "couldn't load")
	}
	t.log.Info("watching for changes")
	t.watching = make(chan struct{})
	go t.run()
	return nil
}

func (t *TLSConfigLoader) load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Load local CA
	if t.serverCACertificateFile != "" {
		certs, err := ioutil.ReadFile(t.serverCACertificateFile)
		if err != nil {
			return err
		}
		// Append local CA cert to the system pool
		if rootCAs.AppendCertsFromPEM(certs) {
			t.log.Info("root CAs loaded")
		} else {
			t.log.Info("No root CAs appended, using system CAs only")
		}
		t.serverCAs = rootCAs
	}

	if t.clientCertificateFile != "" {
		clientCert, err := tls.LoadX509KeyPair(t.clientCertificateFile, t.clientCertificateKeyFile)
		if err == nil {
			t.clientCertificate = &clientCert
			t.log.Info("client certificate and key loaded")
		} else {
			return err
		}
	}

	if t.serverCertificateFile != "" {
		serverCert, err := tls.LoadX509KeyPair(t.serverCertificateFile, t.serverCertificateKeyFile)
		if err == nil {
			t.serverCertificate = &serverCert
			t.log.Info("server certificate and key loaded")
		} else {
			return err
		}
	}

	return nil
}

func (t *TLSConfigLoader) run() {
loop:
	for {
		select {
		case <-t.watching:
			break loop
		case event := <-t.watcher.Events:
			t.log.V(logger.DebugLevel).Info("watch event", "event", event)
			if err := t.load(); err != nil {
				t.log.Error(err, "can't load")
			}
		case err := <-t.watcher.Errors:
			t.log.Error(err, "error watching files")
		}
	}
	t.log.Info("stopped watching")
	t.watcher.Close()
}

// GetTLSConfig returns a tls.Config with auto reloading certs.
func (t *TLSConfigLoader) GetClientTLSConfig() *tls.Config {
	return &tls.Config{
		RootCAs:              t.GetRootCAs(),
		GetClientCertificate: t.GetClientCertificate,
	}
}

func (t *TLSConfigLoader) GetServerTLSConfig(clientAuthType tls.ClientAuthType) *tls.Config {
	conf := &tls.Config{
		GetCertificate: t.GetCertificate,
		ClientAuth:     clientAuthType,
	}
	if clientAuthType == tls.RequireAndVerifyClientCert {
		conf.ClientCAs = t.GetClientCAs()
	}
	return conf
}

// GetRootCAs returns the cert pool to use to verify certificates.
func (t *TLSConfigLoader) GetRootCAs() *x509.CertPool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.serverCAs
}

// GetClientCAs returns the cert pool to use to verify client certificates.
func (t *TLSConfigLoader) GetClientCAs() *x509.CertPool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.clientCAs
}

// GetCertificate returns the loaded certificate for use by
// the TLSConfig fields GetCertificate field in a http.Server.
func (t *TLSConfigLoader) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.serverCertificate, nil
}

// GetClientCertificate returns the loaded certificate for use by
// the TLSConfig fields GetClientCertificate field in a http.Server.
func (t *TLSConfigLoader) GetClientCertificate(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.clientCertificate, nil
}

// Stop stops watching for changes to the
// certificate and key files.
func (t *TLSConfigLoader) Stop() {
	t.watching <- struct{}{}
}
