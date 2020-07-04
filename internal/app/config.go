package app

import "fmt"

const (
	Prod = "production"
	Test = "testing"
	Dev  = "development"
)

// Config represents the application configuration
type Config struct {
	context string
	logPath string
}

func NewConfig(context, logPath string) Config {
	if context != Prod && context != Test {
		context = Dev
	}

	return Config{
		context: context,
		logPath: logPath,
	}
}

type DbConfig struct {
	driver, host, name, user, password, port, mgPath string
}

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
