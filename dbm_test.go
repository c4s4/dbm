package main

import (
	"sort"
	"testing"
)

func TestParseConfiguration(t *testing.T) {
	source := `# platform list
platforms:
- itg
- prp
- prod
# default platform
default-platform: itg
# platform where init is forbidden
critical-platforms: [prod]
# directory for SQL scripts (relative to this file)
sql-dir: sql
# charset of the database
database-charset: utf8
# Database configuration for environments
database:
  itg:
    hostname: localhost
    database: test
    password: test
    username: test
  prp:
    hostname: localhost
    database: test
    password: test
    username: test
  prod:
    hostname: localhost
    database: test
    password: test
    username: test`
	configuration, err := ParseConfiguration([]byte(source))
	if err != nil {
		t.Fatalf("Error parsing configuration: %s", err)
	}
	if configuration.Platforms[0] != "itg" ||
		configuration.Platforms[1] != "prp" ||
		configuration.Platforms[2] != "prod" ||
		configuration.DefaultPlatform != "itg" ||
		len(configuration.CriticalPlatforms) != 1 ||
		configuration.CriticalPlatforms[0] != "prod" ||
		configuration.SqlDir != "sql" ||
		configuration.DatabaseCharset != "utf8" ||
		len(configuration.Database) != 3 ||
		configuration.Database["itg"].Hostname != "localhost" ||
		configuration.Database["itg"].Database != "test" ||
		configuration.Database["itg"].Username != "test" ||
		configuration.Database["itg"].Password != "test" {
		t.Fatal("Loaded bad configuration")
	}
}

func TestNewVersion(t *testing.T) {
	v, err := NewVersion("1.2.3")
	if err != nil {
		t.Errorf("Got error parsing version '1.2.3': %s", err)
	}
	if len(v) != 3 {
		t.Errorf("Bad version length: got %v instead of 3", len(v))
	}
	if v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Error("Bad error parsing")
	}
	v, err = NewVersion("1")
	if err != nil {
		t.Errorf("Got error parsing version '1': %s", err)
	}
	if len(v) != 1 {
		t.Errorf("Bad version length: got %v instead of 1", len(v))
	}
	v, err = NewVersion("init")
	if err != nil {
		t.Error("Could not parse 'init' version")
	}
	if len(v) != 1 || v[0] != -1 {
		t.Error("Bad 'init' version parsing")
	}
}

func TestBadVersion(t *testing.T) {
	_, err := NewVersion("foo")
	if err == nil {
		t.Error("Got no error parsing version 'foo'")
	}
	if err.Error() != "Error parsing version 'foo'" {
		t.Error("Got bad error parsing version 'foo'")
	}
	_, err = NewVersion("")
	if err == nil {
		t.Error("Got no error parsing void version")
	}
	_, err = NewVersion("-1")
	if err == nil {
		t.Error("Got no error parsing negative version")
	}
}

func TestVersionCompareTo(t *testing.T) {
	v1, _ := NewVersion("0")
	v2, _ := NewVersion("0")
	if v1.CompareTo(v2) != 0 {
		t.Fail()
	}
	v2, _ = NewVersion("0.0")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("0.1")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v1, _ = NewVersion("1.2.3")
	v2, _ = NewVersion("2.3.4")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.4")
	if v1.CompareTo(v2) >= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.2")
	if v1.CompareTo(v2) <= 0 {
		t.Fail()
	}
	v2, _ = NewVersion("1.2.3")
	if v1.CompareTo(v2) != 0 {
		t.Fail()
	}
}

func TestVersionString(t *testing.T) {
	s := "1.2.3"
	v, _ := NewVersion(s)
	if v.String() != s {
		t.Fail()
	}
}

func TestSortVersions(t *testing.T) {
	v1, _ := NewVersion("1.2.3")
	v2, _ := NewVersion("1.2")
	v3, _ := NewVersion("1")
	v4, _ := NewVersion("1.2.4")
	v5, _ := NewVersion("2.2")
	v6, _ := NewVersion("0.1")
	v7, _ := NewVersion("init")
	v := []Version{v1, v2, v3, v4, v5, v6, v7}
	sort.Sort(Versions(v))
	if v[0].CompareTo(v7) != 0 ||
		v[1].CompareTo(v6) != 0 ||
		v[2].CompareTo(v3) != 0 ||
		v[3].CompareTo(v2) != 0 ||
		v[4].CompareTo(v1) != 0 ||
		v[5].CompareTo(v4) != 0 ||
		v[6].CompareTo(v5) != 0 {
		t.Fail()
	}
}
