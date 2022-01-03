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
	mu         sync.RWMutex
	caCertFile string
	certFile   string
	keyFile    string
	keyPair    *tls.Certificate
	rootCAs    *x509.CertPool
	watcher    *fsnotify.Watcher
	watching   chan struct{}
	log        logger.Logger
}

func NewTLSConfigLoader(caCertFile, certFile, keyFile string) (*TLSConfigLoader, error) {
	if caCertFile == "" {
		return nil, errors.New("caCertFile must not be empty")
	}
	if certFile == "" {
		return nil, errors.New("certFile must not be empty")
	}
	if keyFile == "" {
		return nil, errors.New("keyFile must not be empty")
	}
	var err error
	caCertFile, err = filepath.Abs(caCertFile)
	if err != nil {
		return nil, err
	}
	certFile, err = filepath.Abs(certFile)
	if err != nil {
		return nil, err
	}
	keyFile, err = filepath.Abs(keyFile)
	if err != nil {
		return nil, err
	}

	return &TLSConfigLoader{
		mu:         sync.RWMutex{},
		caCertFile: caCertFile,
		certFile:   certFile,
		keyFile:    keyFile,
		log:        logger.WithName("tls-config-loader"),
	}, nil
}

// Watch starts watching for changes to the certificate
// and key files. On any change the certificate and key
// are reloaded. If there is an issue the load will fail
// and the old (if any) certificates and keys will continue
// to be used.
func (t *TLSConfigLoader) Watch() error {
	var err error
	if t.watcher, err = fsnotify.NewWatcher(); err != nil {
		return errors.Wrap(err, "can't create watcher")
	}
	if err = t.watcher.Add(t.caCertFile); err != nil {
		return errors.Wrap(err, "can't watch ca cert file")
	}
	if err = t.watcher.Add(t.certFile); err != nil {
		return errors.Wrap(err, "can't watch cert file")
	}
	if err = t.watcher.Add(t.keyFile); err != nil {
		return errors.Wrap(err, "can't watch key file")
	}
	if err := t.load(); err != nil {
		t.log.Error(err, "can't load")
	}
	t.log.Info("watching for ca, cert and key change", "caCertFile", t.caCertFile, "certFile", t.certFile, "keyFile", t.keyFile)
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
	certs, err := ioutil.ReadFile(t.caCertFile)
	if err != nil {
		return err
	}

	// Append local CA cert to the system pool
	if rootCAs.AppendCertsFromPEM(certs) {
		t.log.Info("root CAs loaded")
	} else {
		t.log.Info("No root CAs appended, using system CAs only")
	}

	t.rootCAs = rootCAs

	keyPair, err := tls.LoadX509KeyPair(t.certFile, t.keyFile)
	if err == nil {
		t.keyPair = &keyPair
		t.log.Info("certificate and key loaded")
	}

	return err
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
func (t *TLSConfigLoader) GetTLSConfig() *tls.Config {
	return &tls.Config{
		RootCAs:              t.GetRootCAs(),
		GetCertificate:       t.GetCertificate,
		GetClientCertificate: t.GetClientCertificate,
	}
}

// GetRootCAs returns the cert pool to use to verify certificates.
func (t *TLSConfigLoader) GetRootCAs() *x509.CertPool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.rootCAs
}

// GetCertificate returns the loaded certificate for use by
// the TLSConfig fields GetCertificate field in a http.Server.
func (t *TLSConfigLoader) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.keyPair, nil
}

// GetClientCertificate returns the loaded certificate for use by
// the TLSConfig fields GetClientCertificate field in a http.Server.
func (t *TLSConfigLoader) GetClientCertificate(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.keyPair, nil
}

// Stop stops watching for changes to the
// certificate and key files.
func (t *TLSConfigLoader) Stop() {
	t.watching <- struct{}{}
}
