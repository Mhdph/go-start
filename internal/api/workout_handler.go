package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mhdph/go-start/internal/middleware"
	"github.com/mhdph/go-start/internal/store"
	"github.com/mhdph/go-start/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}

}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDPram: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id"})
		return
	}
	workout, err := wh.workoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workout.UserID = user.ID

	createdWorkout, err := wh.workoutStore.CreateWorkOut(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}
func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")

	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutId)

	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	var updateWorkoutRequest struct {
		Title          *string              `joson:"title"`
		Description    *string              `joson:"description"`
		Duration       *int                 `joson:"duration"`
		CaloriesBurned *int                 `joson:"calories_burned"`
		Entries        []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.Duration != nil {
		existingWorkout.Duration = *updateWorkoutRequest.Duration
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	user := middleware.GetUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workOutOwner, err := wh.workoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		http.Error(w, "failed to fetch workout", http.StatusInternalServerError)
		return
	}

	if workOutOwner.UserID != user.ID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, "failed to update workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")

	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)

	if err != nil {
		http.Error(w, "failed to delete workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
