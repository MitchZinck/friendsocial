package activities

import (
	"encoding/json"
	"net/http"
)

type ActivityService interface {
	Create(activity Activity) (Activity, error)
	ReadAll() ([]Activity, error)
	Read(id string) (Activity, bool, error)
	Update(id string, activity Activity) (Activity, bool, error)
	Delete(id string) (bool, error)
}

type ActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type ActivityHTTPHandler struct {
	activityService ActivityService
}

func NewActivityHTTPHandler(activityService ActivityService) *ActivityHTTPHandler {
	return &ActivityHTTPHandler{
		activityService: activityService,
	}
}

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
