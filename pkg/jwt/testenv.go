package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	env.privateKey = privKey

	privKeyFile, err := ioutil.TempFile("", "private.key")
	if err != nil {
		return nil, err
	}
	defer privKeyFile.Close()

	_, err = privKeyFile.Write(x509.MarshalPKCS1PrivateKey(privKey))
	if err != nil {
		return nil, err
	}
	env.privateKeyFile = privKeyFile.Name()

	pubKeyFile, err := ioutil.TempFile("", "public.key")
	if err != nil {
		return nil, err
	}
	defer pubKeyFile.Close()
	env.publicKeyFile = pubKeyFile.Name()

	pubKeyPem, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}

	_, err = pubKeyFile.Write(pubKeyPem)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (env *TestEnv) Shutdown() error {
	return nil
}
