package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cauldnclark/todo-go/internal/cache"
	"github.com/cauldnclark/todo-go/internal/config"
	"github.com/cauldnclark/todo-go/internal/handlers"
	"github.com/cauldnclark/todo-go/internal/middleware"
	"github.com/cauldnclark/todo-go/internal/redis"
	"github.com/cauldnclark/todo-go/internal/repository"
	"github.com/cauldnclark/todo-go/internal/service"
	"github.com/cauldnclark/todo-go/internal/websocket"
	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbpool, err := config.NewPGXConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	redisHost := cfg.Redis.Host
	redisPort := cfg.Redis.Port
	isProd := cfg.Server.IsProd
	var redisPassword string
	if isProd {
		redisPassword = cfg.Redis.Password
	} else {
		redisPassword = ""
	}

	redisAddr := redisHost + ":" + redisPort
	redisClient, err := redis.NewClient(redisAddr, &redisPassword, 0)

	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}

	log.Println("Redis client initialized")

	redisCache := cache.NewRedisCache(redisClient)
	log.Println("Redis cache initialized")

	hub := websocket.NewHub(redisClient)
	go hub.Run() // Start the hub to handle WebSocket connections
	wsHandler := websocket.NewHandler(hub)
	defer dbpool.Close()

	userRepo := repository.NewUserRepository(dbpool)
	todoRepo := repository.NewTodoRepository(dbpool)

	userService := service.NewUserService(userRepo, cfg.Server.JWTSecret, cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)
	todoService := service.NewTodoService(todoRepo, redisCache, hub)

	authHandler := handlers.NewAuthHandler(userService)
	todoHandler := handlers.NewTodoHandler(todoService, userService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.Server.JWTSecret)

	r := chi.NewRouter()

	r.Use(chimiddle.Logger)
	r.Use(chimiddle.Recoverer)
	r.Use(chimiddle.RequestID)
	r.Use(chimiddle.RealIP)
	r.Use(chimiddle.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/ws", wsHandler.ServeHTTP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/google", authHandler.GoogleSignIn)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Route("/todos", func(r chi.Router) {
			r.Get("/", todoHandler.GetTodos)
			r.Get("/{id}", todoHandler.GetTodoByID)
			r.Post("/", todoHandler.CreateTodo)
			r.Put("/{id}", todoHandler.UpdateTodo)
			r.Delete("/{id}", todoHandler.DeleteTodo)
		})

		r.Get("/me", authHandler.GetCurrentUser)
	})

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logrus.Infof("Server is running on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Error shutting down server: %v", err)
	}
}
