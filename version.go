package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
