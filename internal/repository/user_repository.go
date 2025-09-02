package repository

import (
	"context"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (google_id, email, name, picture_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, user.GoogleID, user.Email, user.Name, user.PictureURL).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `
		SELECT id, google_id, email, name, picture_url, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`
	var user models.User
	err := r.db.QueryRow(ctx, query, googleID).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, google_id, email, name, picture_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.PictureURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET google_id = $1, email = $2, name = $3, picture_url = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, user.GoogleID, user.Email, user.Name, user.PictureURL, user.ID)
	if err != nil {
		return err
	}
	return nil
}
