package activity_participants

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ActivityParticipant struct {
	ID                  int `json:"id"`
	UserID              int `json:"user_id"`
	ScheduledActivityID int `json:"scheduled_activity_id"`
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
		(user_id, scheduled_activity_id) 
		VALUES ($1, $2) 
		RETURNING id`,
		participant.UserID, participant.ScheduledActivityID,
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
		"SELECT id, user_id, scheduled_activity_id FROM activity_participants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		err := rows.Scan(
			&participant.ID, &participant.UserID, &participant.ScheduledActivityID)
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
		`SELECT id, user_id, scheduled_activity_id 
		FROM activity_participants WHERE id = $1`, id).Scan(
		&participant.ID, &participant.UserID, &participant.ScheduledActivityID)

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
		SET user_id = $1, scheduled_activity_id = $2
		WHERE id = $3`,
		participant.UserID, participant.ScheduledActivityID, id)

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
