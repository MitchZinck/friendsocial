package activity_locations

import (
	"encoding/json"
	"net/http"
)

// ActivityLocationService defines the service interface for handling Activity Locations
type ActivityLocationService interface {
	Create(activityLocation ActivityLocation) (ActivityLocation, error)
	ReadAll() ([]ActivityLocation, error)
	Read(id string) (ActivityLocation, bool, error)
	Update(id string, activityLocation ActivityLocation) (ActivityLocation, bool, error)
	Delete(id string) (bool, error)
}

// ActivityLocationError represents the error response structure
type ActivityLocationError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// ActivityLocationHTTPHandler handles HTTP requests for Activity Locations
type ActivityLocationHTTPHandler struct {
	activityLocationService ActivityLocationService
}

// NewActivityLocationHTTPHandler initializes a new ActivityLocationHTTPHandler
func NewActivityLocationHTTPHandler(activityLocationService ActivityLocationService) *ActivityLocationHTTPHandler {
	return &ActivityLocationHTTPHandler{
		activityLocationService: activityLocationService,
	}
}

// HandleHTTPPost handles the creation of a new Activity Location
//	@Summary		Create a new Activity Location
//	@Description	Create a new Activity Location
//	@Tags			activity_locations
//	@Accept			json
//	@Produce		json
//	@Param			activityLocation	body		ActivityLocation	true	"Activity Location data"
//	@Success		201					{object}	ActivityLocation
//	@Failure		400					{object}	ActivityLocationError
//	@Failure		500					{object}	ActivityLocationError
//	@Router			/activity-locations [post]
func (aH *ActivityLocationHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var activityLocation ActivityLocation
	err := json.NewDecoder(r.Body).Decode(&activityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newActivityLocation, err := aH.activityLocationService.Create(activityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newActivityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet handles retrieving all Activity Locations
//	@Summary		Get all Activity Locations
//	@Description	Get all Activity Locations
//	@Tags			activity_locations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		ActivityLocation
//	@Failure		400	{object}	ActivityLocationError
//	@Failure		500	{object}	ActivityLocationError
//	@Router			/activity-locations [get]
func (aH *ActivityLocationHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	activityLocations, err := aH.activityLocationService.ReadAll()
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activityLocations)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID handles retrieving a single Activity Location by ID
//	@Summary		Get an Activity Location by ID
//	@Description	Get a specific Activity Location by ID
//	@Tags			activity_locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Activity Location ID"
//	@Success		200	{object}	ActivityLocation
//	@Failure		400	{object}	ActivityLocationError
//	@Failure		404	{object}	ActivityLocationError
//	@Failure		500	{object}	ActivityLocationError
//	@Router			/activity-locations/{id} [get]
func (aH *ActivityLocationHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	activityLocation, found, err := aH.activityLocationService.Read(id)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut handles updating an existing Activity Location by ID
//	@Summary		Update an existing Activity Location
//	@Description	Update an existing Activity Location by ID
//	@Tags			activity_locations
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string				true	"Activity Location ID"
//	@Param			activityLocation	body		ActivityLocation	true	"Updated Activity Location data"
//	@Success		200					{object}	ActivityLocation
//	@Failure		400					{object}	ActivityLocationError
//	@Failure		404					{object}	ActivityLocationError
//	@Failure		500					{object}	ActivityLocationError
//	@Router			/activity-locations/{id} [put]
func (aH *ActivityLocationHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedActivityLocation ActivityLocation
	err := json.NewDecoder(r.Body).Decode(&updatedActivityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	activityLocation, found, err := aH.activityLocationService.Update(id, updatedActivityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(activityLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete handles deleting an Activity Location by ID
//	@Summary		Delete an Activity Location
//	@Description	Delete an Activity Location by ID
//	@Tags			activity_locations
//	@Param			id	path	string	true	"Activity Location ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ActivityLocationError
//	@Failure		404	{object}	ActivityLocationError
//	@Failure		500	{object}	ActivityLocationError
//	@Router			/activity-locations/{id} [delete]
func (aH *ActivityLocationHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := aH.activityLocationService.Delete(id)
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

// errorResponse sends an error response with the specified status code and message
func (aH *ActivityLocationHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ActivityLocationError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
