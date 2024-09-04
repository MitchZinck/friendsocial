package scheduled_activities

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ScheduledActivity struct {
	ID          int       `json:"id"`
	ActivityID  int       `json:"activity_id"`
	IsActive    bool      `json:"is_active"`
	ScheduledAt time.Time `json:"scheduled_at"` // New field for scheduled_at
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

// Create a new scheduled activity
func (service *Service) Create(scheduledActivity ScheduledActivity) (ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	var id int
	err := service.db.QueryRow(
		context.Background(),
		"INSERT INTO scheduled_activities (activity_id, is_active, scheduled_at) VALUES ($1, $2, $3) RETURNING id",
		scheduledActivity.ActivityID, scheduledActivity.IsActive, scheduledActivity.ScheduledAt,
	).Scan(&id)
	if err != nil {
		return ScheduledActivity{}, err
	}

	scheduledActivity.ID = id
	return scheduledActivity, nil
}

// Read all user activities for a specific user
func (service *Service) ReadAll() ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at FROM scheduled_activities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt); err != nil {
			return nil, err
		}
		scheduledActivities = append(scheduledActivities, scheduledActivity)
	}

	return scheduledActivities, nil
}

// Read a specific user activity by ID
func (service *Service) Read(id string) (ScheduledActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	var scheduledActivity ScheduledActivity
	err := service.db.QueryRow(context.Background(), "SELECT id, activity_id, is_active, scheduled_at FROM scheduled_activities WHERE id = $1", id).Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ScheduledActivity{}, false, nil
		}
		return ScheduledActivity{}, false, err
	}

	return scheduledActivity, true, nil
}

// Update an existing user activity
func (service *Service) Update(id string, scheduledActivity ScheduledActivity) (ScheduledActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "UPDATE scheduled_activities SET activity_id = $1, is_active = $2, scheduled_at = $3 WHERE id = $4",
		scheduledActivity.ActivityID, scheduledActivity.IsActive, scheduledActivity.ScheduledAt, id)
	if err != nil {
		return ScheduledActivity{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return ScheduledActivity{}, false, nil
	}

	return scheduledActivity, true, nil
}

// Delete a user activity by ID
func (service *Service) Delete(id string) (bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "DELETE FROM scheduled_activities WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

// Get all active user activities for a specific user
func (service *Service) GetActiveScheduledActivities(userID string) ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at FROM scheduled_activities WHERE is_active = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeScheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt); err != nil {
			return nil, err
		}
		activeScheduledActivities = append(activeScheduledActivities, scheduledActivity)
	}

	return activeScheduledActivities, nil
}

// Get all inactive user activities for a specific user
func (service *Service) GetInactiveScheduledActivities(userID string) ([]ScheduledActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, activity_id, is_active, scheduled_at FROM scheduled_activities WHERE is_active = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inactiveScheduledActivities []ScheduledActivity
	for rows.Next() {
		var scheduledActivity ScheduledActivity
		if err := rows.Scan(&scheduledActivity.ID, &scheduledActivity.ActivityID, &scheduledActivity.IsActive, &scheduledActivity.ScheduledAt); err != nil {
			return nil, err
		}
		inactiveScheduledActivities = append(inactiveScheduledActivities, scheduledActivity)
	}

	return inactiveScheduledActivities, nil
}
