package activity_participants

import (
	"encoding/json"
	"net/http"
)

type ActivityParticipantService interface {
	Create(participant ActivityParticipant) (ActivityParticipant, error)
	ReadAll() ([]ActivityParticipant, error)
	Read(id string) (ActivityParticipant, bool, error)
	Update(id string, participant ActivityParticipant) (ActivityParticipant, bool, error)
	Delete(id string) (bool, error)
}

type ActivityParticipantError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type ActivityParticipantHTTPHandler struct {
	activityParticipantService ActivityParticipantService
}

func NewActivityParticipantHTTPHandler(activityParticipantService ActivityParticipantService) *ActivityParticipantHTTPHandler {
	return &ActivityParticipantHTTPHandler{
		activityParticipantService: activityParticipantService,
	}
}

func (aH *ActivityParticipantHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var participant ActivityParticipant
	err := json.NewDecoder(r.Body).Decode(&participant)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newParticipant, err := aH.activityParticipantService.Create(participant)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newParticipant)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (aH *ActivityParticipantHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	participants, err := aH.activityParticipantService.ReadAll()
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(participants)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (aH *ActivityParticipantHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	participant, found, err := aH.activityParticipantService.Read(id)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(participant)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (aH *ActivityParticipantHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedParticipant ActivityParticipant
	err := json.NewDecoder(r.Body).Decode(&updatedParticipant)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	participant, found, err := aH.activityParticipantService.Update(id, updatedParticipant)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(participant)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (aH *ActivityParticipantHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := aH.activityParticipantService.Delete(id)
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

func (aH *ActivityParticipantHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(ActivityParticipantError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
