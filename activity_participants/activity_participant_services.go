package activity_participants

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

type ActivityParticipant struct {
	ID                  int    `json:"id"`
	UserID              int    `json:"user_id"`
	ScheduledActivityID int    `json:"scheduled_activity_id"`
	InviteStatus        string `json:"invite_status"`
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
		"SELECT id, user_id, scheduled_activity_id, invite_status FROM activity_participants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		err := rows.Scan(
			&participant.ID, &participant.UserID, &participant.ScheduledActivityID, &participant.InviteStatus)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (s *Service) Read(ids []string) ([]ActivityParticipant, error) {
	s.Lock()
	defer s.Unlock()

	if len(ids) == 0 {
		return []ActivityParticipant{}, nil
	}

	query := "SELECT id, user_id, scheduled_activity_id, invite_status FROM activity_participants WHERE id = ANY($1)"
	rows, err := s.db.Query(context.Background(), query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		if err := rows.Scan(&participant.ID, &participant.UserID, &participant.ScheduledActivityID, &participant.InviteStatus); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

func (s *Service) Update(id string, participant ActivityParticipant) (ActivityParticipant, bool, error) {
	s.Lock()
	defer s.Unlock()

	cmdTag, err := s.db.Exec(
		context.Background(),
		`UPDATE activity_participants 
		SET user_id = $1, scheduled_activity_id = $2, invite_status = $3
		WHERE id = $4`,
		participant.UserID, participant.ScheduledActivityID, participant.InviteStatus, id)

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

func (s *Service) GetActivitiesByUserID(userID string) ([]ActivityParticipant, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(),
		`SELECT id, user_id, scheduled_activity_id, invite_status 
         FROM activity_participants 
         WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		err := rows.Scan(&participant.ID, &participant.UserID, &participant.ScheduledActivityID, &participant.InviteStatus)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (s *Service) GetParticipantsByScheduledActivityID(scheduledActivityID []string) ([]ActivityParticipant, error) {
	s.Lock()
	defer s.Unlock()

	query := `SELECT id, user_id, scheduled_activity_id, invite_status 
         FROM activity_participants 
         WHERE scheduled_activity_id = ANY($1)`
	rows, err := s.db.Query(context.Background(), query, pq.Array(scheduledActivityID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []ActivityParticipant
	for rows.Next() {
		var participant ActivityParticipant
		if err := rows.Scan(&participant.ID, &participant.UserID, &participant.ScheduledActivityID, &participant.InviteStatus); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}
