package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
)

// configuration
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

// version as an array of integers
type Version []int

func NewVersion(v string) (Version, error) {
	version := []int{}
	for _, s := range strings.Split(v, ".") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version '%s'", v)
		}
		version = append(version, i)
	}
	return version, nil
}

func (v1 Version) CompareTo(v2 Version) int {
	l := len(v1)
	if len(v2) < l {
		l = len(v2)
	}
	for i := 0; i < l; i++ {
		d := v1[i] - v2[i]
		if d != 0 {
			return d
		}
	}
	return len(v1) - len(v2)
}

func (v Version) String() string {
	s := ""
	for _, i := range v {
		if s != "" {
			s += "."
		}
		s += strconv.Itoa(i)
	}
	return s
}

// List of versions for sorting
type Versions []Version

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].CompareTo(v[j]) < 0
}
