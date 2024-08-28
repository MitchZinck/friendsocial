package main

import (
	"context"
	"fmt"
	"friendsocial/activity_locations"
	"friendsocial/friends"
	"friendsocial/manual_activities"
	"friendsocial/postgres"
	"friendsocial/user_availability"
	"friendsocial/users"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestSetupAndPopulateDatabase(t *testing.T) {
	// Initialize the database connection
	postgres.InitDB()
	defer postgres.CloseDB()

	shouldWipeDatabase := true

	// Defer the database wipe to ensure it runs at the end, even if tests fail
	defer func() {
		if shouldWipeDatabase {
			err := wipeDatabase(postgres.DB)
			if err != nil {
				t.Errorf("Failed to wipe database: %v", err)
			}
		}
	}()

	// Create services
	locationService := activity_locations.NewService(postgres.DB)
	userService := users.NewService(postgres.DB)
	availabilityService := user_availability.NewService(postgres.DB)
	manualActivityService := manual_activities.NewService(postgres.DB)
	friendService := friends.NewService(postgres.DB)

	// Test data
	locations := []activity_locations.ActivityLocation{
		{Name: "Central Park", Address: "Central Park", City: "New York", State: "NY", Country: "USA"},
		{Name: "Golden Gate Park", Address: "Golden Gate Park", City: "San Francisco", State: "CA", Country: "USA"},
	}

	morningJog := "Morning jog in the park"
	morningJogEstimatedTime := "30"
	afternoonYoga := "Outdoor yoga session"
	afternoonYogaEstimatedTime := "60"
	activities := []manual_activities.ManualActivity{
		{Name: "Jogging", Description: &morningJog, EstimatedTime: &morningJogEstimatedTime},
		{Name: "Yoga", Description: &afternoonYoga, EstimatedTime: &afternoonYogaEstimatedTime},
	}

	users := []users.User{
		{Name: "Alice", Email: "alice@example.com", Password: "password123"},
		{Name: "Bob", Email: "bob@example.com", Password: "password456"},
	}

	// Run subtests
	t.Run("Create Locations", func(t *testing.T) {
		createLocations(t, locationService, locations)
	})

	t.Run("Create Users", func(t *testing.T) {
		createUsers(t, userService, &users)
	})

	t.Run("Set User Availabilities", func(t *testing.T) {
		setUserAvailabilities(t, availabilityService, users)
	})

	t.Run("Create Manual Activities", func(t *testing.T) {
		createManualActivities(t, manualActivityService, activities, users)
	})

	t.Run("Create Friend Relationship", func(t *testing.T) {
		createFriendship(t, friendService, users[0].ID, users[1].ID)
	})

	t.Run("Verify Data", func(t *testing.T) {
		verifyUsers(t, userService, users)
		verifyAvailabilities(t, availabilityService, users)
		verifyFriendship(t, friendService, users[0].ID)
	})
}

func createLocations(t *testing.T, service *activity_locations.Service, locations []activity_locations.ActivityLocation) {
	for _, loc := range locations {
		_, err := service.Create(loc)
		if err != nil {
			t.Fatalf("Failed to create location: %v", err)
		}
	}
}

func createUsers(t *testing.T, service *users.Service, users *[]users.User) {
	for i := range *users {
		createdUser, err := service.Create((*users)[i])
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		(*users)[i] = createdUser
	}
}

func setUserAvailabilities(t *testing.T, service *user_availability.Service, users []users.User) {
	for _, user := range users {
		avail := user_availability.UserAvailability{
			UserID:      user.ID,
			DayOfWeek:   "Monday",
			StartTime:   time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:     time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC),
			IsAvailable: true,
		}
		_, err := service.Create(avail)
		if err != nil {
			t.Fatalf("Failed to create availability: %v", err)
		}
	}
}

func createManualActivities(t *testing.T, service *manual_activities.Service, activities []manual_activities.ManualActivity, users []users.User) {
	for i, activity := range activities {
		activity.UserID = users[i%2].ID // Alternate between users
		_, err := service.Create(activity)
		if err != nil {
			t.Fatalf("Failed to create manual activity: %v", err)
		}
	}
}

func createFriendship(t *testing.T, service *friends.Service, userID, friendID int) {
	_, err := service.Create(fmt.Sprintf("%d", userID), fmt.Sprintf("%d", friendID))
	if err != nil {
		t.Fatalf("Failed to create friendship: %v", err)
	}
}

func verifyUsers(t *testing.T, service *users.Service, users []users.User) {
	for _, user := range users {
		retrievedUser, found, err := service.Read(fmt.Sprintf("%d", user.ID))
		if err != nil {
			t.Errorf("Failed to retrieve user %d: %v", user.ID, err)
		}
		if !found {
			t.Errorf("User %d not found", user.ID)
		}
		if retrievedUser.ID != user.ID {
			t.Errorf("Retrieved user ID %d does not match expected ID %d", retrievedUser.ID, user.ID)
		}
	}
}

func verifyAvailabilities(t *testing.T, service *user_availability.Service, users []users.User) {
	for _, user := range users {
		_, err := service.ReadAll(fmt.Sprintf("%d", user.ID))
		if err != nil {
			t.Errorf("Failed to retrieve availability for user %d: %v", user.ID, err)
		}
	}
}

func verifyFriendship(t *testing.T, service *friends.Service, userID int) {
	friends, err := service.ReadAll(fmt.Sprintf("%d", userID))
	if err != nil || len(friends) == 0 {
		t.Errorf("Failed to retrieve friendship for user %d: %v", userID, err)
	}
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
