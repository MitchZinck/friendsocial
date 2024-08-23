package user_availability

import (
	"encoding/json"
	"net/http"
)

type UserAvailabilityService interface {
	Create(availability UserAvailability) (UserAvailability, error)
	ReadAll() ([]UserAvailability, error)
	Read(id string) (UserAvailability, bool, error)
	Update(id string, availability UserAvailability) (UserAvailability, bool, error)
	Delete(id string) (bool, error)
}

type UserAvailabilityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type UserAvailabilityHTTPHandler struct {
	availabilityService UserAvailabilityService
}

func NewUserAvailabilityHTTPHandler(availabilityService UserAvailabilityService) *UserAvailabilityHTTPHandler {
	return &UserAvailabilityHTTPHandler{
		availabilityService: availabilityService,
	}
}

func (uH *UserAvailabilityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var availability UserAvailability
	err := json.NewDecoder(r.Body).Decode(&availability)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newAvailability, err := uH.availabilityService.Create(availability)

	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newAvailability)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserAvailabilityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	availability, err := uH.availabilityService.ReadAll()
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(availability)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserAvailabilityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	availability, found, err := uH.availabilityService.Read(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(availability)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserAvailabilityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var newAvailability UserAvailability
	err := json.NewDecoder(r.Body).Decode(&newAvailability)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	availability, found, err := uH.availabilityService.Update(id, newAvailability)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(availability)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserAvailabilityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := uH.availabilityService.Delete(id)
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

func (uH *UserAvailabilityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserAvailabilityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
