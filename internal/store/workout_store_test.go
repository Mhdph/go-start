package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		t.Fatalf("openin test db %v", err)
	}

	err = Migrate(db, "../../migraions/")

	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec("TRUNCATE workouts,workout_entries CASCADE")

	if err != nil {
		t.Fatalf("truncate test db error: %v", err)

	}

	return db

}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)

	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:          "Push day",
				Description:    "push for day 1",
				Duration:       60,
				CaloriesBurned: 200,
				Entries: []WorkoutEntry{
					{
						ExerciesName: "Bench press",
						Sets:         3,
						Reps:         IntPtr(10),
					},
				},
			},
			wantErr: false,
		},

		{
			name: "invalid workout",
			workout: &Workout{
				Title: "Push day",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkOut(test.workout)
			if (err != nil) != test.wantErr {
				t.Errorf("CreateWorkout() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				assert.Error(t, err)

			}
			require.NoError(t, err)

			assert.Equal(t, test.workout.Title, createdWorkout.Title)

		})
	}
}

func IntPtr(i int) *int {
	return &i
}
