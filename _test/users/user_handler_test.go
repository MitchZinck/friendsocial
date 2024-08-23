package test_users

import (
	"bytes"
	"encoding/json"
	"friendsocial/users"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserService(ctrl)
	handler := users.NewUserHTTPHandler(mockService)

	newUser := users.User{
		Name:     "Mitchell Zinck22",
		Email:    "mitchellf22zinck@gmail.com",
		Password: "test123",
	}

	mockService.EXPECT().Create(gomock.Any()).Return(newUser, nil)

	body, _ := json.Marshal(newUser)
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.HandleHTTPPost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var createdUser users.User
	err = json.NewDecoder(rr.Body).Decode(&createdUser)
	assert.NoError(t, err)
	assert.Equal(t, newUser, createdUser)
}
