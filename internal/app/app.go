package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mhdph/go-start/internal/api"
	"github.com/mhdph/go-start/internal/middleware"
	"github.com/mhdph/go-start/internal/store"
	"github.com/mhdph/go-start/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddlware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()

	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDb)
	userStore := store.NewPostgresUserStore(pgDb)
	tokenStore := store.NewPostgresTokenStore(pgDb)
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		DB:             pgDb,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available\n")
}
