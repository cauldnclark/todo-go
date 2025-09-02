package models

import "time"

type User struct {
	ID         int       `json:"id" db:"id"`
	GoogleID   string    `json:"google_id" db:"google_id"`
	Email      string    `json:"email" db:"email"`
	Name       string    `json:"name" db:"name"`
	PictureURL string    `json:"picture_url" db:"picture_url"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Todo struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type MetaPagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type TodosPaginated struct {
	Todos []Todo         `json:"todos"`
	Meta  MetaPagination `json:"meta"`
}

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   *bool  `json:"completed"`
}

type GoogleTokenInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	PictureURL string `json:"picture"`
}

type AuthRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
