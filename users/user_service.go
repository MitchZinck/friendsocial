package users

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

func (userService *Service) Create(user User) (User, error) {
	userService.Lock()
	defer userService.Unlock()

	_, err := userService.db.Exec(
		context.Background(),
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, user.Password,
	)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return user, nil
}

func (userService *Service) ReadAll() ([]User, error) {
	userService.Lock()
	defer userService.Unlock()

	rows, err := userService.db.Query(context.Background(), "SELECT id, name, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (userService *Service) Read(id string) (User, bool, error) {
	userService.Lock()
	defer userService.Unlock()

	var user User
	err := userService.db.QueryRow(context.Background(), "SELECT id, name, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, err
	}

	return user, true, nil
}

func (userService *Service) Update(id string, user User) (User, bool, error) {
	userService.Lock()
	defer userService.Unlock()

	cmdTag, err := userService.db.Exec(context.Background(), "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4", user.Name, user.Email, user.Password, id)
	if err != nil {
		return User{}, false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return User{}, false, nil
	}

	return user, true, nil
}

func (userService *Service) Delete(id string) (bool, error) {
	userService.Lock()
	defer userService.Unlock()

	cmdTag, err := userService.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
