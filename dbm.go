package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// print message if error and exit
func printAndExitIfError(message string, err error) {
	if err != nil {
		fmt.Println(message+":", err.Error())
		os.Exit(1)
	}
}

// print message and exit
func printAndExit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

// list script directories to run
func ListDirectories(c *Configuration, v string) []string {
	if _, err := os.Stat(c.SqlDir); os.IsNotExist(err) {
		printAndExitIfError("SQL directory not found", err)
	}
	if fi, _ := os.Stat(c.SqlDir); !fi.IsDir() {
		printAndExit("SQL directory is not a directory")
	}
	dirs, err := ioutil.ReadDir(c.SqlDir)
	printAndExitIfError("Error listing SQL directory", err)
	files := Versions{}
	for _, d := range dirs {
		f, err := NewVersion(d.Name())
		printAndExitIfError(fmt.Sprintf("Bad version directory '%s'", d.Name()), err)
		files = append(files, f)
	}
	sort.Sort(files)
	version, err := NewVersion(v)
	printAndExitIfError("Error parsing version", err)
	selected := []string{}
	for _, v := range files {
		if v.CompareTo(version) <= 0 {
			selected = append(selected, v.Name)
		}
	}
	return selected
}

// select files to run in a given version directory
func ListScriptsForVersion(dir, version, platform string) []string {
	selected := []string{}
	directory := path.Join(dir, version)
	scripts, err := ioutil.ReadDir(directory)
	printAndExitIfError("Error listing script for version", err)
	for _, script := range scripts {
		name := script.Name()
		name = strings.TrimSuffix(name, filepath.Ext(name))
		if name == "all" || name == platform {
			selected = append(selected, script.Name())
		}
	}
	return selected
}

// List migration scripts to run
func ListScripts(configuration *Configuration, platform, version string) []string {
	dirs := ListDirectories(configuration, version)
	files := []string{}
	for _, dir := range dirs {
		scripts := ListScriptsForVersion(configuration.SqlDir, dir, platform)
		for _, script := range scripts {
			file := path.Join(configuration.SqlDir, dir, script)
			files = append(files, file)
		}
	}
	return files
}

// check command line parameters
func checkParameters(configuration Configuration, platform, version string) {
	found := false
	for _, p := range configuration.Platforms {
		if p == platform {
			found = true
			break
		}
	}
	if !found {
		message := fmt.Sprintf("Platform '%s' is unknown, must be one of %s",
			platform, strings.Join(configuration.Platforms, ", "))
		printAndExit(message)
	}
}

// run database migration
func run(platform, version, config string, dryRun bool) {
	configuration, err := LoadConfiguration(config)
	checkParameters(*configuration, platform, version)
	printAndExitIfError("Error loading configuration file", err)
	scripts := ListScripts(configuration, platform, version)
	printAndExitIfError("Error listing migration scripts", err)
	if dryRun {
		fmt.Printf("Script to migrate platform '%s' to version '%s':\n",
			platform, version)
		for _, script := range scripts {
			fmt.Println("- " + script)
		}
		os.Exit(0)
	}
}

// main: parse command line
func main() {
	config := flag.String("config", ".dbm.yml", "DBmigration onfiguration file")
	dryRun := flag.Bool("dry", false, "Dry run (print scripts for migration)")
	flag.Parse()
	if len(flag.Args()) < 2 {
		printAndExit("You must pass platform and version on command line")
	}
	platform := flag.Args()[0]
	version := flag.Args()[1]
	run(platform, version, *config, *dryRun)
}
