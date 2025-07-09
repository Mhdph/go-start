package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mhdph/go-start/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Autheniticate)
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
		r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutById)
		r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutById)
	})

	r.Post("/users", app.UserHandler.HandleRegisterUser)

	r.Post("/tokens/authetication", app.TokenHandler.HandleCreateToken)

	return r
}
