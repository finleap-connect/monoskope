package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const (
	AuthTypeBasic               = "basic"
	AuthTypeBasicSuffixUsername = ".basic.username"
	AuthTypeBasicSuffixPassword = ".basic.password"

	AuthTypeSSH                 = "ssh"
	AuthTypeSSHSuffixPrivateKey = ".ssh.privateKey"
	AuthTypeSSHSuffixPassword   = ".ssh.password"
	AuthTypeSSHSuffixKnownHosts = ".ssh.known_hosts"
)

type GitConfig struct {
	// The (possibly remote) repository URL to clone from.
	URL string `yaml:"url"`
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod `yaml:"-"`
	// Fetch only ReferenceName if true.
	SingleBranch bool `yaml:"singleBranch"`
	// Limit fetching to the specified number of commits.
	Depth int `yaml:"depth"`
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool `yaml:"insecureSkipTLS"`
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte `yaml:"caBundle"`
	// Remote branch to clone. If empty, uses HEAD.
	ReferenceName plumbing.ReferenceName `yaml:"referenceName"`
	// AuthType is the type of authentication
	AuthType string `yaml:"authType"`
	// EnvPrefix is the prefix for environment variables to configure
	EnvPrefix string `yaml:"envPrefix"`

	authMethod   transport.AuthMethod
	cloneOptions *git.CloneOptions
	pullOptions  *git.PullOptions
	pushOptions  *git.PushOptions
}

// configureBasicAuth reads the file containing the basic auth information and unmarshal's it's content into the clone options given.
func (c *GitConfig) configureBasicAuth() error {
	// get env
	usernameKey := fmt.Sprintf("%s%s", c.EnvPrefix, AuthTypeBasicSuffixUsername)
	passwordKey := fmt.Sprintf("%s%s", c.EnvPrefix, AuthTypeBasicSuffixPassword)

	username := os.Getenv(usernameKey)
	if username == "" {
		return fmt.Errorf("%s must not be empty", username)
	}

	password := os.Getenv(passwordKey)
	if password == "" {
		return fmt.Errorf("%s must not be empty", passwordKey)
	}

	// set clone options auth
	c.authMethod = &http.BasicAuth{
		Username: username,
		Password: password,
	}

	return nil
}

// configureSSHAuth reads the file containing the ssh auth information and unmarshal's it's content into the clone options given.
func (c *GitConfig) configureSSHAuth() error {
	privateKeyEnvKey := fmt.Sprintf("%s%s", c.EnvPrefix, AuthTypeSSHSuffixPrivateKey)
	privateKey := os.Getenv(privateKeyEnvKey)
	if privateKey == "" {
		return fmt.Errorf("%s must not be empty", privateKeyEnvKey)
	}

	knownHostsEnvKey := fmt.Sprintf("%s%s", c.EnvPrefix, AuthTypeSSHSuffixKnownHosts)
	knownHosts := os.Getenv(knownHostsEnvKey)
	if knownHosts == "" {
		return fmt.Errorf("%s must not be empty", knownHostsEnvKey)
	}

	// password is optional
	passwordKey := fmt.Sprintf("%s%s", c.EnvPrefix, AuthTypeSSHSuffixPassword)
	password := os.Getenv(passwordKey)

	f, err := os.CreateTemp("", "known-hosts")
	if err != nil {
		return err
	}

	_, err = f.WriteString(knownHosts)
	if err != nil {
		return err
	}

	// configure public key ssh auth
	callback, err := ssh.NewKnownHostsCallback(f.Name())
	if err != nil {
		return err
	}
	publicKeys, err := ssh.NewPublicKeys(ssh.DefaultUsername, []byte(privateKey), password)
	if err != nil {
		return err
	}
	publicKeys.HostKeyCallback = callback
	c.authMethod = publicKeys

	return nil
}

func (c *GitConfig) getAuthMethod() (transport.AuthMethod, error) {
	if c.authMethod != nil {
		return c.authMethod, nil
	}

	// Configure basic auth optionally
	if c.AuthType == AuthTypeBasic {
		if err := c.configureBasicAuth(); err != nil {
			return nil, err
		}
	}

	// Configure ssh auth
	if c.AuthType == AuthTypeSSH {
		if err := c.configureSSHAuth(); err != nil {
			return nil, err
		}
	}

	return c.authMethod, nil
}

func (c *GitConfig) getCloneOptions() (*git.CloneOptions, error) {
	if c.cloneOptions != nil {
		return c.cloneOptions, nil
	}

	authMethod, err := c.getAuthMethod()
	if err != nil {
		return nil, err
	}

	c.cloneOptions = &git.CloneOptions{
		URL:             c.URL,
		SingleBranch:    c.SingleBranch,
		ReferenceName:   c.ReferenceName,
		Depth:           c.Depth,
		Progress:        os.Stdout,
		CABundle:        c.CABundle,
		InsecureSkipTLS: c.InsecureSkipTLS,
		Auth:            authMethod,
	}
	if err := c.cloneOptions.Validate(); err != nil {
		return nil, err
	}
	return c.cloneOptions, nil
}

func (c *GitConfig) getPullOptions() (*git.PullOptions, error) {
	if c.pullOptions != nil {
		return c.pullOptions, nil
	}

	authMethod, err := c.getAuthMethod()
	if err != nil {
		return nil, err
	}

	c.pullOptions = &git.PullOptions{
		SingleBranch:    c.SingleBranch,
		ReferenceName:   c.ReferenceName,
		Depth:           c.Depth,
		Progress:        os.Stdout,
		CABundle:        c.CABundle,
		InsecureSkipTLS: c.InsecureSkipTLS,
		Auth:            authMethod,
	}
	if err := c.pullOptions.Validate(); err != nil {
		return nil, err
	}
	return c.pullOptions, nil
}

func (c *GitConfig) getPushOptions() (*git.PushOptions, error) {
	if c.pushOptions != nil {
		return c.pushOptions, nil
	}

	authMethod, err := c.getAuthMethod()
	if err != nil {
		return nil, err
	}

	c.pushOptions = &git.PushOptions{
		Progress:        os.Stdout,
		CABundle:        c.CABundle,
		InsecureSkipTLS: c.InsecureSkipTLS,
		Auth:            authMethod,
	}
	if err := c.pushOptions.Validate(); err != nil {
		return nil, err
	}
	return c.pushOptions, nil
}
