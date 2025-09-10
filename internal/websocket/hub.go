package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/cauldnclark/todo-go/internal/models"
	"github.com/cauldnclark/todo-go/internal/redis"
)

type Hub struct {
	// Map: UserID → Set of connected clients (for that user)
	clients map[int]map[*Client]bool

	// Channels for communication
	Broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	shutdown   chan struct{}
	wg         sync.WaitGroup
	mu         sync.RWMutex

	// Redis integration for horizontal scaling
	redisClient *redis.Client
	ctx         context.Context
}

func NewHub(redisClient *redis.Client) *Hub {
	hub := &Hub{
		clients:     make(map[int]map[*Client]bool),
		Broadcast:   make(chan Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		shutdown:    make(chan struct{}),
		redisClient: redisClient,
		ctx:         context.Background(),
	}

	if redisClient != nil {
		go hub.redisSubscribe()
	}

	go hub.Run()
	return hub
}

func (h *Hub) Run() {
	defer h.wg.Done()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		case <-h.shutdown:
			log.Println("Hub shutting down")
			h.cleanup()
			return
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.clients[client.UserID]
	if !ok {
		clients = make(map[*Client]bool)
		h.clients[client.UserID] = clients
	}

	clients[client] = true
	log.Printf("Client registered: UserID=%d, TotalClients=%d", client.UserID, len(clients))
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.clients[client.UserID]
	if ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)
			log.Printf("Client unregistered: UserID=%d, RemainingClients=%d", client.UserID, len(clients))
		}
		if len(clients) == 0 {
			delete(h.clients, client.UserID)
		}
	}
}

func (h *Hub) broadcastMessage(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var targetUserID int
	switch v := message.Data.(type) {
	case models.Todo:
		targetUserID = v.UserID
	case map[string]interface{}:
		if uid, ok := v["user_id"].(int); ok {
			targetUserID = uid
		}
	default:
		log.Printf("⚠️  Cannot extract UserID from message data: %+v", message.Data)
		return
	}

	clients, ok := h.clients[targetUserID]
	if !ok || len(clients) == 0 {
		log.Printf("No clients connected for UserID=%d", targetUserID)
		return
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Failed to marshal broadcast message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- msgBytes:
			// Message sent successfully
			log.Printf("Message sent to UserID=%d", targetUserID)
		default:
			log.Printf("❌ Failed to send message to UserID=%d, closing connection", targetUserID)
			delete(clients, client)
			close(client.send)

			if len(clients) == 0 {
				delete(h.clients, targetUserID)
			}
		}
	}
}

func (h *Hub) redisSubscribe() {
	defer h.wg.Done()

	if h.redisClient == nil {
		return
	}

	pubsub := h.redisClient.GetClient().Subscribe(h.ctx, "todo-updates")
	defer pubsub.Close()

	log.Println("Subscribed to Redis channel: todo-updates")

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			var message Message
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Printf("❌ Failed to unmarshal Redis message: %v", err)
				continue
			}
			select {
			case h.Broadcast <- message:
			case <-h.ctx.Done():
				return
			}
		case <-h.ctx.Done():
			log.Println("Redis subscription shutting down")
			return
		}
	}
}

func (h *Hub) cleanup() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Cleaning up all clients")
	for userID, clients := range h.clients {
		for client := range clients {
			select {
			case <-client.send:
				log.Printf("Client already closed: UserID=%d", userID)
			default:
				log.Printf("Closing client connection: UserID=%d", userID)
			}
			close(client.send)

			if client.conn != nil {
				client.conn.Close()
			}

			log.Printf("Client cleaned up: UserID=%d", userID)
		}
	}
}

func (h *Hub) PublishToRedis(event, channel string, data interface{}) error {
	if h.redisClient == nil {
		return nil
	}

	msg := Message{
		Event: event,
		Data:  data,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return h.redisClient.GetClient().Publish(ctx, channel, bytes).Err()
}
