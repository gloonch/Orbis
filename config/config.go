// config/config.go
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration loaded from config/config.yaml
type Config struct {
	Ephemeris struct {
		FilePath      string `yaml:"filePath"`
		HighPrecision bool   `yaml:"highPrecision"`
	} `yaml:"ephemeris"`

	Stream struct {
		IntervalSeconds int `yaml:"intervalSeconds"`
	} `yaml:"stream"`

	Kafka struct {
		Brokers        []string `yaml:"brokers"`
		TopicPositions string   `yaml:"topicPositions"`
		GroupID        string   `yaml:"groupID"`
	} `yaml:"kafka"`

	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

// Load reads the YAML config at the given path and unmarshals it into Config.
func Load(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
