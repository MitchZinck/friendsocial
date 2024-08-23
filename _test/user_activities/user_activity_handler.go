package test_user_activities

import (
	"bytes"
	"encoding/json"
	"friendsocial/user_activities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserActivityService(ctrl)
	handler := user_activities.NewUserActivityHTTPHandler(mockService)

	newUserActivity := user_activities.UserActivity{
		UserID:     1,
		ActivityID: 2,
		IsActive:   true,
	}

	mockService.EXPECT().Create(gomock.Any()).Return(newUserActivity, nil)

	body, _ := json.Marshal(newUserActivity)
	req, err := http.NewRequest(http.MethodPost, "/user_activities", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.HandleHTTPPost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var createdUserActivity user_activities.UserActivity
	err = json.NewDecoder(rr.Body).Decode(&createdUserActivity)
	assert.NoError(t, err)
	assert.Equal(t, newUserActivity, createdUserActivity)
}
