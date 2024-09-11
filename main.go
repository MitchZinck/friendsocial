package main

import (
	"net/http"

	"database/sql"
	"friendsocial/activities"
	"friendsocial/activity_participants"
	"friendsocial/friends"
	"friendsocial/locations"
	"friendsocial/postgres"
	"friendsocial/scheduled_activities"
	"friendsocial/user_activity_preferences"
	"friendsocial/user_availability"
	"friendsocial/users"
	"log"

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
	mux.HandleFunc("GET /user_availability/user/{user_id}", availabilityManager.HandleHTTPGetByUserID)
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

	scheduledActivityService := scheduled_activities.NewService(postgres.DB)
	scheduledActivityManager := scheduled_activities.NewScheduledActivityHTTPHandler(scheduledActivityService)

	mux.HandleFunc("POST /scheduled_activity", scheduledActivityManager.HandleHTTPPost)
	mux.HandleFunc("GET /scheduled_activities", scheduledActivityManager.HandleHTTPGet)
	mux.HandleFunc("GET /scheduled_activity/{id}", scheduledActivityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /scheduled_activity/{id}", scheduledActivityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /scheduled_activity/{id}", scheduledActivityManager.HandleHTTPDelete)

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
	mux.HandleFunc("GET /activity_participants/user/{user_id}", activityParticipantManager.HandleHTTPGetActivitiesByUserID)
	mux.HandleFunc("GET /activity_participants/scheduled_activity/{scheduled_activity_id}", activityParticipantManager.HandleHTTPGetParticipantsByActivityID)

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
}
