package test_activities

import (
	"friendsocial/activities"
	"friendsocial/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCreate(t *testing.T) {
	// Initialize the database connection
	postgres.InitDB()
	defer postgres.CloseDB()

	// Ensure the database connection is successful
	assert.NotNil(t, postgres.DB)

	service := activities.NewService(postgres.DB)

	newActivity := activities.Activity{
		Name:          "Golf",
		Description:   "A game where you hit a golf ball on a golf course.",
		EstimatedTime: "60 Minutes",
		LocationID:    1,
	}

	// Test the Create method
	createdActivity, err := service.Create(newActivity)
	assert.NoError(t, err)

	// Verify that the returned location matches the input
	assert.Equal(t, newActivity.Name, createdActivity.Name)
	assert.Equal(t, newActivity.Description, createdActivity.Description)
	assert.Equal(t, newActivity.EstimatedTime, createdActivity.EstimatedTime)
	assert.Equal(t, newActivity.LocationID, createdActivity.LocationID)
}
