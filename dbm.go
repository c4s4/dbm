package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
