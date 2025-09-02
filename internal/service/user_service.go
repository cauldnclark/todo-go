package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo       *repository.UserRepository
	jwtSecret      string
	googleClientId string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret, googleClientId string) *UserService {
	return &UserService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		googleClientId: googleClientId,
	}
}

func (s *UserService) AuthenticateWithGoogle(ctx context.Context, tokenString string) (*models.AuthResponse, error) {
	googleUser, err := s.verifyGoogleToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to verify google token: %w", err)
	}

	user, err := s.userRepo.GetUserByGoogleID(ctx, googleUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by google id: %w", err)
	}

	if user == nil {
		user = &models.User{
			GoogleID:   googleUser.ID,
			Email:      googleUser.Email,
			Name:       googleUser.Name,
			PictureURL: googleUser.PictureURL,
		}
		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		if user.Email != googleUser.Email || user.Name != googleUser.Name || user.PictureURL != googleUser.PictureURL {
			user.Email = googleUser.Email
			user.Name = googleUser.Name
			user.PictureURL = googleUser.PictureURL
			if err := s.userRepo.UpdateUser(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate jwt: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}

func (s *UserService) ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid user id")
		}
		return int(userID), nil
	}

	return 0, fmt.Errorf("invalid token")
}

func (s *UserService) verifyGoogleToken(tokenString string) (*models.GoogleTokenInfo, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", tokenString)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var googleUser models.GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (s *UserService) generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
