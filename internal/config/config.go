// The config package is responsible for loading configuation files.

package config

import (
	"fmt"
	// embed required for sample configs below.
	_ "embed"

	"github.com/jwalton/kitsch-prompt/internal/modules"
	"github.com/jwalton/kitsch-prompt/sampleconfig"
	"gopkg.in/yaml.v3"
)

// Config represents a configuration file.
type Config struct {
	Styles map[string]yaml.Node `yaml:"styles"`
	// Prompt is the module to use to display the prompt.
	Prompt yaml.Node
}

// GetPromptModule returns the root prompt module for the configuration.
func (c *Config) GetPromptModule() (modules.Module, error) {
	if c == nil || c.Prompt.IsZero() {
		return nil, fmt.Errorf("configuration is missing prompt")
	}
	return modules.CreateModule(&c.Prompt)
}

// LoadFromYaml loads the configuration file from a YAML file.
func (c *Config) LoadFromYaml(yamlData []byte) error {
	return yaml.Unmarshal(yamlData, &c)
}

// LoadDefaultConfig will load a default configuration.
func LoadDefaultConfig() (Config, error) {
	var config = Config{}
	err := config.LoadFromYaml(sampleconfig.DefaultConfig)
	if err != nil {
		// Default config should not have errors!
		fmt.Println("kitch: Error in default configuration", err)
		return Config{}, err
	}
	return config, nil
}
