package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo           *repository.UserRepository
	jwtSecret          string
	googleClientId     string
	googleClientSecret string
	redirectURI        string
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret, googleClientId, googleClientSecret, redirectURI string) *UserService {
	return &UserService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		googleClientId:     googleClientId,
		googleClientSecret: googleClientSecret,
		redirectURI:        redirectURI,
	}
}

func (s *UserService) AuthenticateWithGoogle(ctx context.Context, tokenString string) (*models.AuthResponse, error) {
	// Exchange authorization code for access token
	tokenResp, err := s.exchangeCodeForToken(tokenString, s.redirectURI)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	googleUser, err := s.getGoogleUserInfo(tokenResp.AccessToken)
	if err != nil {
		log.Printf("Error getting user info from Google: %v", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.userRepo.GetUserByGoogleID(ctx, googleUser.ID)
	if err != nil {
		log.Println("Swallowing this error to create new user:", err)
	}

	if user == nil {
		user = &models.User{
			GoogleID:   googleUser.ID,
			Email:      googleUser.Email,
			Name:       googleUser.Name,
			PictureURL: googleUser.PictureURL,
		}
		if errCreate := s.userRepo.CreateUser(ctx, user); errCreate != nil {
			return nil, fmt.Errorf("failed to create user: %w", errCreate)
		}
	} else {
		if user.Email != googleUser.Email || user.Name != googleUser.Name || user.PictureURL != googleUser.PictureURL {
			user.Email = googleUser.Email
			user.Name = googleUser.Name
			user.PictureURL = googleUser.PictureURL
			if errUpdate := s.userRepo.UpdateUser(ctx, user); errUpdate != nil {
				return nil, fmt.Errorf("failed to update user: %w", errUpdate)
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

func (s *UserService) getGoogleUserInfo(accessToken string) (*models.GoogleTokenInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google userinfo request failed: %s", string(body))
	}

	var userInfo models.GoogleTokenInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (s *UserService) exchangeCodeForToken(code, redirectURI string) (*models.GoogleTokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"
	log.Printf("Making token exchange request to: %s", tokenURL)

	data := url.Values{}
	data.Set("client_id", s.googleClientId)
	data.Set("client_secret", s.googleClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	log.Printf("Token request data: client_id=%s, redirect_uri=%s, code=%s", s.googleClientId, redirectURI, code)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Token exchange response status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Token exchange failed with status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("google token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp models.GoogleTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		log.Printf("Failed to decode token response: %v", err)
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	log.Printf("Token exchange successful, access_token length: %d", len(tokenResp.AccessToken))
	return &tokenResp, nil
}

func (s *UserService) generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
