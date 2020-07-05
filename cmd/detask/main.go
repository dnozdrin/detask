//+build !test

package main

import (
	"github.com/dnozdrin/detask/internal/app"
	"os"
)

func main() {
	a := app.App{}
	a.Initialize(
		app.NewDBConfig(
			"postgres",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_MIGRATION_PATH"),
		),
		app.NewConfig(
			os.Getenv("APP_CONTEXT"),
			os.Getenv("APP_LOG_PATH"),
		),
	)

	a.Run(":" + os.Getenv("PORT"))
}
