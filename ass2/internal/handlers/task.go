package handlers

import (
	"ass2/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type TaskHandler struct {
	tasks  []models.Task
	nextID int
	mu     sync.RWMutex
	client *http.Client
}

const MaxTitleLength = 200

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks:  make([]models.Task, 0),
		nextID: 1,
		client: &http.Client{},
	}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		idParam := r.URL.Query().Get("id")
		if idParam != "" {
			h.getTaskByID(w, r, idParam)
		} else {
			h.getAllTasks(w, r)
		}
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodPatch:
		h.updateTask(w, r)
	case http.MethodDelete:
		h.deleteTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /tasks [get]
func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, idParam string) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id: must be a valid integer"})
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, task := range h.tasks {
		if task.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get all tasks with optional filtering by done status
// @Tags tasks
// @Accept json
// @Produce json
// @Param done query boolean false "Filter by done status"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /tasks [get]
func (h *TaskHandler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	doneParam := r.URL.Query().Get("done")
	var filteredTasks []models.Task

	if doneParam == "" {
		filteredTasks = h.tasks
	} else {
		doneValue, err := strconv.ParseBool(doneParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid done parameter: must be true or false"})
			return
		}

		for _, task := range h.tasks {
			if task.Done == doneValue {
				filteredTasks = append(filteredTasks, task)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredTasks)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with a title
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body object{title=string} true "Task title (max 200 characters)"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /tasks [post]
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body: must be valid JSON"})
		return
	}

	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid title: title cannot be empty"})
		return
	}

	if len(req.Title) > MaxTitleLength {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid title: title exceeds maximum length of 200 characters"})
		return
	}

	h.mu.Lock()
	task := models.Task{
		ID:    h.nextID,
		Title: req.Title,
		Done:  false,
	}
	h.tasks = append(h.tasks, task)
	h.nextID++
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask godoc
// @Summary Update task status
// @Description Update the done status of a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query int true "Task ID"
// @Param task body object{done=boolean} true "Task status"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /tasks [patch]
func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing required parameter: id"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id: must be a valid integer"})
		return
	}

	var req struct {
		Done *bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body: must be valid JSON"})
		return
	}

	if req.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing required field: done"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.tasks {
		if h.tasks[i].ID == id {
			h.tasks[i].Done = *req.Done
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]bool{"updated": true})
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /tasks [delete]
func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing required parameter: id"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id: must be a valid integer"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, task := range h.tasks {
		if task.ID == id {
			h.tasks = append(h.tasks[:i], h.tasks[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "task deleted successfully"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
}
