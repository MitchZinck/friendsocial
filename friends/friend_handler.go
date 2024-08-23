package friends

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FriendService interface {
	Create(userID int, friendID int) (Friend, error)
	Delete(userID int, friendID int) (bool, error)
	ReadAll(userID int) ([]Friend, error)
	Read(userID int, friendID int) (bool, error)
}

type FriendError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type FriendHTTPHandler struct {
	friendService FriendService
}

func NewFriendHTTPHandler(friendService FriendService) *FriendHTTPHandler {
	return &FriendHTTPHandler{
		friendService: friendService,
	}
}

func (fH *FriendHTTPHandler) HandleHTTPPostAddFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	friendID, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid friend_id")
		return
	}

	friend, err := fH.friendService.Create(userID, friendID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(friend)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (fH *FriendHTTPHandler) HandleHTTPDeleteRemoveFriend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	friendID, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid friend_id")
		return
	}

	found, err := fH.friendService.Delete(userID, friendID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !found {
		fH.errorResponse(w, http.StatusNotFound, "Friendship not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fH *FriendHTTPHandler) HandleHTTPGetFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	friends, err := fH.friendService.ReadAll(userID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(friends)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (fH *FriendHTTPHandler) HandleHTTPGetIsFriends(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid user_id")
		return
	}

	friendID, err := strconv.Atoi(vars["friend_id"])
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, "Invalid friend_id")
		return
	}

	isFriend, err := fH.friendService.Read(userID, friendID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]bool{"is_friend": isFriend})
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (fH *FriendHTTPHandler) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(FriendError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}
