// todo: review logging format
// todo: review app-wide logging policy
// todo: consider adding development option for zap logger
// todo: consider refactoring to make slim main()
// todo: consider moving db connection data to app config
// todo: review errors handling, consider wrapping and unwrapping
// todo: consider using values instead of pointers
package main

import (
	"github.com/dnozdrin/detask/internal/app"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {
	a := app.App{}
	a.Initialize(app.NewDBconf(
		"postgres",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_PORT"),
	), os.Getenv("APP_CONTEXT"))

	defer a.SyncLogger()
	defer a.CloseDB()

	a.Run(":" + os.Getenv("APP_PORT"))
}
