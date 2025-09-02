package service

import (
	"context"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/repository"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID int, req *models.CreateTodoRequest) (*models.Todo, error) {
	todo := &models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	err := s.todoRepo.CreateTodo(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todoID, userID int, req *models.UpdateTodoRequest) (*models.Todo, error) {
	todo, err := s.todoRepo.GetTodoByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	err = s.todoRepo.UpdateTodo(ctx, todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *TodoService) GetTodos(ctx context.Context, userID int, completed *bool, page, limit int) (*models.TodosPaginated, error) {
	todosPage, err := s.todoRepo.GetTodosPaginated(ctx, userID, completed, page, limit)
	if err != nil {
		return nil, err
	}
	return todosPage, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, todoID, userID int) error {
	return s.todoRepo.DeleteTodo(ctx, todoID, userID)
}
