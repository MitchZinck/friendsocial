package user_availability

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserAvailability struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	DayOfWeek   string    `json:"day_of_week"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsAvailable bool      `json:"is_available"`
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

func (s *Service) Create(availability UserAvailability) (UserAvailability, error) {
	s.Lock()
	defer s.Unlock()

	err := s.db.QueryRow(
		context.Background(),
		`INSERT INTO user_availability (user_id, day_of_week, start_time, end_time, is_available) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		availability.UserID, availability.DayOfWeek, availability.StartTime, availability.EndTime, availability.IsAvailable,
	).Scan(&availability.ID)

	if err != nil {
		return UserAvailability{}, err
	}

	return availability, nil
}

func (s *Service) ReadAll() ([]UserAvailability, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(), "SELECT id, user_id, day_of_week, start_time, end_time, is_available FROM user_availability")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []UserAvailability
	for rows.Next() {
		var availability UserAvailability
		if err := rows.Scan(&availability.ID, &availability.UserID, &availability.DayOfWeek, &availability.StartTime, &availability.EndTime, &availability.IsAvailable); err != nil {
			return nil, err
		}
		availabilities = append(availabilities, availability)
	}

	return availabilities, nil
}

func (s *Service) Read(id int) (UserAvailability, bool, error) {
	s.Lock()
	defer s.Unlock()

	var availability UserAvailability
	err := s.db.QueryRow(context.Background(),
		"SELECT id, user_id, day_of_week, start_time, end_time, is_available FROM user_availability WHERE id = $1",
		id,
	).Scan(&availability.ID, &availability.UserID, &availability.DayOfWeek, &availability.StartTime, &availability.EndTime, &availability.IsAvailable)

	if err != nil {
		if err == pgx.ErrNoRows {
			return UserAvailability{}, false, nil
		}
		return UserAvailability{}, false, err
	}

	return availability, true, nil
}

func (s *Service) Update(id int, availability UserAvailability) (UserAvailability, bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(
		context.Background(),
		`UPDATE user_availability 
		 SET user_id = $1, day_of_week = $2, start_time = $3, end_time = $4, is_available = $5 
		 WHERE id = $6`,
		availability.UserID, availability.DayOfWeek, availability.StartTime, availability.EndTime, availability.IsAvailable, id,
	)

	if err != nil {
		return UserAvailability{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return UserAvailability{}, false, nil
	}

	return availability, true, nil
}

func (s *Service) Delete(id int) (bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(context.Background(), "DELETE FROM user_availability WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
