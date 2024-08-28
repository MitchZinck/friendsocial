package user_activities

import (
	"encoding/json"
	"net/http"
)

// UserActivityService defines the interface for user activity operations.
type UserActivityService interface {
	Create(userActivity UserActivity) (UserActivity, error)
	ReadAll(userID string) ([]UserActivity, error)
	Read(id string) (UserActivity, bool, error)
	Update(id string, userActivity UserActivity) (UserActivity, bool, error)
	Delete(id string) (bool, error)
}

// UserActivityError represents an error response.
//
//	@swagger:model
type UserActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// UserActivityHTTPHandler is the HTTP handler for user activity operations.
type UserActivityHTTPHandler struct {
	userActivityService UserActivityService
}

// NewUserActivityHTTPHandler creates a new UserActivityHTTPHandler.
func NewUserActivityHTTPHandler(userActivityService UserActivityService) *UserActivityHTTPHandler {
	return &UserActivityHTTPHandler{
		userActivityService: userActivityService,
	}
}

// HandleHTTPPost handles the creation of a new user activity.
//
//	@Summary		Create a new user activity
//	@Description	Create a new user activity
//	@Tags			user_activities
//	@Accept			json
//	@Produce		json
//	@Param			userActivity	body		UserActivity	true	"User Activity"
//	@Success		201				{object}	UserActivity
//	@Failure		400				{object}	UserActivityError
//	@Failure		500				{object}	UserActivityError
//	@Router			/user_activities [post]
func (uH *UserActivityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var userActivity UserActivity
	err := json.NewDecoder(r.Body).Decode(&userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newUserActivity, err := uH.userActivityService.Create(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet handles fetching all user activities for a specific user.
//
//	@Summary		Get all user activities
//	@Description	Get all user activities
//	@Tags			user_activities
//	@Produce		json
//	@Param			userID	path	string	true	"User ID"
//	@Success		200	{array}		UserActivity
//	@Failure		400	{object}	UserActivityError
//	@Failure		500	{object}	UserActivityError
//	@Router			/user_activities/user/{userID} [get]
func (uH *UserActivityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	userActivities, err := uH.userActivityService.ReadAll(userID)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivities)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID handles fetching a user activity by ID.
//
//	@Summary		Get a user activity by ID
//	@Description	Get a user activity by ID
//	@Tags			user_activities
//	@Produce		json
//	@Param			id	path		string	true	"User Activity ID"
//	@Success		200	{object}	UserActivity
//	@Failure		400	{object}	UserActivityError
//	@Failure		404	{object}	UserActivityError
//	@Failure		500	{object}	UserActivityError
//	@Router			/user_activities/{id} [get]
func (uH *UserActivityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userActivity, found, err := uH.userActivityService.Read(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut handles updating a user activity by ID.
//
//	@Summary		Update a user activity by ID
//	@Description	Update a user activity by ID
//	@Tags			user_activities
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string			true	"User Activity ID"
//	@Param			userActivity	body		UserActivity	true	"User Activity"
//	@Success		200				{object}	UserActivity
//	@Failure		400				{object}	UserActivityError
//	@Failure		404				{object}	UserActivityError
//	@Failure		500				{object}	UserActivityError
//	@Router			/user_activities/{id} [put]
func (uH *UserActivityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedUserActivity UserActivity
	err := json.NewDecoder(r.Body).Decode(&updatedUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userActivity, found, err := uH.userActivityService.Update(id, updatedUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete handles deleting a user activity by ID.
//
//	@Summary		Delete a user activity by ID
//	@Description	Delete a user activity by ID
//	@Tags			user_activities
//	@Param			id	path	string	true	"User Activity ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	UserActivityError
//	@Failure		404	{object}	UserActivityError
//	@Failure		500	{object}	UserActivityError
//	@Router			/user_activities/{id} [delete]
func (uH *UserActivityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := uH.userActivityService.Delete(id)
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

func (uH *UserActivityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserActivityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
