package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	INIT_NAME   = "init"
	INIT_NUMBER = -1
	SEPARATOR   = "."
)

// version as a name and an array of integers
type Version struct {
	Name    string
	Numbers []int
}

func NewVersion(name string) (Version, error) {
	if name == INIT_NAME {
		return Version{
			Name:    INIT_NAME,
			Numbers: []int{INIT_NUMBER},
		}, nil
	}
	numbers := []int{}
	for _, number := range strings.Split(name, SEPARATOR) {
		integer, err := strconv.Atoi(number)
		if err != nil || integer < 0 {
			return Version{}, fmt.Errorf("Error parsing version '%s'", number)
		}
		numbers = append(numbers, integer)
	}
	return Version{
		Numbers: numbers,
		Name:    name,
	}, nil
}

func (version1 Version) CompareTo(version2 Version) int {
	length := len(version1.Numbers)
	if len(version2.Numbers) < length {
		length = len(version2.Numbers)
	}
	for index := 0; index < length; index++ {
		diff := version1.Numbers[index] - version2.Numbers[index]
		if diff != 0 {
			return diff
		}
	}
	return len(version1.Numbers) - len(version2.Numbers)
}

// List of versions for sorting
type Versions []Version

func NewVersions(names []string) (Versions, error) {
	versions := []Version{}
	for _, name := range names {
		version, err := NewVersion(name)
		if err != nil {
			return []Version{}, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (versions Versions) Len() int {
	return len(versions)
}

func (versions Versions) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}

func (versions Versions) Less(i, j int) bool {
	return versions[i].CompareTo(versions[j]) < 0
}
