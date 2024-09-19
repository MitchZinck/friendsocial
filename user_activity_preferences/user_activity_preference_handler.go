package user_activity_preferences

import (
	"encoding/json"
	"net/http"
)

// UserActivityPreferenceService defines the interface for user activity preference operations
type UserActivityPreferenceService interface {
	Create(preference UserActivityPreference) (UserActivityPreference, error)
	ReadAll() ([]UserActivityPreference, error)
	Read(id string) (UserActivityPreference, bool, error)
	Update(id string, preference UserActivityPreference) (UserActivityPreference, bool, error)
	Delete(id string) (bool, error)
	ReadByUserID(userID string) ([]UserActivityPreference, error)
}

// UserActivityPreferenceError represents the error response
type UserActivityPreferenceError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// UserActivityPreferenceHTTPHandler handles HTTP requests for user activity preferences
type UserActivityPreferenceHTTPHandler struct {
	preferenceService UserActivityPreferenceService
}

// NewUserActivityPreferenceHTTPHandler creates a new HTTP handler for user activity preferences
func NewUserActivityPreferenceHTTPHandler(preferenceService UserActivityPreferenceService) *UserActivityPreferenceHTTPHandler {
	return &UserActivityPreferenceHTTPHandler{
		preferenceService: preferenceService,
	}
}

// HandleHTTPPost creates a new user activity preference
//
//	@Summary		Create a new user activity preference
//	@Description	Create a new user activity preference
//	@Tags			preferences
//	@Accept			json
//	@Produce		json
//	@Param			preference	body		UserActivityPreference	true	"User Activity Preference"
//	@Success		201			{object}	UserActivityPreference
//	@Failure		400			{object}	UserActivityPreferenceError
//	@Failure		500			{object}	UserActivityPreferenceError
//	@Router			/preferences [post]
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var preference UserActivityPreference
	err := json.NewDecoder(r.Body).Decode(&preference)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newPreference, err := h.preferenceService.Create(preference)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newPreference)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet retrieves all user activity preferences
//
//	@Summary		Get all user activity preferences
//	@Description	Retrieve all user activity preferences
//	@Tags			preferences
//	@Produce		json
//	@Success		200	{array}		UserActivityPreference
//	@Failure		400	{object}	UserActivityPreferenceError
//	@Failure		500	{object}	UserActivityPreferenceError
//	@Router			/preferences [get]
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	preferences, err := h.preferenceService.ReadAll()
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(preferences)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID retrieves a user activity preference by ID
// HandleHTTPGetWithID retrieves a user activity preference by ID
//
//	@Summary		Get a user activity preference by ID
//	@Description	Retrieve a user activity preference by ID
//	@Tags			preferences
//	@Produce		json
//	@Param			id	path		string	true	"User Activity Preference ID"
//	@Success		200	{object}	UserActivityPreference
//	@Failure		400	{object}	UserActivityPreferenceError
//	@Failure		404	{object}	UserActivityPreferenceError
//	@Failure		500	{object}	UserActivityPreferenceError
//	@Router			/preferences/{id} [get]
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	preference, found, err := h.preferenceService.Read(id)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		h.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(preference)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut updates a user activity preference by ID
//
//	@Summary		Update a user activity preference by ID
//	@Description	Update a user activity preference by ID
//	@Tags			preferences
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"User Activity Preference ID"
//	@Param			preference	body		UserActivityPreference	true	"Updated User Activity Preference"
//	@Success		200			{object}	UserActivityPreference
//	@Failure		400			{object}	UserActivityPreferenceError
//	@Failure		404			{object}	UserActivityPreferenceError
//	@Failure		500			{object}	UserActivityPreferenceError
//	@Router			/preferences/{id} [put]
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var newPreference UserActivityPreference
	err := json.NewDecoder(r.Body).Decode(&newPreference)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	preference, found, err := h.preferenceService.Update(id, newPreference)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		h.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(preference)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete deletes a user activity preference by ID
//
//	@Summary		Delete a user activity preference by ID
//	@Description	Delete a user activity preference by ID
//	@Tags			preferences
//	@Param			id	path	string	true	"User Activity Preference ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	UserActivityPreferenceError
//	@Failure		404	{object}	UserActivityPreferenceError
//	@Failure		500	{object}	UserActivityPreferenceError
//	@Router			/preferences/{id} [delete]
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := h.preferenceService.Delete(id)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		h.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// errorResponse sends an error response with the given status code and error message
func (h *UserActivityPreferenceHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserActivityPreferenceError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}

// Add a new method to handle getting preferences by user ID
func (h *UserActivityPreferenceHTTPHandler) HandleHTTPGetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	preferences, err := h.preferenceService.ReadByUserID(userID)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(preferences)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
