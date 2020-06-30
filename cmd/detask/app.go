package main

import (
	"database/sql"
	internal "github.com/dnozdrin/detask/internal/app"
	"github.com/dnozdrin/detask/internal/delivery/http"
	"github.com/dnozdrin/detask/internal/delivery/http/rest"
	sv "github.com/dnozdrin/detask/internal/domain/services"
	pg "github.com/dnozdrin/detask/internal/infrastructure/storage/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	mg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	prod = "production"
	test = "testing"
)

// App represents the main application handler
type app struct {
	config appConfig
	dbConf dbConfig

	db *sql.DB

	router *http.Router
	log    *zap.SugaredLogger

	boardService   rest.BoardService
	columnService  rest.ColumnService
	taskService    rest.TaskService
	commentService rest.CommentService
}

// Initialize loads all required for application run dependencies
func (a *app) initialize(dbConf dbConfig, config appConfig) {
	a.config = config
	a.dbConf = dbConf

	a.loadLogger()
	a.connectDB()
	a.migrateDb()
	a.loadServices()
	a.setupDelivery()
}

func (a *app) connectDB() {
	var err error

	a.db, err = sql.Open(a.dbConf.driver, a.dbConf.toConnString())
	if err != nil {
		a.log.Fatalf("DB connection error: %v", err)
	}

	if err = a.db.Ping(); err != nil {
		a.log.Fatalf("DB connection verification error: %v", err)
	}
}

func (a *app) migrateDb() {
	driver, err := mg.WithInstance(a.db, &mg.Config{})
	m, err := migrate.NewWithDatabaseInstance(a.dbConf.mgPath, a.dbConf.driver, driver)
	if err != nil {
		a.log.Fatalf("DB migration: failed: %v", err)
	}
	err = m.Up()
	switch err {
	case migrate.ErrNoChange:
		a.log.Info("DB migration: database schema is already up to date")
	case nil:
		a.log.Info("DB migration: changes applied")
	default:
		a.log.Fatalf("DB migration: failed: %v", err)
	}
}

func (a *app) loadLogger() {
	var logInitFunc func(options ...zap.Option) (*zap.Logger, error)
	switch a.config.context {
	case prod:
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{a.config.logPath}
		logInitFunc = cfg.Build
	default:
		logInitFunc = zap.NewDevelopment
	}

	zapLogger, err := logInitFunc()
	if err != nil {
		a.log.Fatalf("logger initialization error: %v", err)
	}

	a.log = zapLogger.Sugar()
}

func (a *app) loadServices() {
	validatorImpl := internal.NewValidator(validator.New(), a.log)

	var (
		boardStorage   sv.BoardStorage
		columnStorage  sv.ColumnStorage
		taskStorage    sv.TaskStorage
		commentStorage sv.CommentStorage
	)

	switch a.dbConf.driver {
	case "postgres":
		boardStorage = pg.NewBoardDAO(a.db, a.log)
		columnStorage = pg.NewColumnDAO(a.db, a.log)
		taskStorage = pg.NewTaskDAO(a.db, a.log)
		commentStorage = pg.NewCommentsDAO(a.db, a.log)
	default:
		a.log.Fatalf("%s driver support is not implemented", a.dbConf.driver)
	}

	a.boardService = sv.NewBoardService(validatorImpl, boardStorage)
	a.columnService = sv.NewColumnService(validatorImpl, columnStorage)
	a.taskService = sv.NewTaskService(validatorImpl, taskStorage)
	a.commentService = sv.NewCommentService(validatorImpl, commentStorage)
}

func (a *app) setupDelivery() {
	a.router = http.NewRouter()
	subRouter := a.router.GetSubRouter("/api/v1")

	healthCheckHandler := rest.NewHealthCheck(a.log)
	boardHandle := rest.NewBoardHandler(a.boardService, a.log, subRouter)
	columnHandler := rest.NewColumnHandler(a.columnService, a.log, subRouter)
	taskHandler := rest.NewTaskHandler(a.taskService, a.log, subRouter)
	commentHandler := rest.NewCommentHandler(a.commentService, a.log, subRouter)

	var routes = http.Routes{
		http.Route{Pattern: "/health", Method: "GET", Name: "health", HandlerFunc: healthCheckHandler.Status},

		http.Route{Pattern: "/board", Method: "POST", Name: "new_board", HandlerFunc: boardHandle.Create},
		http.Route{Pattern: "/boards", Method: "GET", Name: "get_boards", HandlerFunc: boardHandle.Get},
		http.Route{Pattern: "/boards/{id:[0-9]+}", Method: "GET", Name: "get_board", HandlerFunc: boardHandle.GetOneById},
		http.Route{Pattern: "/boards/{id:[0-9]+}", Method: "PUT", Name: "update_board", HandlerFunc: boardHandle.Update},
		http.Route{Pattern: "/boards/{id:[0-9]+}", Method: "DELETE", Name: "delete_board", HandlerFunc: boardHandle.Delete},

		http.Route{Pattern: "/column", Method: "POST", Name: "new_column", HandlerFunc: columnHandler.Create},
		http.Route{Pattern: "/columns", Method: "GET", Name: "get_columns", HandlerFunc: columnHandler.Get},
		http.Route{Pattern: "/columns/{id:[0-9]+}", Method: "GET", Name: "get_column", HandlerFunc: columnHandler.GetOneById},
		http.Route{Pattern: "/columns/{id:[0-9]+}", Method: "PUT", Name: "update_column", HandlerFunc: columnHandler.Update},
		http.Route{Pattern: "/columns/{id:[0-9]+}", Method: "DELETE", Name: "delete_column", HandlerFunc: columnHandler.Delete},

		http.Route{Pattern: "/task", Method: "POST", Name: "create_task", HandlerFunc: taskHandler.Create},
		http.Route{Pattern: "/tasks", Method: "GET", Name: "get_tasks", HandlerFunc: taskHandler.Get},
		http.Route{Pattern: "/tasks/{id:[0-9]+}", Method: "GET", Name: "get_task", HandlerFunc: taskHandler.GetOneById},
		http.Route{Pattern: "/tasks/{id:[0-9]+}", Method: "PUT", Name: "update_task", HandlerFunc: taskHandler.Update},
		http.Route{Pattern: "/tasks/{id:[0-9]+}", Method: "DELETE", Name: "delete_task", HandlerFunc: taskHandler.Delete},

		http.Route{Pattern: "/comment", Method: "POST", Name: "create_comment", HandlerFunc: commentHandler.Create},
		http.Route{Pattern: "/comments", Method: "GET", Name: "get_comments", HandlerFunc: commentHandler.Get},
		http.Route{Pattern: "/comments/{id:[0-9]+}", Method: "GET", Name: "get_comment", HandlerFunc: commentHandler.GetOneById},
		http.Route{Pattern: "/comments/{id:[0-9]+}", Method: "PUT", Name: "update_comment", HandlerFunc: commentHandler.Update},
		http.Route{Pattern: "/comments/{id:[0-9]+}", Method: "DELETE", Name: "delete_comment", HandlerFunc: commentHandler.Delete},
	}

	for _, route := range routes {
		subRouter.Register(route)
	}
}

// Run will start the web server on the given address
func (a *app) run(addr string) {
	http.NewServer(a.router, a.log).Start(addr)
}

// SyncLogger flushes any buffered log entries. Applications should take care
// to call Sync before exiting. Check for "sync /dev/stderr: invalid argument"
// error is added for development log preset and should be removed as soon as
// this is issue will be fixed in uber-go/zap
func (a *app) syncLogger() {
	if err := a.log.Sync(); err != nil {
		if err.Error() == "sync /dev/stderr: invalid argument" {
			a.log.Debug(err)
		} else {
			a.log.Errorf("logger sync error: %v", err)
		}
	}
}

// CloseDB closes the database and prevents new queries from starting.
// Applications should take care to call CloseDB before exiting.
func (a *app) closeDB() {
	if err := a.db.Close(); err != nil {
		a.log.Errorf("DB close error: %v", err)
	}
}
