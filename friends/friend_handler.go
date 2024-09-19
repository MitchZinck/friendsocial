package friends

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// FriendService defines the interface for the friend service
type FriendService interface {
	Create(userID string, friendID string) (Friend, error)
	ReadByUserID(userID string) ([]Friend, error)
	ReadByFriendID(friendID string) ([]Friend, error)
	UsersAreFriends(userID string, friendID string) (bool, error)
	Delete(userID string, friendID string) (bool, error)
}

// FriendError represents an error response
type FriendError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

// FriendHTTPHandler is the HTTP handler for friend-related operations
type FriendHTTPHandler struct {
	friendService FriendService
}

// NewFriendHTTPHandler creates a new FriendHTTPHandler
func NewFriendHTTPHandler(friendService FriendService) *FriendHTTPHandler {
	return &FriendHTTPHandler{
		friendService: friendService,
	}
}

// HandleHTTPPost creates a new friendship
//
//	@Summary		Create a new friendship
//	@Description	Create a new friendship between two users
//	@Tags			friends
//	@Accept			json
//	@Produce		json
//	@Param			friend	body		Friend	true	"Friendship information"
//	@Success		201		{object}	Friend
//	@Failure		400		{object}	FriendError
//	@Failure		500		{object}	FriendError
//	@Router			/friends [post]
func (fH *FriendHTTPHandler) HandleHTTPPost(w http.ResponseWriter, r *http.Request) {
	var friend Friend
	err := json.NewDecoder(r.Body).Decode(&friend)
	if err != nil {
		fH.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	friendIDStr := strconv.Itoa(friend.UserID)
	friendFriendIDStr := strconv.Itoa(friend.FriendID)

	newFriend, err := fH.friendService.Create(friendIDStr, friendFriendIDStr)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newFriend)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// HandleHTTPGet retrieves all friends of a user
//
//	@Summary		Get all friends of a user
//	@Description	Retrieve all friendships for a given user
//	@Tags			friends
//	@Produce		json
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{array}		Friend
//	@Failure		500		{object}	FriendError
//	@Router			/friends/{user_id} [get]
func (fH *FriendHTTPHandler) HandleHTTPGetByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	friends, err := fH.friendService.ReadByUserID(userID)
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

// HandleHTTPGetByFriendID checks if a specific friendship exists
//
//	@Summary		Check if a specific friendship exists
//	@Description	Check if a friendship exists between two users
//	@Tags			friends
//	@Produce		json
//	@Param			user_id		path	string	true	"User ID"
//	@Param			friend_id	path	string	true	"Friend ID"
//	@Success		200
//	@Failure		404	{object}	FriendError
//	@Failure		500	{object}	FriendError
//	@Router			/friends/{user_id}/{friend_id} [get]
func (fH *FriendHTTPHandler) HandleHTTPGetByFriendID(w http.ResponseWriter, r *http.Request) {
	friendID := r.PathValue("friend_id")

	exists, err := fH.friendService.ReadByFriendID(friendID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(exists) == 0 {
		fH.errorResponse(w, http.StatusNotFound, "Friendship not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleHTTPDelete deletes a friendship
//
//	@Summary		Delete a friendship
//	@Description	Delete an existing friendship between two users
//	@Tags			friends
//	@Param			user_id		path	string	true	"User ID"
//	@Param			friend_id	path	string	true	"Friend ID"
//	@Success		204
//	@Failure		404	{object}	FriendError
//	@Failure		500	{object}	FriendError
//	@Router			/friends/{user_id}/{friend_id} [delete]
func (fH *FriendHTTPHandler) HandleHTTPDelete(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	friendID := r.PathValue("friend_id")

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

// HandleHTTPGetAreFriends checks if two users are friends
//
//	@Summary		Check if two users are friends
//	@Description	Check if two users are friends
//	@Tags			friends
//	@Produce		json
//	@Param			user_id		path	string	true	"User ID"
func (fH *FriendHTTPHandler) HandleHTTPGetAreFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	friendID := r.PathValue("friend_id")

	areFriends, err := fH.friendService.UsersAreFriends(userID, friendID)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(areFriends)
	if err != nil {
		fH.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// errorResponse sends a JSON error response
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
