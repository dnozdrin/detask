package main

import "fmt"

type appConfig struct {
	context string
	logPath string
}

func newAppConfig(context, logPath string) appConfig {
	return appConfig{
		context: context,
		logPath: logPath,
	}
}

type dbConfig struct {
	driver, host, name, user, password, port, mgPath string
}

func newDBConfig(driver, host, name, user, password, port, mgPath string) dbConfig {
	return dbConfig{
		driver:   driver,
		host:     host,
		name:     name,
		user:     user,
		password: password,
		port:     port,
		mgPath:   mgPath,
	}
}

func (c dbConfig) toConnString() string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s port=%s sslmode=disable",
		c.host,
		c.name,
		c.user,
		c.password,
		c.port,
	)
}
