package user_activity_preferences

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserActivityPreference struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	ActivityID      int    `json:"activity_id"`
	Frequency       int    `json:"frequency"`
	FrequencyPeriod string `json:"frequency_period"`
	DaysOfWeek      string `json:"days_of_week"`
}

type Service struct {
	sync.Mutex
	db       *pgxpool.Pool
	services *map[string]interface{}
}

func NewService(db *pgxpool.Pool, services *map[string]interface{}) *Service {
	return &Service{
		db:       db,
		services: services,
	}
}

func (s *Service) Create(preference UserActivityPreference) (UserActivityPreference, error) {
	s.Lock()
	defer s.Unlock()
	var id int
	err := s.db.QueryRow(
		context.Background(),
		`INSERT INTO user_activity_preferences (user_id, activity_id, frequency, frequency_period, days_of_week) 
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		preference.UserID, preference.ActivityID, preference.Frequency, preference.FrequencyPeriod, preference.DaysOfWeek,
	).Scan(&id)
	if err != nil {
		return UserActivityPreference{}, err
	}
	preference.ID = id

	return preference, nil
}

func (s *Service) ReadAll() ([]UserActivityPreference, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(), "SELECT id, user_id, activity_id, frequency, frequency_period, days_of_week FROM user_activity_preferences")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []UserActivityPreference
	for rows.Next() {
		var preference UserActivityPreference
		if err := rows.Scan(&preference.ID, &preference.UserID, &preference.ActivityID, &preference.Frequency, &preference.FrequencyPeriod, &preference.DaysOfWeek); err != nil {
			return nil, err
		}
		preferences = append(preferences, preference)
	}

	return preferences, nil
}

func (s *Service) Read(id string) (UserActivityPreference, bool, error) {
	s.Lock()
	defer s.Unlock()

	var preference UserActivityPreference
	err := s.db.QueryRow(context.Background(), "SELECT id, user_id, activity_id, frequency, frequency_period, days_of_week FROM user_activity_preferences WHERE id = $1", id).Scan(&preference.ID, &preference.UserID, &preference.ActivityID, &preference.Frequency, &preference.FrequencyPeriod, &preference.DaysOfWeek)
	if err != nil {
		if err == pgx.ErrNoRows {
			return UserActivityPreference{}, false, nil
		}
		return UserActivityPreference{}, false, err
	}

	return preference, true, nil
}

func (s *Service) Update(id string, preference UserActivityPreference) (UserActivityPreference, bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(context.Background(), "UPDATE user_activity_preferences SET user_id = $1, activity_id = $2, frequency = $3, frequency_period = $4 WHERE id = $5", preference.UserID, preference.ActivityID, preference.Frequency, preference.FrequencyPeriod, id)
	if err != nil {
		return UserActivityPreference{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return UserActivityPreference{}, false, nil
	}

	return preference, true, nil
}

func (s *Service) Delete(id string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(context.Background(), "DELETE FROM user_activity_preferences WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

// Add a new method to read preferences by user ID
func (s *Service) ReadByUserID(userID string) ([]UserActivityPreference, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(), "SELECT id, user_id, activity_id, frequency, frequency_period, days_of_week FROM user_activity_preferences WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var preferences []UserActivityPreference
	for rows.Next() {
		var preference UserActivityPreference
		if err := rows.Scan(&preference.ID, &preference.UserID, &preference.ActivityID, &preference.Frequency, &preference.FrequencyPeriod, &preference.DaysOfWeek); err != nil {
			return nil, err
		}
		preferences = append(preferences, preference)
	}

	return preferences, nil
}
