package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	m "my_rest_server/model"
	s "my_rest_server/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type response struct {
	Data  interface{}
	Error string
}

type UserHandler struct {
	storage *s.Cache
}

func NewInstance(store *s.Cache) *UserHandler {
	return &UserHandler{storage: store}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.storage.GetAllUsers(r.Context())
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	sendResponseOk(w, *users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		sendError(w, fmt.Errorf("failed to find id"), http.StatusBadRequest)
		return
	}
	user, err := h.storage.GetUser(r.Context(), userId)
	if err != nil {
		sendError(w, err, http.StatusNotFound)
		return
	}
	sendResponseOk(w, []m.User{user})
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	var user m.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, fmt.Errorf("failed to parse json data"), http.StatusBadRequest)
	}
	user.Id = uuid.NewString()
	userId, err := h.storage.SaveUser(r.Context(), user)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	user.Id = userId
	user.CreatedAt = time.Now()

	sendResponseOk(w, []m.User{user})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}
	deletedUser, err := h.storage.DeleteUser(r.Context(), userId)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	sendResponseOk(w, []m.User{deletedUser})
}

func getUserId(r *http.Request) (string, error) {
	paramName := "id"
	userId := mux.Vars(r)[paramName]
	if userId == "" {
		return "", fmt.Errorf("failed to find param: %s", paramName)
	}
	return userId, nil
}

// -----  -----  -----  -----  -----  -----

func sendResponseOk(w http.ResponseWriter, users []m.User) {
	setJsonContent(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{Data: users})
}

func sendError(w http.ResponseWriter, err error, status int) {
	setJsonContent(w)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response{Error: err.Error()})
}

func setJsonContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
