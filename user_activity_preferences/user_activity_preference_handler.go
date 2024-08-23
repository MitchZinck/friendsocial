package user_activity_preferences

import (
	"encoding/json"
	"net/http"
)

type UserActivityPreferenceService interface {
	Create(preference UserActivityPreference) (UserActivityPreference, error)
	ReadAll() ([]UserActivityPreference, error)
	Read(id string) (UserActivityPreference, bool, error)
	Update(id string, preference UserActivityPreference) (UserActivityPreference, bool, error)
	Delete(id string) (bool, error)
}

type UserActivityPreferenceError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type UserActivityPreferenceHTTPHandler struct {
	preferenceService UserActivityPreferenceService
}

func NewUserActivityPreferenceHTTPHandler(preferenceService UserActivityPreferenceService) *UserActivityPreferenceHTTPHandler {
	return &UserActivityPreferenceHTTPHandler{
		preferenceService: preferenceService,
	}
}

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
