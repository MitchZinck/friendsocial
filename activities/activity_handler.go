package activities

import (
	"encoding/json"
	"net/http"
)

// ActivityService defines the interface for activity services
type ActivityService interface {
	Create(activity Activity) (Activity, error)
	ReadAll() ([]Activity, error)
	Read(id string) (Activity, bool, error)
	Update(id string, activity Activity) (Activity, bool, error)
	Delete(id string) (bool, error)
}

// ActivityError represents an error response
type ActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// ActivityHTTPHandler handles HTTP requests for activities
type ActivityHTTPHandler struct {
	activityService ActivityService
}

// NewActivityHTTPHandler creates a new ActivityHTTPHandler
func NewActivityHTTPHandler(activityService ActivityService) *ActivityHTTPHandler {
	return &ActivityHTTPHandler{
		activityService: activityService,
	}
}

// HandleHTTPPost handles the creation of a new activity
//	@Summary		Create a new activity
//	@Description	Create a new activity
//	@Tags			activities
//	@Accept			json
//	@Produce		json
//	@Param			activity	body		Activity	true	"Activity object"
//	@Success		201			{object}	Activity
//	@Failure		400			{object}	ActivityError
//	@Failure		500			{object}	ActivityError
//	@Router			/activities [post]
func (aH *ActivityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var activity Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newActivity, err := aH.activityService.Create(activity)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newActivity)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet handles fetching all activities
//	@Summary		Get all activities
//	@Description	Get all activities
//	@Tags			activities
//	@Produce		json
//	@Success		200	{array}		Activity
//	@Failure		400	{object}	ActivityError
//	@Failure		500	{object}	ActivityError
//	@Router			/activities [get]
func (aH *ActivityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	activities, err := aH.activityService.ReadAll()
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activities)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID handles fetching an activity by ID
//	@Summary		Get an activity by ID
//	@Description	Get an activity by ID
//	@Tags			activities
//	@Produce		json
//	@Param			id	path		string	true	"Activity ID"
//	@Success		200	{object}	Activity
//	@Failure		400	{object}	ActivityError
//	@Failure		404	{object}	ActivityError
//	@Failure		500	{object}	ActivityError
//	@Router			/activities/{id} [get]
func (aH *ActivityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	activity, found, err := aH.activityService.Read(id)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activity)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut handles updating an activity
//	@Summary		Update an activity by ID
//	@Description	Update an activity by ID
//	@Tags			activities
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string		true	"Activity ID"
//	@Param			activity	body		Activity	true	"Updated Activity object"
//	@Success		200			{object}	Activity
//	@Failure		400			{object}	ActivityError
//	@Failure		404			{object}	ActivityError
//	@Failure		500			{object}	ActivityError
//	@Router			/activities/{id} [put]
func (aH *ActivityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedActivity Activity
	err := json.NewDecoder(r.Body).Decode(&updatedActivity)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	activity, found, err := aH.activityService.Update(id, updatedActivity)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activity)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete handles deleting an activity by ID
//	@Summary		Delete an activity by ID
//	@Description	Delete an activity by ID
//	@Tags			activities
//	@Param			id	path	string	true	"Activity ID"
//	@Success		204
//	@Failure		400	{object}	ActivityError
//	@Failure		404	{object}	ActivityError
//	@Failure		500	{object}	ActivityError
//	@Router			/activities/{id} [delete]
func (aH *ActivityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := aH.activityService.Delete(id)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (aH *ActivityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ActivityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
