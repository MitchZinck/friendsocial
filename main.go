package main

import (
	"net/http"

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

	mux.HandleFunc("POST /user_activity_preferences", userActivityPreferenceManager.HandleHTTPPost)
	mux.HandleFunc("GET /user_activity_preferences", userActivityPreferenceManager.HandleHTTPGet)
	mux.HandleFunc("GET /user_activity_preferences/{id}", userActivityPreferenceManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /user_activity_preferences/{id}", userActivityPreferenceManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /user_activity_preferences/{id}", userActivityPreferenceManager.HandleHTTPDelete)

	userActivityService := user_activities.NewService(postgres.DB)
	userActivityManager := user_activities.NewUserActivityHTTPHandler(userActivityService)

	mux.HandleFunc("POST /user_activity", userActivityManager.HandleHTTPPost)
	mux.HandleFunc("GET /user_activity", userActivityManager.HandleHTTPGet)
	mux.HandleFunc("GET /user_activity/{id}", userActivityManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /user_activity/{id}", userActivityManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /user_activity/{id}", userActivityManager.HandleHTTPDelete)

	manualActivityService := manual_activities.NewService(postgres.DB)
	manualActivityManager := manual_activities.NewManualActivityHTTPHandler(manualActivityService)

	mux.HandleFunc("POST /manual_activity", manualActivityManager.HandleHTTPPost)
	mux.HandleFunc("GET /manual_activity", manualActivityManager.HandleHTTPGet)
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
	mux.HandleFunc("GET /activity_participant", activityParticipantManager.HandleHTTPGet)
	mux.HandleFunc("GET /activity_participant/{id}", activityParticipantManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /activity_participant/{id}", activityParticipantManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /activity_participant/{id}", activityParticipantManager.HandleHTTPDelete)

	activityLocationService := activity_locations.NewService(postgres.DB)
	activityLocationManager := activity_locations.NewActivityLocationHTTPHandler(activityLocationService)

	mux.HandleFunc("POST /activity_location", activityLocationManager.HandleHTTPPost)
	mux.HandleFunc("GET /activity_location", activityLocationManager.HandleHTTPGet)
	mux.HandleFunc("GET /activity_location/{id}", activityLocationManager.HandleHTTPGetWithID)
	mux.HandleFunc("PUT /activity_location/{id}", activityLocationManager.HandleHTTPPut)
	mux.HandleFunc("DELETE /activity_location/{id}", activityLocationManager.HandleHTTPDelete)

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
}
