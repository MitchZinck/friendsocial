package user_activity_preferences_participants

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserActivityPreferenceParticipant struct {
	ID                       int `json:"id"`
	UserActivityPreferenceID int `json:"user_activity_preference_id"`
	UserID                   int `json:"user_id"`
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

// Implement Create, ReadAll, Read, Update, Delete methods similar to UserActivityPreference Service

func (s *Service) Create(participant UserActivityPreferenceParticipant) (UserActivityPreferenceParticipant, error) {
	s.Lock()
	defer s.Unlock()

	query := `
		INSERT INTO user_activity_preferences_participants (user_activity_preference_id, user_id)
		VALUES ($1, $2)
		RETURNING id, user_activity_preference_id, user_id
	`

	err := s.db.QueryRow(context.Background(), query, participant.UserActivityPreferenceID, participant.UserID).Scan(&participant.ID, &participant.UserActivityPreferenceID, &participant.UserID)
	if err != nil {
		return UserActivityPreferenceParticipant{}, err
	}

	return participant, nil
}

func (s *Service) ReadAll() ([]UserActivityPreferenceParticipant, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(), "SELECT id, user_activity_preference_id, user_id FROM user_activity_preferences_participants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []UserActivityPreferenceParticipant
	for rows.Next() {
		var participant UserActivityPreferenceParticipant
		if err := rows.Scan(&participant.ID, &participant.UserActivityPreferenceID, &participant.UserID); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func (s *Service) Read(id string) (UserActivityPreferenceParticipant, bool, error) {
	s.Lock()
	defer s.Unlock()

	query := `
		SELECT id, user_activity_preference_id, user_id
		FROM user_activity_preferences_participants
		WHERE id = $1
	`

	var participant UserActivityPreferenceParticipant
	err := s.db.QueryRow(context.Background(), query, id).Scan(&participant.ID, &participant.UserActivityPreferenceID, &participant.UserID)
	if err != nil {
		return UserActivityPreferenceParticipant{}, false, err
	}

	return participant, true, nil
}

func (s *Service) Update(id string, participant UserActivityPreferenceParticipant) (UserActivityPreferenceParticipant, bool, error) {
	s.Lock()
	defer s.Unlock()

	query := `
		UPDATE user_activity_preferences_participants
		SET user_activity_preference_id = $1, user_id = $2
		WHERE id = $3
		RETURNING id, user_activity_preference_id, user_id
	`

	err := s.db.QueryRow(context.Background(), query, participant.UserActivityPreferenceID, participant.UserID, id).Scan(&participant.ID, &participant.UserActivityPreferenceID, &participant.UserID)
	if err != nil {
		return UserActivityPreferenceParticipant{}, false, err
	}

	return participant, true, nil
}

func (s *Service) Delete(id string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	query := `
		DELETE FROM user_activity_preferences_participants 
		WHERE id = $1
	`

	result, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (s *Service) ReadByPreferenceID(preferenceID string) ([]UserActivityPreferenceParticipant, error) {
	s.Lock()
	defer s.Unlock()

	rows, err := s.db.Query(context.Background(), "SELECT id, user_activity_preference_id, user_id FROM user_activity_preferences_participants WHERE user_activity_preference_id = $1", preferenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []UserActivityPreferenceParticipant
	for rows.Next() {
		var participant UserActivityPreferenceParticipant
		if err := rows.Scan(&participant.ID, &participant.UserActivityPreferenceID, &participant.UserID); err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}
