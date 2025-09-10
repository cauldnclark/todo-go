package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/cauldnclark/todo-go/internal/cache"
	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/repository"
	"github.com/cauldnclark/todo-go/internal/websocket"
)

type TodoService struct {
	todoRepo *repository.TodoRepository
	cache    *cache.RedisCache
	hub      *websocket.Hub
}

func NewTodoService(todoRepo *repository.TodoRepository, cache *cache.RedisCache, hub *websocket.Hub) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		cache:    cache,
		hub:      hub,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID int, req *models.CreateTodoRequest) error {
	todo := &models.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	err := s.todoRepo.CreateTodo(ctx, todo)
	if err != nil {
		return err
	}

	cacheKey := "todos_user_" + strconv.Itoa(todo.UserID)
	s.cache.Delete(ctx, cacheKey)

	s.hub.Broadcast <- websocket.Message{
		Event: "todo.created",
		Data:  *todo,
	}

	return nil
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

func (s *TodoService) GetTodoByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	cacheKey := "todos_user_" + strconv.Itoa(todoID)

	var todo *models.Todo

	// try cache first
	err := s.cache.Get(ctx, cacheKey, &todo)
	if err == nil {
		log.Printf("cache hit for key: %s", cacheKey)
		return todo, nil
	}

	// if not in cache, get from db
	todo, err = s.todoRepo.GetTodoByID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}

	// cache the todo
	s.cache.Set(ctx, cacheKey, todo, time.Hour)
	log.Printf("cache set for key: %s", cacheKey)

	return todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, todoID, userID int) error {
	return s.todoRepo.DeleteTodo(ctx, todoID, userID)
}
