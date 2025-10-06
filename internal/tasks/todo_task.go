package tasks

import (
	"context"
	"encoding/json"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/service"
	"github.com/hibiken/asynq"
)

const (
	CreateTodoTask = "create_todo"
)

type CreateTodoTaskPayload struct {
	UserID      int
	Title       string
	Description string
	Completed   bool
}

func NewCreateTodoTask(payload CreateTodoTaskPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(CreateTodoTask, data), nil
}

type TaskHandler struct {
	todoService *service.TodoService
}

func NewTaskHandler(todoService *service.TodoService) *TaskHandler {
	return &TaskHandler{
		todoService: todoService,
	}
}

func (h *TaskHandler) HandleCreateTodoTask(ctx context.Context, t *asynq.Task) error {
	var payload CreateTodoTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	req := &models.CreateTodoRequest{
		Title:       payload.Title,
		Description: payload.Description,
	}

	return h.todoService.CreateTodo(ctx, payload.UserID, req)
}
