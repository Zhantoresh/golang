package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"golang/internal/usecase"
	"golang/pkg/modules"
)

type UserHandler struct {
	uc *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// /users
func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		users, err := h.uc.GetUsers()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, users)

	case http.MethodPost:
		var u modules.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		id, err := h.uc.CreateUser(u)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"id": id})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// /users/{id}
func (h *UserHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {

	case http.MethodGet:
		user, err := h.uc.GetUserByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "user not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, user)

	case http.MethodPut, http.MethodPatch:
		var u modules.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		rows, err := h.uc.UpdateUser(id, u)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if rows == 0 {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "updated"})

	case http.MethodDelete:
		rows, err := h.uc.DeleteUserByID(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if rows == 0 {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": rows})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}