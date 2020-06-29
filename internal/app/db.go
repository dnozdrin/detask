package app

import "fmt"

type DBConf struct {
	driver, host, name, user, password, port string
}

func NewDBconf(driver, host, name, user, password, port string) DBConf {
	return DBConf{driver, host, name, user, password, port}
}

func (c DBConf) toConnString() string {
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s port=%s sslmode=disable",
		c.host,
		c.name,
		c.user,
		c.password,
		c.port,
	)
}
