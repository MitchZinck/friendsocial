package user_availability

import (
	"encoding/json"
	"net/http"
)

type UserAvailabilityService interface {
	Create(availability UserAvailability) (UserAvailability, error)
	ReadAll() ([]UserAvailability, error)
	ReadByUserID(userID string) ([]UserAvailability, error)
	Read(id string) (UserAvailability, bool, error)
	Update(id string, availability UserAvailability) (UserAvailability, bool, error)
	Delete(id string) (bool, error)
}

// UserAvailabilityError represents an error response
type UserAvailabilityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// UserAvailabilityHTTPHandler handles HTTP requests for user availability
type UserAvailabilityHTTPHandler struct {
	availabilityService UserAvailabilityService
}

// NewUserAvailabilityHTTPHandler creates a new UserAvailabilityHTTPHandler
func NewUserAvailabilityHTTPHandler(availabilityService UserAvailabilityService) *UserAvailabilityHTTPHandler {
	return &UserAvailabilityHTTPHandler{
		availabilityService: availabilityService,
	}
}

// HandleHTTPPost handles the creation of new user availability
//
//	@Summary		Create new user availability
//	@Description	Create a new availability record for a user
//	@Tags			User Availability
//	@Accept			json
//	@Produce		json
//	@Param			availability	body		UserAvailability	true	"User Availability"
//	@Success		201				{object}	UserAvailability
//	@Failure		400				{object}	UserAvailabilityError
//	@Failure		500				{object}	UserAvailabilityError
//	@Router			/user_availability [post]
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

// HandleHTTPGet handles retrieving all user availability records
//
//	@Summary		Get all user availability records
//	@Description	Retrieve all availability records for all users
//	@Tags			User Availability
//	@Produce		json
//	@Success		200	{array}		UserAvailability
//	@Failure		400	{object}	UserAvailabilityError
//	@Failure		500	{object}	UserAvailabilityError
//	@Router			/user_availability [get]
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

// HandleHTTPGetWithID handles retrieving a user availability record by ID
//
//	@Summary		Get a user availability record by ID
//	@Description	Retrieve a specific availability record by its ID
//	@Tags			User Availability
//	@Produce		json
//	@Param			id	path		string	true	"User Availability ID"
//	@Success		200	{object}	UserAvailability
//	@Failure		400	{object}	UserAvailabilityError
//	@Failure		404	{object}	UserAvailabilityError
//	@Failure		500	{object}	UserAvailabilityError
//	@Router			/user_availability/{id} [get]
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

// HandleHTTPPut handles updating a user availability record by ID
//
//	@Summary		Update a user availability record by ID
//	@Description	Update a specific availability record by its ID
//	@Tags			User Availability
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string				true	"User Availability ID"
//	@Param			availability	body		UserAvailability	true	"Updated User Availability"
//	@Success		200				{object}	UserAvailability
//	@Failure		400				{object}	UserAvailabilityError
//	@Failure		404				{object}	UserAvailabilityError
//	@Failure		500				{object}	UserAvailabilityError
//	@Router			/user_availability/{id} [put]
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

// HandleHTTPDelete handles deleting a user availability record by ID
//
//	@Summary		Delete a user availability record by ID
//	@Description	Delete a specific availability record by its ID
//	@Tags			User Availability
//	@Param			id	path	string	true	"User Availability ID"
//	@Success		204
//	@Failure		400	{object}	UserAvailabilityError
//	@Failure		404	{object}	UserAvailabilityError
//	@Failure		500	{object}	UserAvailabilityError
//	@Router			/user_availability/{id} [delete]
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

// HandleHTTPGetByUserID handles retrieving a user availability record by user ID
//
//	@Summary		Get user availability by user ID
//	@Description	Retrieve availability records for a specific user by their ID
//	@Tags			User Availability
//	@Produce		json
//	@Param			user_id	path	string	true	"User ID"
//	@Success		200	{array}		UserAvailability
//	@Failure		400	{object}	UserAvailabilityError
//	@Failure		404	{object}	UserAvailabilityError
//	@Failure		500	{object}	UserAvailabilityError
//	@Router			/user_availability/user/{user_id} [get]
func (uH *UserAvailabilityHTTPHandler) HandleHTTPGetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	availability, err := uH.availabilityService.ReadByUserID(userID)
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

// errorResponse sends a JSON error response
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
