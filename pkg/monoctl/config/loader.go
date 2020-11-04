package config

import (
	"io/ioutil"
	"os"
	"path"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/homedir"
)

const (
	RecommendedConfigPathEnvVar = "MONOSKOPECONFIG"
	RecommendedHomeDir          = ".monoskope"
	RecommendedFileName         = "config"
)

var (
	RecommendedConfigDir = path.Join(homedir.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile  = path.Join(RecommendedConfigDir, RecommendedFileName)
)

type ClientConfigLoader struct {
	// Logger interface
	log          logger.Logger
	config       *Config
	configPath   string
	ExplicitFile string
}

// NewLoader is a convenience function that returns a new ClientConfigLoader object with defaults
func NewLoader() *ClientConfigLoader {
	return &ClientConfigLoader{
		log: logger.WithName("client-config-loader"),
	}
}

func (l *ClientConfigLoader) InitConifg(config *Config) error {
	if l.ExplicitFile != "" {
		exists, err := util.FileExists(l.ExplicitFile)
		if err != nil {
			return err
		}
		if exists {
			return ErrAlreadyInitialized
		}
		l.config = config
		l.configPath = l.ExplicitFile
		return l.SaveToFile(config, l.configPath, 0644)
	}

	exists, err := util.FileExists(RecommendedHomeFile)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyInitialized
	}
	err = util.CreateDir(RecommendedConfigDir, 0700)
	if err != nil {
		return err
	}

	l.config = config
	l.configPath = RecommendedHomeFile
	return l.SaveToFile(config, l.configPath, 0600)
}

// LoadAndStoreConfig loads and stores the config either from env or home file.
func (l *ClientConfigLoader) LoadAndStoreConfig() error {
	if l.ExplicitFile != "" {
		if err := l.loadAndStoreConfig(l.ExplicitFile); err != nil {
			return err
		}
		l.configPath = l.ExplicitFile
		return nil
	}

	// Load config from envvar path if provided
	envVarFile := os.Getenv(RecommendedConfigPathEnvVar)
	if len(envVarFile) != 0 {
		if err := l.loadAndStoreConfig(envVarFile); err != nil {
			return err
		}
		l.configPath = envVarFile
		return nil
	}

	// Load recommended home file if present
	if err := l.loadAndStoreConfig(RecommendedHomeFile); err != nil {
		if os.IsNotExist(err) {
			return ErrNoConfigExists
		}
		return err
	}
	l.configPath = RecommendedHomeFile

	return nil
}

// loadAndStoreConfig checks if the given file exists and loads it's contents
func (l *ClientConfigLoader) loadAndStoreConfig(filename string) error {
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return err
	}
	l.config, err = l.LoadFromFile(filename)
	return err
}

// GetConfigPath returns the path of the previously loaded config
func (l *ClientConfigLoader) GetConfigPath() string {
	return l.configPath
}

// GetConfig returns the previously loaded config
func (l *ClientConfigLoader) GetConfig() *Config {
	return l.config
}

// LoadFromFile takes a filename and deserializes the contents into Config object
func (l *ClientConfigLoader) LoadFromFile(filename string) (*Config, error) {
	monoconfigBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := l.LoadFromBytes(monoconfigBytes)
	if err != nil {
		return nil, err
	}
	l.log.Info("Config loaded from file", "filename", filename)

	return config, nil
}

// Load takes a byte slice and deserializes the contents into Config object.
// Encapsulates deserialization without assuming the source is a file.
func (*ClientConfigLoader) LoadFromBytes(data []byte) (*Config, error) {
	config := NewConfig()

	err := yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (l *ClientConfigLoader) SaveToFile(config *Config, filename string, permission os.FileMode) error {
	bytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, permission)
	if err != nil {
		return err
	}
	l.log.Info("Config saved to file", "filename", filename)

	return nil
}
