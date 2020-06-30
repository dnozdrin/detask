// todo: review logging format
// todo: review app-wide logging policy
// todo: review errors handling, consider wrapping and unwrapping
// todo: consider using values instead of pointers
package main

import "os"

func main() {
	a := app{}
	a.initialize(
		newDBConfig(
			"postgres",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_MIGRATION_PATH"),
		),
		newAppConfig(
			os.Getenv("APP_CONTEXT"),
			os.Getenv("APP_LOG_PATH"),
		),
	)

	defer a.syncLogger()
	defer a.closeDB()

	a.run(":" + os.Getenv("APP_PORT"))
}
