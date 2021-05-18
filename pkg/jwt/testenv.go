package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/fs"
	"io/ioutil"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
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

	privKeyFile, err := ioutil.TempFile("", "private.key")
	if err != nil {
		return env, err
	}
	defer privKeyFile.Close()
	env.privateKeyFile = privKeyFile.Name()

	pubKeyFile, err := ioutil.TempFile("", "public.key")
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

	err = ioutil.WriteFile(env.privateKeyFile, x509.MarshalPKCS1PrivateKey(privKey), fs.ModeAppend)
	if err != nil {
		return err
	}

	pubKeyPem, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(env.publicKeyFile, pubKeyPem, fs.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

func (env *TestEnv) Shutdown() error {
	return nil
}
