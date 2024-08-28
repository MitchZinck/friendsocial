package users

import (
	"encoding/json"
	"net/http"
)

// UserService defines the interface for user-related operations
type UserService interface {
	Create(user User) (User, error)
	ReadAll() ([]User, error)
	Read(id string) (User, bool, error)
	Update(id string, user User) (User, bool, error)
	Delete(id string) (bool, error)
	PartialUpdate(id string, updates map[string]interface{}) (User, bool, error) // Add this line
}

// UserError defines the structure for an error response
type UserError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// UserHTTPHandler handles HTTP requests related to users
type UserHTTPHandler struct {
	userService UserService
}

// NewUserHTTPHandler creates a new UserHTTPHandler
func NewUserHTTPHandler(userService UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

// HandleHTTPPost creates a new user
//
//	@Summary		Create a new user
//	@Description	Create a new user with the input payload
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User	true	"User to be created"
//	@Success		201		{object}	User
//	@Failure		400		{object}	UserError
//	@Failure		500		{object}	UserError
//	@Router			/users [post]
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

// HandleHTTPGet retrieves all users
//
//	@Summary		Get all users
//	@Description	Retrieve all users
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		User
//	@Failure		400	{object}	UserError
//	@Failure		500	{object}	UserError
//	@Router			/users [get]
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

// HandleHTTPGetWithID retrieves a user by ID
//
//	@Summary		Get a user by ID
//	@Description	Retrieve a user by their ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	User
//	@Failure		400	{object}	UserError
//	@Failure		404	{object}	UserError
//	@Failure		500	{object}	UserError
//	@Router			/users/{id} [get]
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

// HandleHTTPPut updates a user by ID
//
//	@Summary		Update a user by ID
//	@Description	Update an existing user with the provided payload
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"User ID"
//	@Param			user	body		User	true	"Updated user data"
//	@Success		200		{object}	User
//	@Failure		400		{object}	UserError
//	@Failure		404		{object}	UserError
//	@Failure		500		{object}	UserError
//	@Router			/users/{id} [put]
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

// HandleHTTPDelete deletes a user by ID
//
//	@Summary		Delete a user by ID
//	@Description	Delete an existing user by their ID
//	@Tags			users
//	@Param			id	path	string	true	"User ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	UserError
//	@Failure		404	{object}	UserError
//	@Router			/users/{id} [delete]
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

// HandleHTTPPatch updates a user partially by ID
//
//	@Summary		Partially update a user by ID
//	@Description	Update specific fields of an existing user with the provided payload
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"User ID"
//	@Param			updates	body		object	true	"Partial user data to update"
//	@Success		200		{object}	User
//	@Failure		400		{object}	UserError
//	@Failure		404		{object}	UserError
//	@Failure		500		{object}	UserError
//	@Router			/users/{id} [patch]
func (uH *UserHTTPHandler) HandleHTTPPatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updates map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		uH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, found, err := uH.userService.PartialUpdate(id, updates)
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
