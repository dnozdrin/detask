package app

import "fmt"

const (
	// Prod is a context for production usage
	Prod = "production"
	// Test is a context for running automatic tests
	Test = "testing"
	// Dev is the default context
	Dev = "development"
)

// Config represents the application configuration
type Config struct {
	context string
	logPath string
}

// NewConfig is a Config constructor
func NewConfig(context, logPath string) Config {
	if context != Prod && context != Test {
		context = Dev
	}

	return Config{
		context: context,
		logPath: logPath,
	}
}

// DbConfig represents configuration required for DB connection
type DbConfig struct {
	driver, host, name, user, password, port, mgPath string
}

// NewDBConfig is a DbConfig constructor
func NewDBConfig(driver, host, name, user, password, port, mgPath string) DbConfig {
	return DbConfig{
		driver:   driver,
		host:     host,
		name:     name,
		user:     user,
		password: password,
		port:     port,
		mgPath:   mgPath,
	}
}

func (c DbConfig) toConnString() string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s port=%s sslmode=disable",
		c.host,
		c.name,
		c.user,
		c.password,
		c.port,
	)
}
