package test_activities

import (
	"bytes"
	"encoding/json"
	"friendsocial/activities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockActivityService(ctrl)
	handler := activities.NewActivityHTTPHandler(mockService)

	newActivity := activities.Activity{
		Name:          "Golf",
		Description:   "A game where you hit a golf ball on a golf course.",
		EstimatedTime: "60 Minutes",
		LocationID:    1,
	}

	mockService.EXPECT().Create(gomock.Any()).Return(newActivity, nil)

	body, _ := json.Marshal(newActivity)
	req, err := http.NewRequest(http.MethodPost, "/activities", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.HandleHTTPPost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var createdActivity activities.Activity
	err = json.NewDecoder(rr.Body).Decode(&createdActivity)
	assert.NoError(t, err)
	assert.Equal(t, newActivity, createdActivity)
}
