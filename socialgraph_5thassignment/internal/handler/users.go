package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"socialgraph_5thassignment/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}

	filters := map[string]string{
		"id":         r.URL.Query().Get("id"),
		"name":       r.URL.Query().Get("name"),
		"email":      r.URL.Query().Get("email"),
		"gender":     r.URL.Query().Get("gender"),
		"birth_date": r.URL.Query().Get("birth_date"),
	}

	orderBy := r.URL.Query().Get("order_by")
	status := r.URL.Query().Get("status")

	resp, err := h.repo.GetPaginatedUsers(page, limit, filters, orderBy, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user1, err1 := strconv.Atoi(r.URL.Query().Get("user1"))
	user2, err2 := strconv.Atoi(r.URL.Query().Get("user2"))

	if err1 != nil || err2 != nil || user1 <= 0 || user2 <= 0 || user1 == user2 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid user ids"})
		return
	}

	users, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	if err := h.repo.SoftDeleteUser(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "soft deleted"})
}