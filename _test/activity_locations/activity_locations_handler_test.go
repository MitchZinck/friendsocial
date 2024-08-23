package test_activity_locations

import (
	"bytes"
	"encoding/json"
	"friendsocial/activity_locations"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockActivityLocationService(ctrl)
	handler := activity_locations.NewActivityLocationHTTPHandler(mockService)

	newLocation := activity_locations.ActivityLocation{
		Name:      "Park",
		Address:   "123 Green St",
		City:      "Greenwood",
		State:     "GW",
		ZipCode:   "12345",
		Country:   "Wonderland",
		Latitude:  1.2345,
		Longitude: 6.7890,
	}

	mockService.EXPECT().Create(gomock.Any()).Return(newLocation, nil)

	body, _ := json.Marshal(newLocation)
	req, err := http.NewRequest(http.MethodPost, "/activity-locations", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.HandleHTTPPost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var createdLocation activity_locations.ActivityLocation
	err = json.NewDecoder(rr.Body).Decode(&createdLocation)
	assert.NoError(t, err)
	assert.Equal(t, newLocation, createdLocation)
}
