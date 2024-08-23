package user_activities

import (
	"encoding/json"
	"net/http"
)

type UserActivityService interface {
	Create(userActivity UserActivity) (UserActivity, error)
	ReadAll() ([]UserActivity, error)
	Read(id string) (UserActivity, bool, error)
	Update(id string, userActivity UserActivity) (UserActivity, bool, error)
	Delete(id string) (bool, error)
}

type UserActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type UserActivityHTTPHandler struct {
	userActivityService UserActivityService
}

func NewUserActivityHTTPHandler(userActivityService UserActivityService) *UserActivityHTTPHandler {
	return &UserActivityHTTPHandler{
		userActivityService: userActivityService,
	}
}

func (uH *UserActivityHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var userActivity UserActivity
	err := json.NewDecoder(r.Body).Decode(&userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newUserActivity, err := uH.userActivityService.Create(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserActivityHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	userActivities, err := uH.userActivityService.ReadAll()
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivities)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserActivityHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	userActivity, found, err := uH.userActivityService.Read(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserActivityHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedUserActivity UserActivity
	err := json.NewDecoder(r.Body).Decode(&updatedUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userActivity, found, err := uH.userActivityService.Update(id, updatedUserActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userActivity)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserActivityHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := uH.userActivityService.Delete(id)
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

func (uH *UserActivityHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserActivityError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
