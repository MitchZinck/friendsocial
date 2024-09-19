package scheduled_activities

import (
	"encoding/json"
	"friendsocial/user_activity_preferences"
	"net/http"
	"strconv"
)

// ScheduledActivityService defines the interface for scheduled activity operations.
type ScheduledActivityService interface {
	Create(scheduledActivity ScheduledActivity) (ScheduledActivity, error)
	CreateMultiple(activityID int, selectedDates []string, startTime string, endTime string, timeZone string) ([]ScheduledActivity, error)
	ReadAll() ([]ScheduledActivity, error)
	Read(id string) (ScheduledActivity, bool, error)
	Update(id string, scheduledActivity ScheduledActivity) (ScheduledActivity, bool, error)
	Delete(id string) (bool, error)
	CreateRepeatingScheduledActivity(preference user_activity_preferences.UserActivityPreference, startTime string, timeZone string) ([]ScheduledActivity, error)
	DeclineRepeatedActivity(userID int, scheduledActivityID int) error
}

// ScheduledActivityError represents an error response.
//
//	@swagger:model
type ScheduledActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// ScheduledActivityHTTPHandler is the HTTP handler for scheduled activity operations.
type ScheduledActivityHTTPHandler struct {
	scheduledActivityService ScheduledActivityService
	services                 *map[string]interface{}
}

// NewScheduledActivityHTTPHandler creates a new ScheduledActivityHTTPHandler.
func NewScheduledActivityHTTPHandler(scheduledActivityService ScheduledActivityService, services *map[string]interface{}) *ScheduledActivityHTTPHandler {
	return &ScheduledActivityHTTPHandler{
		scheduledActivityService: scheduledActivityService,
		services:                 services,
	}
}

// HandleHTTPPost handles the creation of a new user activity.
//
//	@Summary		Create a new scheduled activity
//	@Description	Create a new scheduled activity
//	@Tags			scheduled_activities
//	@Accept			json
//	@Produce		json
//	@Param			userActivity	body		UserActivity	true	"User Activity"
//	@Success		201				{object}	ScheduledActivity
//	@Failure		400				{object}	ScheduledActivityError
//	@Failure		500				{object}	ScheduledActivityError
//	@Router			/scheduled_activity [post]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var scheduledActivity ScheduledActivity
	err := json.NewDecoder(r.Body).Decode(&scheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newScheduledActivity, err := uH.scheduledActivityService.Create(scheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newScheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

type CreateMultipleRequest struct {
	ActivityID    int      `json:"activity_id"`
	SelectedDates []string `json:"selected_dates"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
	TimeZone      string   `json:"time_zone"`
}

// HandleHTTPPostMultiple handles the creation of multiple scheduled activities.
//
//	@Summary		Create multiple scheduled activities
//	@Description	Create multiple scheduled activities
//	@Tags			scheduled_activities
//	@Accept			json
//	@Produce		json
//	@Param			scheduledActivities	body		[]ScheduledActivity	true	"Scheduled Activities"
//	@Success		201				{object}	ScheduledActivity
//	@Failure		400				{object}	ScheduledActivityError
//	@Failure		500				{object}	ScheduledActivityError
//	@Router			/scheduled_activities [post]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPPostMultiple(w http.ResponseWriter, r *http.Request) {
	var createMultipleRequest CreateMultipleRequest
	err := json.NewDecoder(r.Body).Decode(&createMultipleRequest)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newScheduledActivities, err := uH.scheduledActivityService.CreateMultiple(
		createMultipleRequest.ActivityID,
		createMultipleRequest.SelectedDates,
		createMultipleRequest.StartTime,
		createMultipleRequest.EndTime,
		createMultipleRequest.TimeZone,
	)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newScheduledActivities)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet handles fetching all scheduled activities.
//
//	@Summary		Get all scheduled activities
//	@Description	Get all scheduled activities
//	@Tags			scheduled_activities
//	@Produce		json
//	@Success		200	{array}		ScheduledActivity
//	@Failure		400	{object}	ScheduledActivityError
//	@Failure		500	{object}	ScheduledActivityError
//	@Router			/scheduled_activity [get]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	scheduledActivities, err := uH.scheduledActivityService.ReadAll()
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(scheduledActivities)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID handles fetching a scheduled activity by ID.
//
//	@Summary		Get a scheduled activity by ID
//	@Description	Get a scheduled activity by ID
//	@Tags			scheduled_activities
//	@Produce		json
//	@Param			id	path		string	true	"Scheduled Activity ID"
//	@Success		200	{object}	ScheduledActivity
//	@Failure		400	{object}	ScheduledActivityError
//	@Failure		404	{object}	ScheduledActivityError
//	@Failure		500	{object}	ScheduledActivityError
//	@Router			/scheduled_activities/{id} [get]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	scheduledActivity, found, err := uH.scheduledActivityService.Read(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(scheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut handles updating a user activity by ID.
//
//	@Summary		Update a user activity by ID
//	@Description	Update a user activity by ID
//	@Tags			scheduled_activities
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string			true	"User Activity ID"
//	@Param			userActivity	body		UserActivity	true	"User Activity"
//	@Success		200				{object}	ScheduledActivity
//	@Failure		400				{object}	ScheduledActivityError
//	@Failure		404				{object}	ScheduledActivityError
//	@Failure		500				{object}	ScheduledActivityError
//	@Router			/scheduled_activities/{id} [put]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedScheduledActivity ScheduledActivity
	err := json.NewDecoder(r.Body).Decode(&updatedScheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	scheduledActivity, found, err := uH.scheduledActivityService.Update(id, updatedScheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(scheduledActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete handles deleting a user activity by ID.
//
//	@Summary		Delete a user activity by ID
//	@Description	Delete a user activity by ID
//	@Tags			scheduled_activities
//	@Param			id	path	string	true	"Scheduled Activity ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ScheduledActivityError
//	@Failure		404	{object}	ScheduledActivityError
//	@Failure		500	{object}	ScheduledActivityError
//	@Router			/scheduled_activities/{id} [delete]
func (uH *ScheduledActivityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := uH.scheduledActivityService.Delete(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uH *ScheduledActivityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ScheduledActivityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}

type RepeatScheduledActivityRequest struct {
	PreferenceID string `json:"preference_id"`
	StartTime    string `json:"start_time"`
	TimeZone     string `json:"time_zone"`
}

// HandleHTTPPostRepeatScheduledActivity handles the request to repeat a scheduled activity
//
//	@Summary		Repeat a scheduled activity
//	@Description	Repeat a scheduled activity
//	@Tags			scheduled_activities
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RepeatScheduledActivityRequest	true	"Repeat Scheduled Activity Request"
//	@Success		201		{array}		scheduled_activities.ScheduledActivity
//	@Failure		400		{object}	ScheduledActivityError
//	@Failure		500		{object}	ScheduledActivityError
//	@Router			/scheduled_activity/repeat [post]
func (h *ScheduledActivityHTTPHandler) HandleHTTPPostRepeatScheduledActivity(w http.ResponseWriter, r *http.Request) {
	var request RepeatScheduledActivityRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	preference, _, err := (*h.services)["user_activity_preferences"].(user_activity_preferences.UserActivityPreferenceService).Read(request.PreferenceID)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	scheduledActivities, err := h.scheduledActivityService.CreateRepeatingScheduledActivity(preference, request.StartTime, request.TimeZone)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(scheduledActivities)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

type DeclineRepeatedActivityRequest struct {
	UserID              string `json:"user_id"`
	ScheduledActivityID string `json:"scheduled_activity_id"`
}

// HandleHTTPDeclineRepeatedActivity handles the request to decline a repeated activity
//
//	@Summary		Decline a repeated activity
//	@Description	Decline a repeated activity
//	@Tags			preferences
//	@Accept			json
//	@Produce		json
//	@Param			request	body		DeclineRepeatedActivityRequest	true	"Decline Repeated Activity Request"
//	@Success		200		"No Content"
//	@Failure		400		{object}	ScheduledActivityError
//	@Failure		500		{object}	ScheduledActivityError
//	@Router			/scheduled_activity/repeat/decline [post]
func (h *ScheduledActivityHTTPHandler) HandleHTTPPostDeclineRepeatedActivity(w http.ResponseWriter, r *http.Request) {
	var request DeclineRepeatedActivityRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := strconv.Atoi(request.UserID)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	scheduledActivityID, err := strconv.Atoi(request.ScheduledActivityID)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.scheduledActivityService.DeclineRepeatedActivity(userID, scheduledActivityID)

	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
