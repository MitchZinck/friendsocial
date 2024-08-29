package user_activities

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserActivity struct {
	ID         int  `json:"id"`
	UserID     int  `json:"user_id"`
	ActivityID int  `json:"activity_id"`
	IsActive   bool `json:"is_active"`
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

// Create a new user activity
func (service *Service) Create(userActivity UserActivity) (UserActivity, error) {
	service.Lock()
	defer service.Unlock()

	var id int
	err := service.db.QueryRow(
		context.Background(),
		"INSERT INTO user_activities (user_id, activity_id, is_active) VALUES ($1, $2, $3) RETURNING id",
		userActivity.UserID, userActivity.ActivityID, userActivity.IsActive,
	).Scan(&id)
	if err != nil {
		return UserActivity{}, err
	}

	userActivity.ID = id
	return userActivity, nil
}

// Read all user activities for a specific user
func (service *Service) ReadAll(userID string) ([]UserActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, user_id, activity_id, is_active FROM user_activities WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userActivities []UserActivity
	for rows.Next() {
		var userActivity UserActivity
		if err := rows.Scan(&userActivity.ID, &userActivity.UserID, &userActivity.ActivityID, &userActivity.IsActive); err != nil {
			return nil, err
		}
		userActivities = append(userActivities, userActivity)
	}

	return userActivities, nil
}

// Read a specific user activity by ID
func (service *Service) Read(id string) (UserActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	var userActivity UserActivity
	err := service.db.QueryRow(context.Background(), "SELECT id, user_id, activity_id, is_active FROM user_activities WHERE id = $1", id).Scan(&userActivity.ID, &userActivity.UserID, &userActivity.ActivityID, &userActivity.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return UserActivity{}, false, nil
		}
		return UserActivity{}, false, err
	}

	return userActivity, true, nil
}

// Update an existing user activity
func (service *Service) Update(id string, userActivity UserActivity) (UserActivity, bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "UPDATE user_activities SET user_id = $1, activity_id = $2, is_active = $3 WHERE id = $4",
		userActivity.UserID, userActivity.ActivityID, userActivity.IsActive, id)
	if err != nil {
		return UserActivity{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return UserActivity{}, false, nil
	}

	return userActivity, true, nil
}

// Delete a user activity by ID
func (service *Service) Delete(id string) (bool, error) {
	service.Lock()
	defer service.Unlock()

	cmdTag, err := service.db.Exec(context.Background(), "DELETE FROM user_activities WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

// Get all active user activities for a specific user
func (service *Service) GetActiveUserActivities(userID string) ([]UserActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, user_id, activity_id, is_active FROM user_activities WHERE user_id = $1 AND is_active = TRUE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activeUserActivities []UserActivity
	for rows.Next() {
		var userActivity UserActivity
		if err := rows.Scan(&userActivity.ID, &userActivity.UserID, &userActivity.ActivityID, &userActivity.IsActive); err != nil {
			return nil, err
		}
		activeUserActivities = append(activeUserActivities, userActivity)
	}

	return activeUserActivities, nil
}

// Get all inactive user activities for a specific user
func (service *Service) GetInactiveUserActivities(userID string) ([]UserActivity, error) {
	service.Lock()
	defer service.Unlock()

	rows, err := service.db.Query(context.Background(), "SELECT id, user_id, activity_id, is_active FROM user_activities WHERE user_id = $1 AND is_active = FALSE", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inactiveUserActivities []UserActivity
	for rows.Next() {
		var userActivity UserActivity
		if err := rows.Scan(&userActivity.ID, &userActivity.UserID, &userActivity.ActivityID, &userActivity.IsActive); err != nil {
			return nil, err
		}
		inactiveUserActivities = append(inactiveUserActivities, userActivity)
	}

	return inactiveUserActivities, nil
}
