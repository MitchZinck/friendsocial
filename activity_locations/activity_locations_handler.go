package activity_locations

import (
	"encoding/json"
	"net/http"
)

type ActivityLocationService interface {
	Create(activityLocation ActivityLocation) (ActivityLocation, error)
	ReadAll() ([]ActivityLocation, error)
	Read(id string) (ActivityLocation, bool, error)
	Update(id string, activityLocation ActivityLocation) (ActivityLocation, bool, error)
	Delete(id string) (bool, error)
}

type ActivityLocationError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type ActivityLocationHTTPHandler struct {
	activityLocationService ActivityLocationService
}

func NewActivityLocationHTTPHandler(activityLocationService ActivityLocationService) *ActivityLocationHTTPHandler {
	return &ActivityLocationHTTPHandler{
		activityLocationService: activityLocationService,
	}
}

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
