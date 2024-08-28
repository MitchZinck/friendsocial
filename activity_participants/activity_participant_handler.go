package activity_participants

import (
	"encoding/json"
	"net/http"
)

// ActivityParticipantService defines the methods for handling activity participants
type ActivityParticipantService interface {
	Create(participant ActivityParticipant) (ActivityParticipant, error)
	ReadAll() ([]ActivityParticipant, error)
	Read(id string) (ActivityParticipant, bool, error)
	Update(id string, participant ActivityParticipant) (ActivityParticipant, bool, error)
	Delete(id string) (bool, error)
}

// ActivityParticipantError represents the structure of an error response
// swagger:model
type ActivityParticipantError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// ActivityParticipantHTTPHandler handles HTTP requests for activity participants
type ActivityParticipantHTTPHandler struct {
	activityParticipantService ActivityParticipantService
}

// NewActivityParticipantHTTPHandler creates a new handler for activity participants
func NewActivityParticipantHTTPHandler(activityParticipantService ActivityParticipantService) *ActivityParticipantHTTPHandler {
	return &ActivityParticipantHTTPHandler{
		activityParticipantService: activityParticipantService,
	}
}

// HandleHTTPPost creates a new activity participant
//	@Summary		Create a new activity participant
//	@Description	Create a new activity participant
//	@Tags			participants
//	@Accept			json
//	@Produce		json
//	@Param			participant	body		ActivityParticipant	true	"Activity Participant"
//	@Success		201			{object}	ActivityParticipant
//	@Failure		400			{object}	ActivityParticipantError
//	@Failure		500			{object}	ActivityParticipantError
//	@Router			/participants [post]
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

// HandleHTTPGet retrieves all activity participants
//	@Summary		Get all activity participants
//	@Description	Retrieve all activity participants
//	@Tags			participants
//	@Produce		json
//	@Success		200	{array}		ActivityParticipant
//	@Failure		400	{object}	ActivityParticipantError
//	@Failure		500	{object}	ActivityParticipantError
//	@Router			/participants [get]
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

// HandleHTTPGetWithID retrieves a specific activity participant by ID
//	@Summary		Get an activity participant by ID
//	@Description	Retrieve an activity participant by ID
//	@Tags			participants
//	@Produce		json
//	@Param			id	path		string	true	"Activity Participant ID"
//	@Success		200	{object}	ActivityParticipant
//	@Failure		400	{object}	ActivityParticipantError
//	@Failure		404	{object}	ActivityParticipantError
//	@Failure		500	{object}	ActivityParticipantError
//	@Router			/participants/{id} [get]
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

// HandleHTTPPut updates an existing activity participant
//	@Summary		Update an activity participant
//	@Description	Update an existing activity participant
//	@Tags			participants
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string				true	"Activity Participant ID"
//	@Param			participant	body		ActivityParticipant	true	"Updated Activity Participant"
//	@Success		200			{object}	ActivityParticipant
//	@Failure		400			{object}	ActivityParticipantError
//	@Failure		404			{object}	ActivityParticipantError
//	@Failure		500			{object}	ActivityParticipantError
//	@Router			/participants/{id} [put]
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

// HandleHTTPDelete deletes an activity participant
//	@Summary		Delete an activity participant
//	@Description	Delete an activity participant by ID
//	@Tags			participants
//	@Param			id	path	string	true	"Activity Participant ID"
//	@Success		204
//	@Failure		400	{object}	ActivityParticipantError
//	@Failure		404	{object}	ActivityParticipantError
//	@Router			/participants/{id} [delete]
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
