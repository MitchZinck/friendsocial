package manual_activities

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ManualActivity struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	ActivityID    *int      `json:"activity_id,omitempty"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	EstimatedTime *string   `json:"estimated_time,omitempty"` // Interval type stored as string for simplicity
	LocationID    *int      `json:"location_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	IsActive      bool      `json:"is_active"`
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

func (s *Service) Create(manualActivity ManualActivity) (ManualActivity, error) {
	s.Lock()
	defer s.Unlock()

	err := s.db.QueryRow(
		context.Background(),
		`INSERT INTO manual_activities 
		(user_id, activity_id, name, description, estimated_time, location_id, is_active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, created_at`,
		manualActivity.UserID, manualActivity.ActivityID, manualActivity.Name,
		manualActivity.Description, manualActivity.EstimatedTime, manualActivity.LocationID, manualActivity.IsActive,
	).Scan(&manualActivity.ID, &manualActivity.CreatedAt)

	if err != nil {
		return ManualActivity{}, err
	}

	return manualActivity, nil
}

func (s *Service) ReadAll() ([]ManualActivity, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(),
		"SELECT id, user_id, activity_id, name, description, estimated_time, location_id, created_at, is_active FROM manual_activities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manualActivities []ManualActivity
	for rows.Next() {
		var manualActivity ManualActivity
		err := rows.Scan(
			&manualActivity.ID, &manualActivity.UserID, &manualActivity.ActivityID, &manualActivity.Name,
			&manualActivity.Description, &manualActivity.EstimatedTime, &manualActivity.LocationID,
			&manualActivity.CreatedAt, &manualActivity.IsActive)
		if err != nil {
			return nil, err
		}
		manualActivities = append(manualActivities, manualActivity)
	}

	return manualActivities, nil
}

func (s *Service) Read(id string) (ManualActivity, bool, error) {
	s.Lock()
	defer s.Unlock()

	var manualActivity ManualActivity
	err := s.db.QueryRow(
		context.Background(),
		`SELECT id, user_id, activity_id, name, description, estimated_time, location_id, created_at, is_active 
		FROM manual_activities WHERE id = $1`, id).Scan(
		&manualActivity.ID, &manualActivity.UserID, &manualActivity.ActivityID, &manualActivity.Name,
		&manualActivity.Description, &manualActivity.EstimatedTime, &manualActivity.LocationID,
		&manualActivity.CreatedAt, &manualActivity.IsActive)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ManualActivity{}, false, nil
		}
		return ManualActivity{}, false, err
	}

	return manualActivity, true, nil
}

func (s *Service) Update(id string, manualActivity ManualActivity) (ManualActivity, bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(
		context.Background(),
		`UPDATE manual_activities 
		SET user_id = $1, activity_id = $2, name = $3, description = $4, 
		estimated_time = $5, location_id = $6, is_active = $7 
		WHERE id = $8`,
		manualActivity.UserID, manualActivity.ActivityID, manualActivity.Name,
		manualActivity.Description, manualActivity.EstimatedTime, manualActivity.LocationID,
		manualActivity.IsActive, id)

	if err != nil {
		return ManualActivity{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return ManualActivity{}, false, nil
	}

	return manualActivity, true, nil
}

func (s *Service) Delete(id string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(context.Background(), "DELETE FROM manual_activities WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
