package manual_activities

import (
	"encoding/json"
	"net/http"
)

// ManualActivityService defines the interface for CRUD operations on manual activities
type ManualActivityService interface {
	Create(manualActivity ManualActivity) (ManualActivity, error)
	ReadAll() ([]ManualActivity, error)
	Read(id string) (ManualActivity, bool, error)
	Update(id string, manualActivity ManualActivity) (ManualActivity, bool, error)
	Delete(id string) (bool, error)
}

// ManualActivityError represents an error response
type ManualActivityError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// ManualActivityHTTPHandler handles HTTP requests for manual activities
type ManualActivityHTTPHandler struct {
	manualActivityService ManualActivityService
}

// NewManualActivityHTTPHandler creates a new ManualActivityHTTPHandler
func NewManualActivityHTTPHandler(manualActivityService ManualActivityService) *ManualActivityHTTPHandler {
	return &ManualActivityHTTPHandler{
		manualActivityService: manualActivityService,
	}
}

// HandleHTTPPost handles the creation of a new manual activity
//	@Summary		Create a new manual activity
//	@Description	Create a new manual activity
//	@Tags			manual-activity
//	@Accept			json
//	@Produce		json
//	@Param			manualActivity	body		ManualActivity	true	"Manual Activity"
//	@Success		201				{object}	ManualActivity
//	@Failure		400				{object}	ManualActivityError
//	@Failure		500				{object}	ManualActivityError
//	@Router			/manual-activities [post]
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

// HandleHTTPGet handles the retrieval of all manual activities
//	@Summary		Get all manual activities
//	@Description	Retrieve all manual activities
//	@Tags			manual-activity
//	@Produce		json
//	@Success		200	{array}		ManualActivity
//	@Failure		400	{object}	ManualActivityError
//	@Failure		500	{object}	ManualActivityError
//	@Router			/manual-activities [get]
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

// HandleHTTPGetWithID handles the retrieval of a specific manual activity by ID
//	@Summary		Get a manual activity by ID
//	@Description	Retrieve a specific manual activity by its ID
//	@Tags			manual-activity
//	@Produce		json
//	@Param			id	path		string	true	"Manual Activity ID"
//	@Success		200	{object}	ManualActivity
//	@Failure		400	{object}	ManualActivityError
//	@Failure		404	{object}	ManualActivityError
//	@Failure		500	{object}	ManualActivityError
//	@Router			/manual-activities/{id} [get]
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

// HandleHTTPPut handles the updating of an existing manual activity
//	@Summary		Update a manual activity
//	@Description	Update an existing manual activity by its ID
//	@Tags			manual-activity
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string			true	"Manual Activity ID"
//	@Param			manualActivity	body		ManualActivity	true	"Manual Activity"
//	@Success		200				{object}	ManualActivity
//	@Failure		400				{object}	ManualActivityError
//	@Failure		404				{object}	ManualActivityError
//	@Failure		500				{object}	ManualActivityError
//	@Router			/manual-activities/{id} [put]
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

// HandleHTTPDelete handles the deletion of a manual activity
//	@Summary		Delete a manual activity
//	@Description	Delete a manual activity by its ID
//	@Tags			manual-activity
//	@Param			id	path	string	true	"Manual Activity ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ManualActivityError
//	@Failure		404	{object}	ManualActivityError
//	@Failure		500	{object}	ManualActivityError
//	@Router			/manual-activities/{id} [delete]
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
