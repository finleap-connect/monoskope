package config

import (
	"io/ioutil"
	"os"
	"path"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"gopkg.in/yaml.v2"
)

const (
	RecommendedConfigPathEnvVar = "MONOSKOPECONFIG"
	RecommendedHomeDir          = ".monoskope"
	RecommendedFileName         = "config"
	FileMode                    = 0644
)

var (
	RecommendedConfigDir = path.Join(util.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile  = path.Join(RecommendedConfigDir, RecommendedFileName)
)

type ClientConfigManager struct {
	// Logger interface
	log          logger.Logger
	config       *Config
	configPath   string
	explicitFile string
}

// NewLoader is a convenience function that returns a new ClientConfigManager object with defaults
func NewLoader() *ClientConfigManager {
	return &ClientConfigManager{
		log: logger.WithName("client-config-loader"),
	}
}

// NewLoaderFromExplicitFile is a convenience function that returns a new ClientConfigManager object with explicitFile set
func NewLoaderFromExplicitFile(explicitFile string) *ClientConfigManager {
	loader := NewLoader()
	loader.explicitFile = explicitFile
	return loader
}

// loadAndStoreConfig checks if the given file exists and loads it's contents
func (l *ClientConfigManager) saveAndStoreConfig(filename string, config *Config) error {
	exists, err := util.FileExists(filename)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyInitialized
	}
	l.config = config
	l.configPath = filename
	return l.SaveToFile(config, l.configPath, FileMode)
}

// loadAndStoreConfig checks if the given file exists and loads it's contents
func (l *ClientConfigManager) loadAndStoreConfig(filename string) error {
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return err
	}
	l.config, err = l.LoadFromFile(filename)
	return err
}

// GetConfigPath returns the path of the previously loaded config
func (l *ClientConfigManager) GetConfigPath() string {
	return l.configPath
}

// GetConfig returns the previously loaded config
func (l *ClientConfigManager) GetConfig() *Config {
	return l.config
}

func (l *ClientConfigManager) InitConifg(config *Config) error {
	if l.explicitFile != "" {
		return l.saveAndStoreConfig(l.explicitFile, config)
	}

	envVarFile := os.Getenv(RecommendedConfigPathEnvVar)
	if len(envVarFile) != 0 {
		return l.saveAndStoreConfig(envVarFile, config)
	}

	return l.saveAndStoreConfig(RecommendedHomeFile, config)
}

func (l *ClientConfigManager) SaveConfig() error {
	if l.configPath == "" || l.config == nil {
		return ErrNoConfigExists
	}
	return l.SaveToFile(l.config, l.configPath, FileMode)
}

// LoadAndStoreConfig loads and stores the config either from env or home file.
func (l *ClientConfigManager) LoadAndStoreConfig() error {
	if l.explicitFile != "" {
		if err := l.loadAndStoreConfig(l.explicitFile); err != nil {
			return err
		}
		l.configPath = l.explicitFile
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

// LoadFromFile takes a filename and deserializes the contents into Config object
func (l *ClientConfigManager) LoadFromFile(filename string) (*Config, error) {
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

// LoadFromBytes takes a byte slice and deserializes the contents into Config object.
// Encapsulates deserialization without assuming the source is a file.
func (*ClientConfigManager) LoadFromBytes(data []byte) (*Config, error) {
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

// SaveToFile takes a config, serializes the contents and stores them into a file.
func (l *ClientConfigManager) SaveToFile(config *Config, filename string, permission os.FileMode) error {
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
