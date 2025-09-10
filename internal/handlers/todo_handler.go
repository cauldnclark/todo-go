package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cauldnclark/todo-go/internal/middleware"
	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	todoService *service.TodoService
	userService *service.UserService
}

func NewTodoHandler(todoService *service.TodoService, userService *service.UserService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		userService: userService,
	}
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	completed := r.URL.Query().Get("completed")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	var completedBool *bool
	switch completed {
	case "true":
		completedBool = new(bool)
		*completedBool = true
	case "false":
		completedBool = new(bool)
		*completedBool = false
	default:
		completedBool = nil
	}

	todoPage, err := h.todoService.GetTodos(r.Context(), userID, completedBool, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoPage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.todoService.CreateTodo(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, "Failed to create todo "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Todo created successfully"}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	todoIDStr := chi.URLParam(r, "id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTodoRequest
	if errDecode := json.NewDecoder(r.Body).Decode(&req); errDecode != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.UpdateTodo(r.Context(), userID, todoID, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	todoIDStr := chi.URLParam(r, "id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.todoService.DeleteTodo(r.Context(), userID, todoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	todoIDStr := chi.URLParam(r, "id")
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetTodoByID(r.Context(), todoID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
