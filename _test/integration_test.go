package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"friendsocial/activities"
	"friendsocial/activity_participants"
	"friendsocial/friends"
	"friendsocial/locations"
	"friendsocial/manual_activities"
	"friendsocial/postgres"
	"friendsocial/user_activities"
	"friendsocial/user_activity_preferences"
	"friendsocial/user_availability"
	"friendsocial/users"

	"github.com/jackc/pgx/v4/pgxpool"
)

const baseURL = "http://localhost:8080"

type TestIDs struct {
	UserID                   int
	UserAvailabilityID       int
	FriendID                 int
	LocationIDs              []int
	ActivityID               int
	UserActivityPreferenceID int
	ManualActivityID         int
	ActivityParticipantID    int
	UserActivityID           int
}

func TestIntegration(t *testing.T) {
	// Defer the database wipe to ensure it runs at the end, even if tests fail
	shouldWipeDatabase := true
	ids := TestIDs{}
	defer func() {
		if shouldWipeDatabase {
			// Perform delete tests here, right before wiping the database
			deleteAllEntities(t, ids)
			postgres.InitDB()
			defer postgres.CloseDB()

			err := wipeDatabase(postgres.DB)
			if err != nil {
				t.Errorf("Failed to wipe database: %v", err)
			}
		}
	}()

	// Create a set of locations to use throughout the tests
	ids.LocationIDs = createTestLocations(t)

	// Test user endpoints
	user1 := testUserEndpoints(t, ids.LocationIDs[0])
	ids.UserID = user1.ID

	// Create a second user
	user2 := testUserEndpoints(t, ids.LocationIDs[1])

	// Test friend endpoints
	ids.FriendID = testFriendEndpoints(t, user1.ID, user2.ID)

	// Test user availability endpoints
	ids.UserAvailabilityID = testUserAvailabilityEndpoints(t, user1.ID)

	// Test activity endpoints
	activity := testActivityEndpoints(t, ids.LocationIDs[2])
	ids.ActivityID = activity.ID

	// Test user activity preference endpoints
	ids.UserActivityPreferenceID = testUserActivityPreferenceEndpoints(t, user1.ID, activity.ID)

	// Test manual activity endpoints
	ids.ManualActivityID = testManualActivityEndpoints(t, user1.ID)

	// Test activity participant endpoints
	ids.ActivityParticipantID = testActivityParticipantEndpoints(t, user1.ID, activity.ID)

	// Test user activity endpoints
	ids.UserActivityID = testUserActivityEndpoints(t, user1.ID, activity.ID)
}

// wipeDatabase deletes all data from the tables
func wipeDatabase(db *pgxpool.Pool) error {
	tables := []string{
		"activity_participants",
		"manual_activities",
		"user_activity_preferences",
		"user_availability",
		"friends",
		"user_activities",
		"activities",
		"users",
		"locations",
	}
	fmt.Printf("Wiping tables: %v\n", tables)

	for _, table := range tables {
		_, err := db.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to wipe table %s: %v", table, err)
		}
	}

	return nil
}

func deleteAllEntities(t *testing.T, ids TestIDs) {
	t.Run("Delete Tests", func(t *testing.T) {
		// Delete in reverse order of creation
		testDeleteUserActivity(t, fmt.Sprintf("%d", ids.UserActivityID))
		testDeleteActivityParticipant(t, fmt.Sprintf("%d", ids.ActivityParticipantID))
		testDeleteManualActivity(t, fmt.Sprintf("%d", ids.ManualActivityID))
		testDeleteUserActivityPreference(t, fmt.Sprintf("%d", ids.UserActivityPreferenceID))
		testDeleteActivity(t, fmt.Sprintf("%d", ids.ActivityID))
		testDeleteFriend(t, fmt.Sprintf("%d", ids.UserID), fmt.Sprintf("%d", ids.FriendID))
		testDeleteUserAvailability(t, fmt.Sprintf("%d", ids.UserAvailabilityID))
		testDeleteUser(t, fmt.Sprintf("%d", ids.UserID))
		testDeleteUser(t, fmt.Sprintf("%d", ids.FriendID))
		for _, locationID := range ids.LocationIDs {
			testDeleteLocation(t, fmt.Sprintf("%d", locationID))
		}
	})
}

func testUserEndpoints(t *testing.T, locationID int) users.User {
	var updatedUser users.User

	t.Run("User Endpoints", func(t *testing.T) {
		// Test creating a user
		user := users.User{
			Name:       "Test User",
			Email:      fmt.Sprintf("testuser%d@example.com", time.Now().UnixNano()),
			Password:   "testpassword",
			LocationID: &locationID,
		}
		createdUser := testCreateUser(t, user)

		// Test getting the user
		testGetUser(t, fmt.Sprintf("%d", createdUser.ID))

		// Test full update of the user
		updatedLocationID := locationID + 1
		updatedUserData := users.User{
			Name:       "Updated Test User",
			Email:      fmt.Sprintf("updatedtestuser%d@example.com", time.Now().UnixNano()),
			Password:   "updatedtestpassword",
			LocationID: &updatedLocationID,
		}
		updatedUser = testUpdateUser(t, fmt.Sprintf("%d", createdUser.ID), updatedUserData)

		// Test partial update of the user
		partialUpdate := map[string]interface{}{
			"name": "Partially Updated Test User",
		}
		updatedUser = testPartialUpdateUser(t, fmt.Sprintf("%d", createdUser.ID), partialUpdate)
	})

	return updatedUser
}

// Add this new function to test partial updates
func testPartialUpdateUser(t *testing.T, userID string, updates map[string]interface{}) users.User {
	resp, body := makeRequest(t, "PATCH", fmt.Sprintf("/users/%s", userID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var updatedUser users.User
	err := json.Unmarshal(body, &updatedUser)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check if the partial update was applied correctly
	if updatedUser.Name != updates["name"] {
		t.Fatalf("Partial update failed: expected name %s, got %s", updates["name"], updatedUser.Name)
	}

	return updatedUser
}

func testCreateUser(t *testing.T, user users.User) users.User {
	resp, body := makeRequest(t, "POST", "/users", user)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdUser users.User
	err := json.Unmarshal(body, &createdUser)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdUser.ID == 0 {
		t.Fatalf("Created user ID is 0")
	}

	return createdUser
}

func testGetUser(t *testing.T, userID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/users/%s", userID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUser(t *testing.T, userID string, updates users.User) users.User {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/users/%s", userID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var updatedUser users.User
	err := json.Unmarshal(body, &updatedUser)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	return updatedUser
}

func testDeleteUser(t *testing.T, userID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/users/%s", userID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testUserAvailabilityEndpoints(t *testing.T, userID int) int {
	var createdAvailabilityID int
	t.Run("User Availability Endpoints", func(t *testing.T) {
		availability := user_availability.UserAvailability{
			UserID:      userID,
			DayOfWeek:   "Monday",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(time.Hour * 2),
			IsAvailable: true,
		}

		// Create
		createdAvailability := testCreateUserAvailability(t, availability)
		createdAvailabilityID = createdAvailability.ID

		// Read
		testGetUserAvailability(t, fmt.Sprintf("%d", createdAvailability.ID))

		// Update
		updatedAvailability := user_availability.UserAvailability{
			UserID:      userID,
			DayOfWeek:   "Monday",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(time.Hour * 2),
			IsAvailable: false,
		}
		testUpdateUserAvailability(t, fmt.Sprintf("%d", createdAvailability.ID), updatedAvailability)
	})
	return createdAvailabilityID
}

func testActivityEndpoints(t *testing.T, locationID int) activities.Activity {
	var updatedActivity activities.Activity

	t.Run("Activity Endpoints", func(t *testing.T) {
		activity := activities.Activity{
			Name:          "Test Activity",
			Description:   "This is a test activity",
			EstimatedTime: "120",
			LocationID:    locationID,
		}

		// Create
		createdActivity := testCreateActivity(t, activity)

		// Read
		testGetActivity(t, fmt.Sprintf("%d", createdActivity.ID))

		// Update
		updatedActivity = activities.Activity{
			ID:            createdActivity.ID,
			Name:          "Updated Test Activity",
			Description:   "This is an updated test activity",
			EstimatedTime: "3 hours",
			LocationID:    locationID,
		}
		updatedActivity = testUpdateActivity(t, fmt.Sprintf("%d", createdActivity.ID), updatedActivity)
	})

	return updatedActivity
}

func testUserActivityPreferenceEndpoints(t *testing.T, userID, activityID int) int {
	var createdPreferenceID int
	t.Run("User Activity Preference Endpoints", func(t *testing.T) {
		preference := user_activity_preferences.UserActivityPreference{
			UserID:          userID,
			ActivityID:      activityID,
			Frequency:       2,
			FrequencyPeriod: "week",
		}

		// Create
		createdPreference := testCreateUserActivityPreference(t, preference)
		createdPreferenceID = createdPreference.ID

		// Read
		testGetUserActivityPreference(t, fmt.Sprintf("%d", createdPreference.ID))

		// Update
		updatedPreference := user_activity_preferences.UserActivityPreference{
			ID:              createdPreference.ID,
			UserID:          userID,
			ActivityID:      activityID,
			Frequency:       3,
			FrequencyPeriod: "month",
		}
		testUpdateUserActivityPreference(t, fmt.Sprintf("%d", createdPreference.ID), updatedPreference)
	})
	return createdPreferenceID
}

func testUserActivityEndpoints(t *testing.T, userID, activityID int) int {
	var createdUserActivityID int
	t.Run("User Activity Endpoints", func(t *testing.T) {
		userActivity := user_activities.UserActivity{
			UserID:      userID,
			ActivityID:  activityID,
			IsActive:    true,
			ScheduledAt: time.Now().Add(24 * time.Hour), // Set scheduled_at to 24 hours in the future
		}

		// Create
		createdUserActivity := testCreateUserActivity(t, userActivity)
		createdUserActivityID = createdUserActivity.ID

		// Read
		testGetUserActivity(t, fmt.Sprintf("%d", createdUserActivity.ID))

		// Update
		updatedUserActivity := user_activities.UserActivity{
			ID:          createdUserActivity.ID,
			UserID:      userID,
			ActivityID:  activityID,
			IsActive:    false,
			ScheduledAt: time.Now().Add(48 * time.Hour), // Update scheduled_at to 48 hours in the future
		}
		testUpdateUserActivity(t, fmt.Sprintf("%d", createdUserActivity.ID), updatedUserActivity)
	})
	return createdUserActivityID
}

func testManualActivityEndpoints(t *testing.T, userID int) int {
	var createdManualActivityID int
	t.Run("Manual Activity Endpoints", func(t *testing.T) {
		description := "This is a test manual activity"
		estimatedTime := "1 hour"
		manualActivity := manual_activities.ManualActivity{
			UserID:        userID,
			Name:          "Test Manual Activity",
			Description:   &description,
			EstimatedTime: &estimatedTime,
			ScheduledAt:   time.Now().Add(24 * time.Hour),
			IsActive:      true,
		}

		// Create
		createdManualActivity := testCreateManualActivity(t, manualActivity)
		createdManualActivityID = createdManualActivity.ID

		// Read
		testGetManualActivity(t, fmt.Sprintf("%d", createdManualActivity.ID))

		// Update
		updatedManualActivity := manual_activities.ManualActivity{
			ID:            createdManualActivity.ID,
			UserID:        userID,
			Name:          "Updated Test Manual Activity",
			Description:   &description,
			EstimatedTime: &estimatedTime,
			ScheduledAt:   time.Now().Add(48 * time.Hour),
			IsActive:      false,
		}
		testUpdateManualActivity(t, fmt.Sprintf("%d", createdManualActivity.ID), updatedManualActivity)
	})
	return createdManualActivityID
}

func testFriendEndpoints(t *testing.T, userID1, userID2 int) int {
	var createdFriendID int
	t.Run("Friend Endpoints", func(t *testing.T) {
		friend := friends.Friend{
			UserID:   userID1,
			FriendID: userID2,
		}

		// Create
		createdFriend := testCreateFriend(t, friend)
		createdFriendID = createdFriend.FriendID

		// Read
		testGetFriend(t, fmt.Sprintf("%d", createdFriend.UserID), fmt.Sprintf("%d", createdFriend.FriendID))
	})
	return createdFriendID
}

func testActivityParticipantEndpoints(t *testing.T, userID, activityID int) int {
	var createdParticipantID int
	t.Run("Activity Participant Endpoints", func(t *testing.T) {
		participant := activity_participants.ActivityParticipant{
			UserID:     userID,
			ActivityID: &activityID,
			IsCreator:  true,
			IsActive:   true,
		}

		// Create
		createdParticipant := testCreateActivityParticipant(t, participant)
		createdParticipantID = createdParticipant.ID

		// Read
		testGetActivityParticipant(t, fmt.Sprintf("%d", createdParticipant.ID))

		// Update
		updatedParticipant := activity_participants.ActivityParticipant{
			ID:         createdParticipant.ID,
			UserID:     userID,
			ActivityID: &activityID,
			IsCreator:  false,
			IsActive:   false,
		}
		testUpdateActivityParticipant(t, fmt.Sprintf("%d", createdParticipant.ID), updatedParticipant)
	})
	return createdParticipantID
}

func createTestLocations(t *testing.T) []int {
	locations := []locations.Location{
		{
			Name:      "Test Location 1",
			Address:   "123 Test St",
			City:      "Test City 1",
			State:     "TS1",
			ZipCode:   "12345",
			Country:   "Test Country 1",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		{
			Name:      "Test Location 2",
			Address:   "456 Test Ave",
			City:      "Test City 2",
			State:     "TS2",
			ZipCode:   "67890",
			Country:   "Test Country 2",
			Latitude:  34.0522,
			Longitude: -118.2437,
		},
		{
			Name:      "Test Location 3",
			Address:   "789 Test Blvd",
			City:      "Test City 3",
			State:     "TS3",
			ZipCode:   "13579",
			Country:   "Test Country 3",
			Latitude:  41.8781,
			Longitude: -87.6298,
		},
	}

	var locationIDs []int
	for _, loc := range locations {
		createdLocation := testCreateLocation(t, loc)
		locationIDs = append(locationIDs, createdLocation.ID)

		testGetLocation(t, fmt.Sprintf("%d", createdLocation.ID))

		loc.Name = loc.Name + " Updated"
		testUpdateLocation(t, fmt.Sprintf("%d", createdLocation.ID), loc)
	}

	return locationIDs
}

// Helper functions for each endpoint

func testCreateUserAvailability(t *testing.T, availability user_availability.UserAvailability) user_availability.UserAvailability {
	resp, body := makeRequest(t, "POST", "/user_availability", availability)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v. Response body: %s", resp.Status, string(body))
	}

	var createdAvailability user_availability.UserAvailability
	err := json.Unmarshal(body, &createdAvailability)
	if err != nil {
		t.Fatalf("Failed to parse response: %v. Response body: %s", err, string(body))
	}

	return createdAvailability
}

func testGetUserAvailability(t *testing.T, id string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/user_availability/%s", id), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUserAvailability(t *testing.T, id string, updates user_availability.UserAvailability) {
	resp, _ := makeRequest(t, "PUT", fmt.Sprintf("/user_availability/%s", id), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testDeleteUserAvailability(t *testing.T, id string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/user_availability/%s", id), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateActivity(t *testing.T, activity activities.Activity) activities.Activity {
	resp, body := makeRequest(t, "POST", "/activity", activity)
	if resp.StatusCode != http.StatusCreated {
		t.Logf("Failed to create activity. Activity: %+v", activity)
		t.Fatalf("Expected status Created, got %v. Response body: %s", resp.Status, string(body))
	}

	var createdActivity activities.Activity
	err := json.Unmarshal(body, &createdActivity)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdActivity.ID == 0 {
		t.Fatalf("Created activity has ID 0. Activity: %+v", activity)
	}

	return createdActivity
}

func testGetActivity(t *testing.T, activityID string) {
	resp, body := makeRequest(t, "GET", fmt.Sprintf("/activity/%s", activityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testUpdateActivity(t *testing.T, activityID string, updates activities.Activity) activities.Activity {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/activity/%s", activityID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var updatedActivity activities.Activity
	err := json.Unmarshal(body, &updatedActivity)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	return updatedActivity
}

func testDeleteActivity(t *testing.T, activityID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/activity/%s", activityID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateUserActivityPreference(t *testing.T, preference user_activity_preferences.UserActivityPreference) user_activity_preferences.UserActivityPreference {
	resp, body := makeRequest(t, "POST", "/user_activity_preference", preference)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v. Response body: %s", resp.Status, string(body))
	}

	var createdPreference user_activity_preferences.UserActivityPreference
	err := json.Unmarshal(body, &createdPreference)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdPreference.ID == 0 {
		t.Fatalf("Created user activity preference ID is 0")
	}

	return createdPreference
}

func testGetUserActivityPreference(t *testing.T, preferenceID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/user_activity_preference/%s", preferenceID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUserActivityPreference(t *testing.T, preferenceID string, updates user_activity_preferences.UserActivityPreference) {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/user_activity_preference/%s", preferenceID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testDeleteUserActivityPreference(t *testing.T, preferenceID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/user_activity_preference/%s", preferenceID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateUserActivity(t *testing.T, userActivity user_activities.UserActivity) user_activities.UserActivity {
	resp, body := makeRequest(t, "POST", "/user_activity", userActivity)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdUserActivity user_activities.UserActivity
	err := json.Unmarshal(body, &createdUserActivity)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdUserActivity.ID == 0 {
		t.Fatalf("Created user activity ID is 0")
	}

	return createdUserActivity
}

func testGetUserActivity(t *testing.T, userActivityID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/user_activity/%s", userActivityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUserActivity(t *testing.T, userActivityID string, updates user_activities.UserActivity) {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/user_activity/%s", userActivityID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testDeleteUserActivity(t *testing.T, userActivityID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/user_activity/%s", userActivityID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateManualActivity(t *testing.T, manualActivity manual_activities.ManualActivity) manual_activities.ManualActivity {
	resp, body := makeRequest(t, "POST", "/manual_activity", manualActivity)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v. Response body: %s", resp.Status, string(body))
	}

	var createdManualActivity manual_activities.ManualActivity
	err := json.Unmarshal(body, &createdManualActivity)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdManualActivity.ID == 0 {
		t.Fatalf("Created manual activity ID is 0")
	}

	return createdManualActivity
}

func testGetManualActivity(t *testing.T, manualActivityID string) {
	resp, body := makeRequest(t, "GET", fmt.Sprintf("/manual_activity/%s", manualActivityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testUpdateManualActivity(t *testing.T, manualActivityID string, updates manual_activities.ManualActivity) {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/manual_activity/%s", manualActivityID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testDeleteManualActivity(t *testing.T, manualActivityID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/manual_activity/%s", manualActivityID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateFriend(t *testing.T, friend friends.Friend) friends.Friend {
	resp, body := makeRequest(t, "POST", "/friend", friend)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v. Response body: %s", resp.Status, string(body))
	}

	var createdFriend friends.Friend
	err := json.Unmarshal(body, &createdFriend)
	if err != nil {
		t.Fatalf("Failed to parse response: %v. Response body: %s", err, string(body))
	}

	return createdFriend
}

func testGetFriend(t *testing.T, userID, friendID string) {
	resp, body := makeRequest(t, "GET", fmt.Sprintf("/friend/%s/%s", userID, friendID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testDeleteFriend(t *testing.T, userID, friendID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/friend/%s/%s", userID, friendID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateActivityParticipant(t *testing.T, participant activity_participants.ActivityParticipant) activity_participants.ActivityParticipant {
	resp, body := makeRequest(t, "POST", "/activity_participant", participant)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdParticipant activity_participants.ActivityParticipant
	err := json.Unmarshal(body, &createdParticipant)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdParticipant.ID == 0 {
		t.Fatalf("Created activity participant ID is 0")
	}

	return createdParticipant
}

func testGetActivityParticipant(t *testing.T, participantID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/activity_participant/%s", participantID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateActivityParticipant(t *testing.T, participantID string, updates activity_participants.ActivityParticipant) {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/activity_participant/%s", participantID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v. Response body: %s", resp.Status, string(body))
	}
}

func testDeleteActivityParticipant(t *testing.T, participantID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/activity_participant/%s", participantID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateLocation(t *testing.T, location locations.Location) locations.Location {
	resp, body := makeRequest(t, "POST", "/location", location)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdLocation locations.Location
	err := json.Unmarshal(body, &createdLocation)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdLocation.ID == 0 {
		t.Logf("Created location ID is 0. Response body: %s", string(body))
		t.Fatalf("Created location ID is 0")
	}

	return createdLocation
}

func testGetLocation(t *testing.T, locationID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/location/%s", locationID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateLocation(t *testing.T, locationID string, updates locations.Location) locations.Location {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/location/%s", locationID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var updatedLocation locations.Location
	err := json.Unmarshal(body, &updatedLocation)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	return updatedLocation
}

func testDeleteLocation(t *testing.T, locationID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/location/%s", locationID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func makeRequest(t *testing.T, method, path string, body interface{}) (*http.Response, []byte) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	return resp, respBody
}
