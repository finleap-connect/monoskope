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

package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/fs"
	"os"

	"github.com/finleap-connect/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv

	privateKeyFile string
	publicKeyFile  string

	privateKey *rsa.PrivateKey
}

func NewTestEnv(testEnv *test.TestEnv) (*TestEnv, error) {
	env := &TestEnv{
		TestEnv: testEnv,
	}

	privKeyFile, err := os.CreateTemp("", "private.key")
	if err != nil {
		return env, err
	}
	defer privKeyFile.Close()
	env.privateKeyFile = privKeyFile.Name()

	pubKeyFile, err := os.CreateTemp("", "public.key")
	if err != nil {
		return nil, err
	}
	defer pubKeyFile.Close()
	env.publicKeyFile = pubKeyFile.Name()

	err = env.RotateCertificate()
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) RotateCertificate() error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	env.privateKey = privKey

	err = os.WriteFile(env.privateKeyFile, x509.MarshalPKCS1PrivateKey(privKey), fs.ModeAppend)
	if err != nil {
		return err
	}

	pubKeyPem, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	err = os.WriteFile(env.publicKeyFile, pubKeyPem, fs.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

func (env *TestEnv) CreateSigner() JWTSigner {
	return NewSigner(env.privateKeyFile)
}

func (env *TestEnv) CreateVerifier() (JWTVerifier, error) {
	return NewVerifier(env.publicKeyFile)
}

func (env *TestEnv) Shutdown() error {
	return nil
}
