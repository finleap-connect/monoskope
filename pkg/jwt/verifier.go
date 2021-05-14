package jwt

import (
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Verifier verifies a JWT and parses claims
type Verifier interface {
	Verify(string, interface{}) error
	Close() error
}

type jwtVerifier struct {
	log       logger.Logger
	publicKey []interface{}
	watcher   *fsnotify.Watcher
}

// NewVerifier creates a new verifier for raw JWT
func NewVerifier(publicKeyFilename string) (Verifier, error) {
	v := &jwtVerifier{
		log:       logger.WithName("jwt-verifier"),
		publicKey: make([]interface{}, 0),
	}

	v.log.Info("Loading public key...", "publicKeyFilename", publicKeyFilename)
	err := v.rotatePublicKey(publicKeyFilename)
	if err != nil {
		return nil, err
	}

	v.log.Info("Setting up watcher...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go v.loadPublicKeyOnFileChange()

	err = watcher.Add(publicKeyFilename)
	if err != nil {
		return nil, err
	}
	v.watcher = watcher

	return v, nil
}

func (v *jwtVerifier) loadPublicKeyOnFileChange() {
	for {
		select {
		case event, ok := <-v.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				v.log.Info("Public key has been changed. Updating...")
				err := v.rotatePublicKey(event.Name)
				if err != nil {
					v.log.Error(err, "Error rotating public key.")
				}
				v.log.Info("Public key has been updated.", "KeyCount", len(v.publicKey))
			}
		case err, ok := <-v.watcher.Errors:
			if !ok {
				return
			}
			v.log.Error(err, "Error from watcher.")
		}
	}
}

// loadPublicKey loads the public key
func (v *jwtVerifier) rotatePublicKey(filename string) error {
	pubKeyBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	pubKey, err := LoadPublicKey(pubKeyBytes)
	if err != nil {
		return err
	}

	v.publicKey = append(v.publicKey, pubKey)

	return nil
}

// Verify parses the raw JWT, verifies the content against the public key of the verifier and parses the claims
func (v *jwtVerifier) Verify(rawJWT string, claims interface{}) error {
	parsedJWT, err := jwt.ParseSigned(rawJWT)
	if err != nil {
		return err
	}

	for _, pubKey := range v.publicKey {
		if err := parsedJWT.Claims(pubKey, claims); err != nil {
			return err
		}
	}

	return nil
}

// Close closes file watcher
func (v *jwtVerifier) Close() error {
	return v.watcher.Close()
}
