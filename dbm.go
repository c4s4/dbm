package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
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
	if v == "init" {
		version = append(version, -1)
		return version, nil
	}
	for _, s := range strings.Split(v, ".") {
		i, err := strconv.Atoi(s)
		if err != nil || i < 0 {
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
	if len(v) == 1 && v[0] == -1 {
		return "init"
	}
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

func NewVersions(l []string) (Versions, error) {
	versions := []Version{}
	for _, s := range l {
		v, err := NewVersion(s)
		if err != nil {
			return []Version{}, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].CompareTo(v[j]) < 0
}

// print if error and exit
func printError(message string, err error) {
	if err != nil {
		fmt.Println(message+":", err.Error())
		os.Exit(1)
	}
}

func printAndExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

// list script directories to run
func ListDirectories(c *Configuration, v string) ([]string, error) {
	if _, err := os.Stat(c.SqlDir); os.IsNotExist(err) {
		printError("SQL directory not found", err)
	}
	if fi, _ := os.Stat(c.SqlDir); !fi.IsDir() {
		printAndExit("SQL directory is not a directory")
	}
	dirs, err := ioutil.ReadDir(c.SqlDir)
	printError("Error listing SQL directory", err)
	files := Versions{}
	for _, d := range dirs {
		f, err := NewVersion(d.Name())
		printError(fmt.Sprintf("Bad version directory '%s'", d.Name()), err)
		files = append(files, f)
	}
	sort.Sort(files)
	version, err := NewVersion(v)
	printError("Error parsing version", err)
	selected := []string{}
	for _, v := range files {
		if v.CompareTo(version) <= 0 {
			selected = append(selected, v.String())
		}
	}
	return selected, nil
}

// select files to run in a given version directory
func ListScriptsForVersion(dir, version, platform string) ([]string, error) {
	selected := []string{}
	directory := path.Join(dir, version)
	scripts, err := ioutil.ReadDir(directory)
	printError("Error listing script for version", err)
	for _, script := range scripts {
		name := script.Name()
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if name == "all" || name == platform {
			selected = append(selected, script.Name())
		}
	}
	return selected, nil
}

// List migration scripts to run
func ListScripts(configuration *Configuration, platform, version string) ([]string, error) {
	dirs, err := ListDirectories(configuration, version)
	printError("Error listing SQL directories", err)
	files := []string{}
	for _, dir := range dirs {
		fs, err := ListScriptsForVersion(configuration.SqlDir, dir, platform)
		printError("Error listing scripts for version", err)
		for _, f := range fs {
			files = append(files, path.Join(configuration.SqlDir, dir, f))
		}
	}
	return files, nil
}

// run database migration
func run(config, platform, version string) {
	configuration, err := LoadConfiguration(config)
	printError("Error loading configuration file", err)
	scripts, err := ListScripts(configuration, platform, version)
	printError("Error listing migration scripts", err)
	fmt.Printf("scripts: %#v\n", scripts)
}

// main: parse command line
func main() {
	config := ".dbm.yml"
	platform := os.Args[1]
	version := os.Args[2]
	run(config, platform, version)
}
