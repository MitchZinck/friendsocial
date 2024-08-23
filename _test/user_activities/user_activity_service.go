package test_user_activities

import (
	"friendsocial/postgres"
	"friendsocial/user_activities"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserActivityService(t *testing.T) {
	// Initialize the database connection
	postgres.InitDB()
	defer postgres.CloseDB()

	// Ensure the database connection is successful
	assert.NotNil(t, postgres.DB)

	service := user_activities.NewService(postgres.DB)

	// Test data
	userID := 1
	activityID := 1

	// Test Create method
	t.Run("Create", func(t *testing.T) {
		newUserActivity := user_activities.UserActivity{
			UserID:     userID,
			ActivityID: activityID,
			IsActive:   true,
		}

		createdUserActivity, err := service.Create(newUserActivity)
		assert.NoError(t, err)
		assert.Equal(t, newUserActivity.UserID, createdUserActivity.UserID)
		assert.Equal(t, newUserActivity.ActivityID, createdUserActivity.ActivityID)
		assert.Equal(t, newUserActivity.IsActive, createdUserActivity.IsActive)
	})

	// Test ActiveUserActivities method
	t.Run("Retrieve", func(t *testing.T) {
		userActivities, err := service.GetActiveUserActivities(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, userActivities)

		// Verify that at least one user activity is associated with the user
		found := false
		for _, ua := range userActivities {
			if ua.ActivityID == activityID {
				found = true
				break
			}
		}
		assert.True(t, found, "Activity should be associated with the user")
	})

}
