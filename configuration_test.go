package main

import (
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
