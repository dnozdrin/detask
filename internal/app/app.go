package app

import (
	"database/sql"
	"github.com/dnozdrin/detask/internal/app/log"
	"github.com/dnozdrin/detask/internal/delivery/http"
	"github.com/dnozdrin/detask/internal/delivery/http/rest"
	sv "github.com/dnozdrin/detask/internal/domain/services"
	pg "github.com/dnozdrin/detask/internal/infrastructure/storage/postgres"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	mg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // migrations from files
	_ "github.com/joho/godotenv/autoload"                // automatic env vars from .env files
	_ "github.com/lib/pq"                                // postgres driver
	"github.com/rs/cors"
	"go.uber.org/zap"
	stdhttp "net/http"
)

// App represents the main application handler
type App struct {
	config Config
	dbConf DbConfig

	DB     *sql.DB
	router *http.Router

	log *zap.SugaredLogger

	boardService   rest.BoardService
	columnService  rest.ColumnService
	taskService    rest.TaskService
	commentService rest.CommentService
}

// Initialize loads all required for application run dependencies
func (a *App) Initialize(dbConf DbConfig, config Config) {
	a.config = config
	a.dbConf = dbConf

	a.loadLogger()
	a.connectDB()
	a.migrateDb()
	a.loadServices()
	a.setupDelivery()
}

func (a *App) connectDB() {
	var err error

	a.DB, err = sql.Open(a.dbConf.driver, a.dbConf.toConnString())
	if err != nil {
		a.log.Fatalf("DB connection error: %v", err)
	}

	if err = a.DB.Ping(); err != nil {
		a.log.Fatalf("DB connection verification error: %v", err)
	}
}

func (a *App) migrateDb() {
	driver, err := mg.WithInstance(a.DB, &mg.Config{})
	if err != nil {
		a.log.Fatalf("DB migration: failed: %v", err)
	}
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

func (a *App) loadLogger() {
	defer func() {
		if r := recover(); r != nil {
			a.log.Fatalf("logger initialization failed: %v", r)
		}
	}()
	var cfg zap.Config
	switch a.config.context {
	case Prod:
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{a.config.logPath}
	case Test:
		cfg = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.FatalLevel),
			Development:       true,
			DisableStacktrace: true,
			Encoding:          "console",
			EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
			OutputPaths:       []string{a.config.logPath},
			ErrorOutputPaths:  []string{"stderr"},
		}
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		a.log.Fatalf("logger initialization failed: %v", err)
	}

	a.log = zapLogger.Sugar()
}

func (a *App) loadServices() {
	validatorImpl := NewValidator(validator.New(), a.log)

	var (
		boardStorage   sv.BoardStorage
		columnStorage  sv.ColumnStorage
		taskStorage    sv.TaskStorage
		commentStorage sv.CommentStorage
	)

	switch a.dbConf.driver {
	case "postgres":
		boardStorage = pg.NewBoardDAO(a.DB, a.log)
		columnStorage = pg.NewColumnDAO(a.DB, a.log)
		taskStorage = pg.NewTaskDAO(a.DB, a.log)
		commentStorage = pg.NewCommentsDAO(a.DB, a.log)
	default:
		a.log.Fatalf("%s driver support is not implemented", a.dbConf.driver)
	}

	a.boardService = sv.NewBoardService(validatorImpl, boardStorage)
	a.columnService = sv.NewColumnService(validatorImpl, columnStorage)
	a.taskService = sv.NewTaskService(validatorImpl, taskStorage)
	a.commentService = sv.NewCommentService(validatorImpl, commentStorage)
}

func (a *App) setupDelivery() {
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

func (a *App) addCORSMiddleware(handler stdhttp.Handler) stdhttp.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: a.config.allowedOrigins,
		AllowedMethods: []string{"HEAD", "GET", "POST", "DELETE", "PUT"},
		Debug:          a.config.context == Dev,
	})

	c.Log = log.NewCORSLogger(a.log)
	return c.Handler(handler)
}

// Run will start the web server on the given address
func (a *App) Run(addr string) {
	if err := http.NewServer(a.addCORSMiddleware(a.router), a.log).Start(addr); err != nil {
		a.log.Fatalf("http: server: listen and server: %v", err)
	}

	a.syncLogger()
	a.closeDB()
}

// ServeHTTPInternal is used for end to end tests
func (a *App) ServeHTTPInternal(w stdhttp.ResponseWriter, req *stdhttp.Request) {
	a.router.ServeHTTP(w, req)
}

// syncLogger flushes any buffered log entries. Applications should take care
// to call Sync before exiting. Check for "sync /dev/stderr: invalid argument"
// error is added for development log preset and should be removed as soon as
// this is issue will be fixed in uber-go/zap
func (a *App) syncLogger() {
	if err := a.log.Sync(); err != nil {
		if err.Error() == "sync /dev/stderr: invalid argument" {
			a.log.Debug(err)
		} else {
			a.log.Errorf("logger sync error: %v", err)
		}
	}
}

// closeDB closes the database and prevents new queries from starting.
// Applications should take care to call CloseDB before exiting.
func (a *App) closeDB() {
	if err := a.DB.Close(); err != nil {
		a.log.Errorf("DB close error: %v", err)
	}
}
