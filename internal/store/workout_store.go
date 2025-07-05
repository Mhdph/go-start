package store

import (
	"database/sql"
)

type Workout struct {
	ID             int            `json:"id"`
	Title          string         `joson:"title"`
	Description    string         `joson:"description"`
	Duration       int            `joson:"duration"`
	CaloriesBurned int            `joson:"calories_burned"`
	Entries        []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID           int    `json:"id"`
	ExerciesName string `json:"exercise_name"`
	Sets         int    `json:"sets"`
	Reps         *int   `json:"reps"`
	Duration     *int   `json:"duration"`
	Weight       *int   `json:"weight"`
	Notes        string `json:"notes"`
	OrderIndex   int    `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkOut(*Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkOut(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := ` 
	INSERT INTO workouts (title, description, duration, calories_burned) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
	`

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.Duration, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query = `INSERT INTO workout_entries (workout_id,exercise_name,sets,reps,duration,weight,notes,order_index)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciesName, entry.Sets, entry.Reps, entry.Duration, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `SELECT id,title,description,duration,calories_burned FROM workouts WHERE id = $1`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.Duration, &workout.CaloriesBurned)
	if err != nil {
		return nil, err
	}

	entryQuery := `SELECT id,exercise_name, sets, reps, duration, wiegh, note, order_index
	FROM workout_entries 
	WHERE workout_id = $1 
	ORDER BY order_inedx
	`
	rows, err := pg.db.Query(entryQuery, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(
			&entry.ID,
			&entry.ExerciesName,
			&entry.Sets,
			&entry.Reps,
			&entry.Duration,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := ` 
	UPDATE workouts 
	SET title = $1, description = $2, duraion = $3, calories_burned = $4
	WHERE id = $5
	`

	result, err := tx.Exec(query, workout.Title, workout.Description, workout.Duration, workout.CaloriesBurned, workout.ID)

	if err != nil {
		return err
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}
	if RowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM workout_entrues WHERE workout_Id = $1`, workout.ID)

	if err != nil {
		return err
	}

	for _, entry := range workout.Entries {
		query := `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = tx.Exec(query, workout.ID, entry.ExerciesName, entry.Sets, entry.Reps, entry.Duration, entry.Weight, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts WHERE id = $1`
	_, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
