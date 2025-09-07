package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, todo.UserID, todo.Title, todo.Description, todo.Completed).Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// paginated search of todos
func (r *TodoRepository) GetTodosPaginated(ctx context.Context, userID int, completed *bool, page, limit int) (*models.TodosPaginated, error) {
	query := `
			SELECT id, user_id, title, description, completed, created_at, updated_at
			FROM todos
			WHERE user_id = $1
			AND completed = $2
			ORDER BY id
			LIMIT $3
			OFFSET $4
		`
	if completed == nil {
		query = `
			SELECT id, user_id, title, description, completed, created_at, updated_at
			FROM todos
			WHERE user_id = $1
			ORDER BY id
			LIMIT $2
			OFFSET $3
		`
	}
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		errScan := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if errScan != nil {
			return nil, errScan
		}
		todos = append(todos, todo)
	}
	if errRows := rows.Err(); errRows != nil {
		return nil, errRows
	}

	// get total count
	query = `
		SELECT COUNT(*)
		FROM todos
		WHERE user_id = $1
		AND completed = $2
	`
	var total int
	err = r.db.QueryRow(ctx, query, userID, completed).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.TodosPaginated{
		Todos: todos,
		Meta: models.MetaPagination{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}, nil
}

func (r *TodoRepository) GetTodoByID(ctx context.Context, id, userID int) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $3, description = $4, completed = $5, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, todo.ID, todo.UserID, todo.Title, todo.Description, todo.Completed).
		Scan(&todo.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id, userID int) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}
