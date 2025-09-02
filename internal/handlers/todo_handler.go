package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cauldnclark/todo-go/internal/middleware"
	"github.com/cauldnclark/todo-go/internal/service"
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
