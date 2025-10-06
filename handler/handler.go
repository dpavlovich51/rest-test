package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	e "my_rest_server/error"
	m "my_rest_server/model"
	s "my_rest_server/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type response struct {
	Data  interface{}
	Error string
}

type UserHandler struct {
	storage *s.UserService
}

func NewInstance(store *s.UserService) *UserHandler {
	return &UserHandler{storage: store}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.storage.GetAllUsers(r.Context())
	if err != nil {
		sendError(w, err)
		return
	}
	sendResponseOk(w, *users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		sendError(w, e.NewError2(http.StatusBadRequest, "failed to find id", err))
		return
	}
	user, err := h.storage.GetUser(r.Context(), userId)
	if err != nil {
		sendError(w, e.NewError2(http.StatusBadRequest, "failed to find id", err))
		return
	}
	sendResponseOk(w, []m.User{*user})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		sendError(w, e.NewError2(http.StatusBadRequest, "parameter: userId is missing", err))
		return
	}
	// isExists, err := h.storage.IsUserExists(r.Context(), userId)
	// if err != nil {
	// 	sendError(w, err, http.StatusInternalServerError)
	// 	return
	// }
	// if !isExists {
	// 	sendError(w, fmt.Errorf("user: %s does not exist", userId), http.StatusInternalServerError)
	// 	return
	// }
	var oldUser *m.User
	oldUser, err = h.storage.GetUser(r.Context(), userId)
	if err != nil {
		sendError(w, err)
		return
	}
	h.saveUserWithFields(w, r, func(user *m.User) {
		(*user).Id = oldUser.Id
		(*user).CreatedAt = oldUser.CreatedAt
	})
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	h.saveUserWithFields(w, r, func(user *m.User) {
		(*user).Id = uuid.NewString()
		(*user).CreatedAt = time.Now()
	})
}

func (h *UserHandler) saveUserWithFields(
	w http.ResponseWriter,
	r *http.Request,
	override func(user *m.User),
) {
	var user m.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendError(w, e.NewError2(http.StatusBadRequest, "failed to parse json data", err))
	}
	override(&user)

	_, err = h.storage.SaveUser(r.Context(), user)
	if err != nil {
		sendError(w, err)
		return
	}
	sendResponseOk(w, []m.User{user})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		sendError(w, e.NewError2(http.StatusBadRequest, "parameter: userId is missing", err))
		return
	}
	deletedUser, err := h.storage.DeleteUser(r.Context(), userId)
	if err != nil {
		sendError(w, err)
		return
	}
	sendResponseOk(w, []m.User{*deletedUser})
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

func sendError(w http.ResponseWriter, err error) {
	var domainErr *e.DomainError

	setJsonContent(w)
	if errors.As(err, &domainErr) {
		w.WriteHeader(domainErr.Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response{Error: err.Error()})
}

func setJsonContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
