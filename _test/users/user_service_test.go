package test_users

import (
	"friendsocial/postgres"
	"friendsocial/users"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	// Initialize the database connection
	postgres.InitDB()
	defer postgres.CloseDB()

	// Ensure the database connection is successful
	assert.NotNil(t, postgres.DB)

	service := users.NewService(postgres.DB)

	newUser := users.User{
		Name:     "Mitchell Zinck",
		Email:    "mitchellfzinck@gmail.com",
		Password: "test123",
	}

	// Test the Create method
	createdUser, err := service.Create(newUser)
	assert.NoError(t, err)

	// Verify that the returned location matches the input
	assert.Equal(t, newUser.Name, createdUser.Name)
	assert.Equal(t, newUser.Email, createdUser.Email)
	// assert.Equal(t, newUser.Password, createdUser.Password)
}
