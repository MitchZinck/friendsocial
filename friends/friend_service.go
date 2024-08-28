package friends

import (
	"context"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Friend struct {
	UserID    int    `json:"user_id"`
	FriendID  int    `json:"friend_id"`
	CreatedAt string `json:"created_at"`
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

// Creates a friendship between two users
func (friendService *Service) Create(userID string, friendID string) (Friend, error) {
	friendService.Lock()
	defer friendService.Unlock()

	_, err := friendService.db.Exec(
		context.Background(),
		"INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)",
		userID, friendID,
	)
	if err != nil {
		return Friend{}, err
	}

	userIDStr, _ := strconv.Atoi(userID)
	friendIDStr, _ := strconv.Atoi(friendID)

	// Return the friendship details without querying again
	return Friend{
		UserID:    userIDStr,
		FriendID:  friendIDStr,
		CreatedAt: "", // We assume created_at is automatically handled by the database.
	}, nil
}

// Removes a friendship between two users
func (friendService *Service) Delete(userID string, friendID string) (bool, error) {
	friendService.Lock()
	defer friendService.Unlock()

	cmdTag, err := friendService.db.Exec(
		context.Background(),
		"DELETE FROM friends WHERE user_id = $1 AND friend_id = $2",
		userID, friendID,
	)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

// Retrieves all friends of a given user
func (friendService *Service) ReadAll(userID string) ([]Friend, error) {
	friendService.Lock()
	defer friendService.Unlock()

	rows, err := friendService.db.Query(
		context.Background(),
		"SELECT user_id, friend_id, created_at::text FROM friends WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []Friend
	for rows.Next() {
		var friend Friend
		if err := rows.Scan(&friend.UserID, &friend.FriendID, &friend.CreatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// Checks if two users are friends
func (friendService *Service) Read(userID string, friendID string) (bool, error) {
	friendService.Lock()
	defer friendService.Unlock()

	var exists bool
	err := friendService.db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM friends WHERE user_id = $1 AND friend_id = $2)",
		userID, friendID,
	).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
