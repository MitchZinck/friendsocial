package user_activity_preferences_participants

import (
	"encoding/json"
	"net/http"
)

type UserActivityPreferenceParticipantService interface {
	Create(participant UserActivityPreferenceParticipant) (UserActivityPreferenceParticipant, error)
	ReadAll() ([]UserActivityPreferenceParticipant, error)
	Read(id string) (UserActivityPreferenceParticipant, bool, error)
	Update(id string, participant UserActivityPreferenceParticipant) (UserActivityPreferenceParticipant, bool, error)
	Delete(id string) (bool, error)
	ReadByPreferenceID(preferenceID string) ([]UserActivityPreferenceParticipant, error)
}

// UserActivityPreferenceParticipantError represents the error response
type UserActivityPreferenceParticipantError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type UserActivityPreferenceParticipantHTTPHandler struct {
	participantService UserActivityPreferenceParticipantService
}

func NewUserActivityPreferenceParticipantHTTPHandler(participantService UserActivityPreferenceParticipantService) *UserActivityPreferenceParticipantHTTPHandler {
	return &UserActivityPreferenceParticipantHTTPHandler{
		participantService: participantService,
	}
}

// Implement HandleHTTPPost, HandleHTTPGet, HandleHTTPGetWithID, HandleHTTPPut, HandleHTTPDelete methods similar to UserActivityPreferenceHTTPHandler
func (h *UserActivityPreferenceParticipantHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	participant := UserActivityPreferenceParticipant{}
	err := json.NewDecoder(r.Body).Decode(&participant)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	createdParticipant, err := h.participantService.Create(participant)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(createdParticipant)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserActivityPreferenceParticipantHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	participants, err := h.participantService.ReadAll()
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(participants)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserActivityPreferenceParticipantHTTPHandler) HandleHTTPGetByPreferenceID(w http.ResponseWriter, r *http.Request) {
	preferenceID := r.PathValue("preference_id")

	participants, err := h.participantService.ReadByPreferenceID(preferenceID)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(participants)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserActivityPreferenceParticipantHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	success, err := h.participantService.Delete(id)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !success {
		h.errorResponse(w, http.StatusNotFound, "Participant not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserActivityPreferenceParticipantHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	participant := UserActivityPreferenceParticipant{}
	err := json.NewDecoder(r.Body).Decode(&participant)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedParticipant, success, err := h.participantService.Update(id, participant)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !success {
		h.errorResponse(w, http.StatusNotFound, "Participant not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedParticipant)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// errorResponse sends an error response with the given status code and error message
func (h *UserActivityPreferenceParticipantHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserActivityPreferenceParticipantError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
