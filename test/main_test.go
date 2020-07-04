// +build integrational

// todo: consider replacing maps with concrete models
// todo: consider optimizing seeding
// todo: cover with test boundary cases: validation, etc
package test

import (
	"github.com/dnozdrin/detask/internal/app"
	"os"
	"testing"
)

var a app.App

func TestMain(m *testing.M) {
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
			app.Test,
			"stderr",
		),
	)
	code := m.Run()
	os.Exit(code)
}
