package main

import (
	"net/http"

	"database/sql"
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
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	postgres.InitDB()
	defer postgres.CloseDB()

	mux := http.NewServeMux()

	userServices := users.NewService(postgres.DB)
	userManager := users.NewUserHTTPHandler(userServices)

	mux.HandleFunc("POST /users", userManager.HandleHTTPPost)
	mux.HandleFunc("GET /users", userManager.HandleHTTPGet)
	mux.HandleFunc("GET /users/{id}", userManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /users/{id}", userManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /users/{id}", userManager.HandleHTTPDelete)
	mux.HandleFunc("PATCH /users/{id}", userManager.HandleHTTPPatch)

	// User availability services and handlers
	availabilityService := user_availability.NewService(postgres.DB)
	availabilityManager := user_availability.NewUserAvailabilityHTTPHandler(availabilityService)

	mux.HandleFunc("POST /user_availability", availabilityManager.HandleHTTPPost)
	mux.HandleFunc("GET /user_availability", availabilityManager.HandleHTTPGet)
	mux.HandleFunc("GET /user_availability/{id}", availabilityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /user_availability/{id}", availabilityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /user_availability/{id}", availabilityManager.HandleHTTPDelete)

	userActivityPreferenceService := user_activity_preferences.NewService(postgres.DB)
	userActivityPreferenceManager := user_activity_preferences.NewUserActivityPreferenceHTTPHandler(userActivityPreferenceService)

	mux.HandleFunc("POST /user_activity_preference", userActivityPreferenceManager.HandleHTTPPost)
	mux.HandleFunc("GET /user_activity_preferences", userActivityPreferenceManager.HandleHTTPGet)
	mux.HandleFunc("GET /user_activity_preference/{id}", userActivityPreferenceManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /user_activity_preference/{id}", userActivityPreferenceManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /user_activity_preference/{id}", userActivityPreferenceManager.HandleHTTPDelete)

	userActivityService := user_activities.NewService(postgres.DB)
	userActivityManager := user_activities.NewUserActivityHTTPHandler(userActivityService)

	mux.HandleFunc("POST /user_activity", userActivityManager.HandleHTTPPost)
	mux.HandleFunc("GET /user_activities", userActivityManager.HandleHTTPGet)
	mux.HandleFunc("GET /user_activity/{id}", userActivityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /user_activity/{id}", userActivityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /user_activity/{id}", userActivityManager.HandleHTTPDelete)

	manualActivityService := manual_activities.NewService(postgres.DB)
	manualActivityManager := manual_activities.NewManualActivityHTTPHandler(manualActivityService)

	mux.HandleFunc("POST /manual_activity", manualActivityManager.HandleHTTPPost)
	mux.HandleFunc("GET /manual_activities", manualActivityManager.HandleHTTPGet)
	mux.HandleFunc("GET /manual_activity/{id}", manualActivityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /manual_activity/{id}", manualActivityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /manual_activity/{id}", manualActivityManager.HandleHTTPDelete)

	friendService := friends.NewService(postgres.DB)
	friendManager := friends.NewFriendHTTPHandler(friendService)

	mux.HandleFunc("POST /friend", friendManager.HandleHTTPPost)
	mux.HandleFunc("GET /friend", friendManager.HandleHTTPGet)
	mux.HandleFunc("GET /friend/{user_id}/{friend_id}", friendManager.HandleHTTPGetWithID)
	mux.HandleFunc("DELETE /friend/{user_id}/{friend_id}", friendManager.HandleHTTPDelete)

	activityParticipantService := activity_participants.NewService(postgres.DB)
	activityParticipantManager := activity_participants.NewActivityParticipantHTTPHandler(activityParticipantService)

	mux.HandleFunc("POST /activity_participant", activityParticipantManager.HandleHTTPPost)
	mux.HandleFunc("GET /activity_participants", activityParticipantManager.HandleHTTPGet)
	mux.HandleFunc("GET /activity_participant/{id}", activityParticipantManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /activity_participant/{id}", activityParticipantManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /activity_participant/{id}", activityParticipantManager.HandleHTTPDelete)

	locationService := locations.NewService(postgres.DB)
	locationManager := locations.NewLocationHTTPHandler(locationService)

	mux.HandleFunc("POST /location", locationManager.HandleHTTPPost)
	mux.HandleFunc("GET /locations", locationManager.HandleHTTPGet)
	mux.HandleFunc("GET /location/{id}", locationManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /location/{id}", locationManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /location/{id}", locationManager.HandleHTTPDelete)

	activityService := activities.NewService(postgres.DB)
	activityManager := activities.NewActivityHTTPHandler(activityService)

	mux.HandleFunc("POST /activity", activityManager.HandleHTTPPost)
	mux.HandleFunc("GET /activities", activityManager.HandleHTTPGet)
	mux.HandleFunc("GET /activity/{id}", activityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /activity/{id}", activityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /activity/{id}", activityManager.HandleHTTPDelete)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

	// Initialize the database connection
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userID := 1 // Example user ID
	err = scheduleCommonActivities(db, userID)
	if err != nil {
		log.Fatal(err)
	}
}

func scheduleCommonActivities(db *sql.DB, userID int) error {
	// Step 1: Identify common activities between the user and their friends
	query := `
		SELECT a.id, a.name
		FROM user_activity_preferences uap
		JOIN activities a ON uap.activity_id = a.id
		WHERE uap.user_id = $1
		AND uap.activity_id IN (
			SELECT uap2.activity_id
			FROM friends f
			JOIN user_activity_preferences uap2 ON f.friend_id = uap2.user_id
			WHERE f.user_id = $1
		)
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var activities []struct {
		ID   int
		Name string
	}
	for rows.Next() {
		var activity struct {
			ID   int
			Name string
		}
		if err := rows.Scan(&activity.ID, &activity.Name); err != nil {
			return err
		}
		activities = append(activities, activity)
	}

	// Step 2: Check for overlapping availability
	for _, activity := range activities {
		query := `
			SELECT ua.user_id, ua.day_of_week, ua.start_time, ua.end_time
			FROM user_availability ua
			JOIN friends f ON ua.user_id = f.friend_id
			WHERE f.user_id = $1
			AND ua.day_of_week IN (
				SELECT ua2.day_of_week
				FROM user_availability ua2
				WHERE ua2.user_id = $1
			)
			AND ua.start_time <= (
				SELECT ua2.end_time
				FROM user_availability ua2
				WHERE ua2.user_id = $1
				AND ua2.day_of_week = ua.day_of_week
			)
			AND ua.end_time >= (
				SELECT ua2.start_time
				FROM user_availability ua2
				WHERE ua2.user_id = $1
				AND ua2.day_of_week = ua.day_of_week
			)
		`
		rows, err := db.Query(query, userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		var availableFriends []int
		for rows.Next() {
			var friendID int
			var dayOfWeek string
			var startTime, endTime time.Time
			if err := rows.Scan(&friendID, &dayOfWeek, &startTime, &endTime); err != nil {
				return err
			}
			availableFriends = append(availableFriends, friendID)
		}

		// Step 3: Insert the identified activities into the user_activities table
		if len(availableFriends) > 0 {
			query := `
				INSERT INTO user_activities (user_id, activity_id, is_active, scheduled_at)
				VALUES ($1, $2, TRUE, NOW())
				RETURNING id
			`
			var userActivityID int
			err := db.QueryRow(query, userID, activity.ID).Scan(&userActivityID)
			if err != nil {
				return err
			}

			// Step 4: Link the participants by inserting entries into the activity_participants table
			for _, friendID := range availableFriends {
				query := `
					INSERT INTO activity_participants (user_id, activity_id, is_creator, is_active)
					VALUES ($1, $2, FALSE, TRUE)
				`
				_, err := db.Exec(query, friendID, activity.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
