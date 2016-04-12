package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DatabaseConfiguration struct {
	Hostname string
	Database string
	Username string
	Password string
}

type Configuration struct {
	Platforms         []string
	DefaultPlatform   string   `yaml:"default-platform"`
	CriticalPlatforms []string `yaml:"critical-platforms"`
	SqlDir            string   `yaml:"sql-dir"`
	DatabaseCharset   string   `yaml:"database-charset"`
	Database          map[string]DatabaseConfiguration
}

func ParseConfiguration(source []byte) (*Configuration, error) {
	var configuration Configuration
	err := yaml.Unmarshal(source, &configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func LoadConfiguration(file string) (*Configuration, error) {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	configuration, err := ParseConfiguration(source)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}
