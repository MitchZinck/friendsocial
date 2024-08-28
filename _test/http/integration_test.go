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
	"friendsocial/activity_locations"
	"friendsocial/activity_participants"
	"friendsocial/friends"
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
	ActivityLocationID       int
	ActivityID               int
	UserActivityPreferenceID int
	ManualActivityID         int
	ActivityParticipantID    int
	UserActivityID           int
}

func TestIntegration(t *testing.T) {
	// Defer the database wipe to ensure it runs at the end, even if tests fail
	shouldWipeDatabase := true
	defer func() {
		if shouldWipeDatabase {
			postgres.InitDB()
			defer postgres.CloseDB()

			err := wipeDatabase(postgres.DB)
			if err != nil {
				t.Errorf("Failed to wipe database: %v", err)
			}
		}
	}()

	ids := TestIDs{}

	// Test user endpoints
	user1 := testUserEndpoints(t)
	ids.UserID = user1.ID

	// Create a second user
	user2 := testUserEndpoints(t)

	// Test friend endpoints
	ids.FriendID = testFriendEndpoints(t, user1.ID, user2.ID)

	// Test user availability endpoints
	ids.UserAvailabilityID = testUserAvailabilityEndpoints(t, user1.ID)

	// Test activity location endpoints
	activityLocation := testActivityLocationEndpoints(t)
	ids.ActivityLocationID = activityLocation.ID

	// Test activity endpoints
	activity := testActivityEndpoints(t, activityLocation.ID)
	ids.ActivityID = activity.ID

	// Test user activity preference endpoints
	ids.UserActivityPreferenceID = testUserActivityPreferenceEndpoints(t, user1.ID, activity.ID)

	// Test manual activity endpoints
	ids.ManualActivityID = testManualActivityEndpoints(t, user1.ID)

	// Test activity participant endpoints
	ids.ActivityParticipantID = testActivityParticipantEndpoints(t, user1.ID, activity.ID)

	// Test user activity endpoints
	ids.UserActivityID = testUserActivityEndpoints(t, user1.ID, activity.ID)

	// Perform delete tests here, right before wiping the database
	deleteAllEntities(t, ids)
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
		"activity_locations",
		"users",
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
		testDeleteActivityLocation(t, fmt.Sprintf("%d", ids.ActivityLocationID))
		testDeleteFriend(t, fmt.Sprintf("%d", ids.UserID), fmt.Sprintf("%d", ids.FriendID))
		testDeleteUserAvailability(t, fmt.Sprintf("%d", ids.UserAvailabilityID))
		testDeleteUser(t, fmt.Sprintf("%d", ids.UserID))
		testDeleteUser(t, fmt.Sprintf("%d", ids.FriendID))
	})
}

func testUserEndpoints(t *testing.T) users.User {
	var updatedUser users.User

	t.Run("User Endpoints", func(t *testing.T) {
		// Test creating a user
		user := users.User{
			Name:     "Test User",
			Email:    fmt.Sprintf("testuser%d@example.com", time.Now().UnixNano()),
			Password: "testpassword",
		}
		createdUser := testCreateUser(t, user)

		// Test getting the user
		testGetUser(t, fmt.Sprintf("%d", createdUser.ID))

		// Test full update of the user
		updatedUserData := users.User{
			Name:     "Updated Test User",
			Email:    fmt.Sprintf("updatedtestuser%d@example.com", time.Now().UnixNano()),
			Password: "updatedtestpassword",
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
			Frequency: 3,
		}
		testUpdateUserActivityPreference(t, fmt.Sprintf("%d", createdPreference.ID), updatedPreference)
	})
	return createdPreferenceID
}

func testUserActivityEndpoints(t *testing.T, userID, activityID int) int {
	var createdUserActivityID int
	t.Run("User Activity Endpoints", func(t *testing.T) {
		userActivity := user_activities.UserActivity{
			UserID:     userID,
			ActivityID: activityID,
			IsActive:   true,
		}

		// Create
		createdUserActivity := testCreateUserActivity(t, userActivity)
		createdUserActivityID = createdUserActivity.ID

		// Read
		testGetUserActivity(t, fmt.Sprintf("%d", createdUserActivity.ID))

		// Update
		updatedUserActivity := user_activities.UserActivity{
			IsActive: false,
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
			IsActive:      true,
		}

		// Create
		createdManualActivity := testCreateManualActivity(t, manualActivity)
		createdManualActivityID = createdManualActivity.ID

		// Read
		testGetManualActivity(t, fmt.Sprintf("%d", createdManualActivity.ID))

		// Update
		updatedManualActivity := manual_activities.ManualActivity{
			Name: "Updated Test Manual Activity",
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
			IsActive: false,
		}
		testUpdateActivityParticipant(t, fmt.Sprintf("%d", createdParticipant.ID), updatedParticipant)
	})
	return createdParticipantID
}

func testActivityLocationEndpoints(t *testing.T) activity_locations.ActivityLocation {
	var updatedLocation activity_locations.ActivityLocation

	t.Run("Activity Location Endpoints", func(t *testing.T) {
		location := activity_locations.ActivityLocation{
			Name:      "Test Location",
			Address:   "123 Test St",
			City:      "Test City",
			State:     "TS",
			ZipCode:   "12345",
			Country:   "Test Country",
			Latitude:  40.7128,
			Longitude: -74.0060,
		}

		// Create
		createdLocation := testCreateActivityLocation(t, location)

		// Read
		testGetActivityLocation(t, fmt.Sprintf("%d", createdLocation.ID))

		// Update
		updatedLocation = activity_locations.ActivityLocation{
			Name:      "Updated Test Location",
			Address:   "456 Updated St",
			City:      "Updated City",
			State:     "US",
			ZipCode:   "67890",
			Country:   "Updated Country",
			Latitude:  34.0522,
			Longitude: -118.2437,
		}
		updatedLocation = testUpdateActivityLocation(t, fmt.Sprintf("%d", createdLocation.ID), updatedLocation)
	})

	return updatedLocation
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
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdActivity activities.Activity
	err := json.Unmarshal(body, &createdActivity)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdActivity.ID == 0 {
		t.Fatalf("Created activity ID is 0")
	}

	return createdActivity
}

func testGetActivity(t *testing.T, activityID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/activity/%s", activityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
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
	resp, body := makeRequest(t, "POST", "/user_activity_preferences", preference)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
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
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/user_activity_preferences/%s", preferenceID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUserActivityPreference(t *testing.T, preferenceID string, updates user_activity_preferences.UserActivityPreference) {
	resp, _ := makeRequest(t, "PUT", fmt.Sprintf("/user_activity_preferences/%s", preferenceID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testDeleteUserActivityPreference(t *testing.T, preferenceID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/user_activity_preferences/%s", preferenceID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateUserActivity(t *testing.T, userActivity user_activities.UserActivity) user_activities.UserActivity {
	resp, body := makeRequest(t, "POST", "/user_activities", userActivity)
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
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/user_activities/%s", userActivityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateUserActivity(t *testing.T, userActivityID string, updates user_activities.UserActivity) {
	resp, _ := makeRequest(t, "PUT", fmt.Sprintf("/user_activities/%s", userActivityID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testDeleteUserActivity(t *testing.T, userActivityID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/user_activities/%s", userActivityID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateManualActivity(t *testing.T, manualActivity manual_activities.ManualActivity) manual_activities.ManualActivity {
	resp, body := makeRequest(t, "POST", "/manual_activities", manualActivity)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
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
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/manual_activities/%s", manualActivityID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateManualActivity(t *testing.T, manualActivityID string, updates manual_activities.ManualActivity) {
	resp, _ := makeRequest(t, "PUT", fmt.Sprintf("/manual_activities/%s", manualActivityID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testDeleteManualActivity(t *testing.T, manualActivityID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/manual_activities/%s", manualActivityID), nil)
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
	resp, body := makeRequest(t, "POST", "/activity_participants", participant)
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
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/activity_participants/%s", participantID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateActivityParticipant(t *testing.T, participantID string, updates activity_participants.ActivityParticipant) {
	resp, _ := makeRequest(t, "PUT", fmt.Sprintf("/activity_participants/%s", participantID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testDeleteActivityParticipant(t *testing.T, participantID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/activity_participants/%s", participantID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %v", resp.Status)
	}
}

func testCreateActivityLocation(t *testing.T, location activity_locations.ActivityLocation) activity_locations.ActivityLocation {
	resp, body := makeRequest(t, "POST", "/activity_locations", location)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %v", resp.Status)
	}

	var createdLocation activity_locations.ActivityLocation
	err := json.Unmarshal(body, &createdLocation)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if createdLocation.ID == 0 {
		t.Fatalf("Created activity location ID is 0")
	}

	return createdLocation
}

func testGetActivityLocation(t *testing.T, locationID string) {
	resp, _ := makeRequest(t, "GET", fmt.Sprintf("/activity_locations/%s", locationID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}
}

func testUpdateActivityLocation(t *testing.T, locationID string, updates activity_locations.ActivityLocation) activity_locations.ActivityLocation {
	resp, body := makeRequest(t, "PUT", fmt.Sprintf("/activity_locations/%s", locationID), updates)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", resp.Status)
	}

	var updatedLocation activity_locations.ActivityLocation
	err := json.Unmarshal(body, &updatedLocation)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	return updatedLocation
}

func testDeleteActivityLocation(t *testing.T, locationID string) {
	resp, _ := makeRequest(t, "DELETE", fmt.Sprintf("/activity_locations/%s", locationID), nil)
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
