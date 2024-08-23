package activity_participants

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ActivityParticipant struct {
	ID               int  `json:"id"`
	UserID           int  `json:"user_id"`
	ActivityID       *int `json:"activity_id,omitempty"`
	ManualActivityID *int `json:"manual_activity_id,omitempty"`
	IsCreator        bool `json:"is_creator"`
	IsActive         bool `json:"is_active"`
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

func (s *Service) Create(participant ActivityParticipant) (ActivityParticipant, error) {
	s.Lock()
	defer s.Unlock()

	err := s.db.QueryRow(
		context.Background(),
		`INSERT INTO activity_participants 
		(user_id, activity_id, manual_activity_id, is_creator, is_active) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`,
		participant.UserID, participant.ActivityID, participant.ManualActivityID, participant.IsCreator, participant.IsActive,
	).Scan(&participant.ID)

	if err != nil {
		return ActivityParticipant{}, err
	}

	return participant, nil
}

func (s *Service) ReadAll() ([]ActivityParticipant, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(),
		"SELECT id, user_id, activity_id, manual_activity_id, is_creator, is_active FROM activity_participants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		err := rows.Scan(
			&participant.ID, &participant.UserID, &participant.ActivityID, &participant.ManualActivityID,
			&participant.IsCreator, &participant.IsActive)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (s *Service) Read(id string) (ActivityParticipant, bool, error) {
	s.Lock()
	defer s.Unlock()

	var participant ActivityParticipant
	err := s.db.QueryRow(
		context.Background(),
		`SELECT id, user_id, activity_id, manual_activity_id, is_creator, is_active 
		FROM activity_participants WHERE id = $1`, id).Scan(
		&participant.ID, &participant.UserID, &participant.ActivityID, &participant.ManualActivityID,
		&participant.IsCreator, &participant.IsActive)

	if err != nil {
		if err == pgx.ErrNoRows {
			return ActivityParticipant{}, false, nil
		}
		return ActivityParticipant{}, false, err
	}

	return participant, true, nil
}

func (s *Service) Update(id string, participant ActivityParticipant) (ActivityParticipant, bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(
		context.Background(),
		`UPDATE activity_participants 
		SET user_id = $1, activity_id = $2, manual_activity_id = $3, 
		is_creator = $4, is_active = $5 
		WHERE id = $6`,
		participant.UserID, participant.ActivityID, participant.ManualActivityID, participant.IsCreator, participant.IsActive, id)

	if err != nil {
		return ActivityParticipant{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return ActivityParticipant{}, false, nil
	}

	return participant, true, nil
}

func (s *Service) Delete(id string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(context.Background(), "DELETE FROM activity_participants WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
