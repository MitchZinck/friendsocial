package locations

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// LocationService defines the service interface for handling Locations
type LocationService interface {
	Create(location Location) (Location, error)
	ReadAll() ([]Location, error)
	Read(ids []int) ([]Location, error)
	Update(id string, location Location) (Location, bool, error)
	Delete(id string) (bool, error)
}

// LocationError represents the error response structure
type LocationError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// LocationHTTPHandler handles HTTP requests for Locations
type LocationHTTPHandler struct {
	locationService LocationService
}

// NewLocationHTTPHandler initializes a new LocationHTTPHandler
func NewLocationHTTPHandler(locationService LocationService) *LocationHTTPHandler {
	return &LocationHTTPHandler{
		locationService: locationService,
	}
}

// HandleHTTPPost handles the creation of a new Location
//
//	@Summary		Create a new Location
//	@Description	Create a new Location
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			location	body		Location	true	"Location data"
//	@Success		201			{object}	Location
//	@Failure		400			{object}	LocationError
//	@Failure		500			{object}	LocationError
//	@Router			/location [post]
func (aH *LocationHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var location Location
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newLocation, err := aH.locationService.Create(location)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet handles retrieving all Locations
//
//	@Summary		Get all Locations
//	@Description	Get all Locations
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		Location
//	@Failure		400	{object}	LocationError
//	@Failure		500	{object}	LocationError
//	@Router			/locations [get]
func (aH *LocationHTTPHandler) HandleHTTPGet(w http.ResponseWriter, r *http.Request) {
	locations, err := aH.locationService.ReadAll()
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGetWithID handles retrieving a single Location by ID
//
//	@Summary		Get a Location by ID
//	@Description	Get a specific Location by ID
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Location ID"
//	@Success		200	{object}	Location
//	@Failure		400	{object}	LocationError
//	@Failure		404	{object}	LocationError
//	@Failure		500	{object}	LocationError
//	@Router			/location/{id} [get]
func (aH *LocationHTTPHandler) HandleHTTPGetWithID(w http.ResponseWriter, r *http.Request) {
	ids := r.PathValue("ids")

	idList := strings.Split(ids, ",")
	var intIDs []int
	for _, id := range idList {
		intID, err := strconv.Atoi(id)
		if err != nil {
			aH.errorResponse(w, http.StatusBadRequest, "Invalid ID format")
			return
		}
		intIDs = append(intIDs, intID)
	}
	locations, err := aH.locationService.Read(intIDs)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(locations) == 0 {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(locations)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPPut handles updating an existing Location by ID
//
//	@Summary		Update an existing Location
//	@Description	Update an existing Location by ID
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string		true	"Location ID"
//	@Param			location	body		Location	true	"Updated Location data"
//	@Success		200			{object}	Location
//	@Failure		400			{object}	LocationError
//	@Failure		404			{object}	LocationError
//	@Failure		500			{object}	LocationError
//	@Router			/location/{id} [put]
func (aH *LocationHTTPHandler) HandleHTTPPut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedLocation Location
	err := json.NewDecoder(r.Body).Decode(&updatedLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	location, found, err := aH.locationService.Update(id, updatedLocation)
	if err != nil {
		aH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !found {
		aH.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(location)
	if err != nil {
		aH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPDelete handles deleting a Location by ID
//
//	@Summary		Delete a Location
//	@Description	Delete a Location by ID
//	@Tags			locations
//	@Param			id	path	string	true	"Location ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	LocationError
//	@Failure		404	{object}	LocationError
//	@Failure		500	{object}	LocationError
//	@Router			/location/{id} [delete]
func (aH *LocationHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	found, err := aH.locationService.Delete(id)
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

// errorResponse sends an error response with the specified status code and message
func (aH *LocationHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(LocationError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
