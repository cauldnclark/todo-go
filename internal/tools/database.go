package tools

import "time"

type UserDetails struct {
	ID        int
	Name      string
	Email     string
	GoogleID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type TodoDetails struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DataInterface interface {
	CreateUser(user *UserDetails) error
	GetUser(id int) (*UserDetails, error)
	UpdateUser(user *UserDetails) error
	DeleteUser(id int) error
	CreateTodo(todo *TodoDetails) error
	GetTodo(id int) (*TodoDetails, error)
	UpdateTodo(todo *TodoDetails) error
	DeleteTodo(id int) error
	GetTodos(userID int) ([]*TodoDetails, error)
	SetupDatabase() error
}
