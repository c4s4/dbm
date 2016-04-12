package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

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
