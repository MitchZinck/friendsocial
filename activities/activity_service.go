package activities

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Activity struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Emoji         string `json:"emoji"` // Add this field
	Description   string `json:"description"`
	EstimatedTime string `json:"estimated_time"` // Interval type stored as string for simplicity
	LocationID    int    `json:"location_id"`
	UserCreated   bool   `json:"user_created"` // Add this field
}

type Service struct {
	sync.Mutex
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		db: db,
	}
}

func (activityService *Service) Create(activity Activity) (Activity, error) {
	activityService.Lock()
	defer activityService.Unlock()

	err := activityService.db.QueryRow(
		context.Background(),
		"INSERT INTO activities (name, emoji, description, estimated_time, location_id, user_created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		activity.Name, activity.Emoji, activity.Description, activity.EstimatedTime, activity.LocationID, activity.UserCreated,
	).Scan(&activity.ID)

	if err != nil {
		return Activity{}, err
	}

	return activity, nil
}

func (activityService *Service) ReadAll() ([]Activity, error) {
	activityService.Lock()
	defer activityService.Unlock()

	rows, err := activityService.db.Query(context.Background(), "SELECT id, name, emoji, description, estimated_time::text, location_id, user_created FROM activities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.ID, &activity.Name, &activity.Emoji, &activity.Description, &activity.EstimatedTime, &activity.LocationID, &activity.UserCreated); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (activityService *Service) Read(id string) (Activity, bool, error) {
	activityService.Lock()
	defer activityService.Unlock()

	var activity Activity
	err := activityService.db.QueryRow(context.Background(),
		`SELECT id, name, emoji, description, estimated_time::text, location_id, user_created 
		FROM activities 
		WHERE id = $1`, id).Scan(&activity.ID, &activity.Name, &activity.Emoji, &activity.Description, &activity.EstimatedTime, &activity.LocationID, &activity.UserCreated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Activity{}, false, nil
		}
		return Activity{}, false, err
	}

	return activity, true, nil
}

func (activityService *Service) Update(id string, activity Activity) (Activity, bool, error) {
	activityService.Lock()
	defer activityService.Unlock()

	cmdTag, err := activityService.db.Exec(context.Background(),
		"UPDATE activities SET name = $1, emoji = $2, description = $3, estimated_time = $4, location_id = $5, user_created = $6 WHERE id = $7",
		activity.Name, activity.Emoji, activity.Description, activity.EstimatedTime, activity.LocationID, activity.UserCreated, id)

	if err != nil {
		return Activity{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return Activity{}, false, nil
	}

	return activity, true, nil
}

func (activityService *Service) Delete(id string) (bool, error) {
	activityService.Lock()
	defer activityService.Unlock()

	cmdTag, err := activityService.db.Exec(context.Background(), "DELETE FROM activities WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
