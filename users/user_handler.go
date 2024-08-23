package users

import (
	"encoding/json"
	"net/http"
)

type UserService interface {
	Create(user User) (User, error)
	ReadAll() ([]User, error)
	Read(id string) (User, bool, error)
	Update(id string, user User) (User, bool, error)
	Delete(id string) (bool, error)
}

type UserError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type UserHTTPHandler struct {
	userService UserService
}

func NewUserHTTPHandler(userService UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

func (uH *UserHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, err := uH.userService.Create(user)

	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newUser)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	users, err := uH.userService.ReadAll()
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, found, err := uH.userService.Read(id)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, found, err := uH.userService.Update(id, newUser)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		uH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		uH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uH *UserHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := uH.userService.Delete(id)
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

func (uH *UserHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(UserError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
