package manual_activities

import (
	"encoding/json"
	"net/http"
)

type ManualActivityService interface {
	Create(manualActivity ManualActivity) (ManualActivity, error)
	ReadAll() ([]ManualActivity, error)
	Read(id string) (ManualActivity, bool, error)
	Update(id string, manualActivity ManualActivity) (ManualActivity, bool, error)
	Delete(id string) (bool, error)
}

type ManualActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type ManualActivityHTTPHandler struct {
	manualActivityService ManualActivityService
}

func NewManualActivityHTTPHandler(manualActivityService ManualActivityService) *ManualActivityHTTPHandler {
	return &ManualActivityHTTPHandler{
		manualActivityService: manualActivityService,
	}
}

func (mH *ManualActivityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var manualActivity ManualActivity
	err := json.NewDecoder(r.Body).Decode(&manualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newManualActivity, err := mH.manualActivityService.Create(manualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newManualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (mH *ManualActivityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	manualActivities, err := mH.manualActivityService.ReadAll()
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(manualActivities)
	if err != nil {
		mH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (mH *ManualActivityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	manualActivity, found, err := mH.manualActivityService.Read(id)
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		mH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(manualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (mH *ManualActivityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedManualActivity ManualActivity
	err := json.NewDecoder(r.Body).Decode(&updatedManualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	manualActivity, found, err := mH.manualActivityService.Update(id, updatedManualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		mH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(manualActivity)
	if err != nil {
		mH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (mH *ManualActivityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := mH.manualActivityService.Delete(id)
	if err != nil {
		mH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		mH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (mH *ManualActivityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ManualActivityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
