package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message types for communication
type MessageType string

const (
	MsgEvaluate MessageType = "evaluate"
	MsgResult   MessageType = "result"
	MsgError    MessageType = "error"
	MsgLog      MessageType = "log"
	MsgStatus   MessageType = "status"
)

// WebSocket message structure
type WSMessage struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content"`
	Time    time.Time   `json:"time"`
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Connected clients
type Client struct {
	conn   *websocket.Conn
	send   chan WSMessage
	active bool
}

var clients = make(map[*Client]bool)
var broadcast = make(chan WSMessage)

func main() {
	// Initialize audio engine
	if err := initAudioEngine(); err != nil {
		log.Fatalf("Failed to initialize audio engine: %v", err)
	}
	
	r := mux.NewRouter()
	
	// WebSocket endpoint
	r.HandleFunc("/ws", handleWebSocket)
	
	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Strudel Backend Server is running")
	})
	
	// Start the broadcast handler
	go handleBroadcast()
	
	log.Println("🎵 Strudel Backend Server starting on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws")
	log.Println("Health check: http://localhost:8080/health")
	
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	
	client := &Client{
		conn:   conn,
		send:   make(chan WSMessage, 256),
		active: true,
	}
	
	clients[client] = true
	
	log.Printf("Client connected. Total clients: %d", len(clients))
	
	// Send welcome message
	welcomeMsg := WSMessage{
		Type:    MsgStatus,
		Content: "Connected to Strudel backend!",
		Time:    time.Now(),
	}
	client.send <- welcomeMsg
	
	// Start goroutines for this client
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		delete(clients, c)
		c.conn.Close()
		log.Printf("Client disconnected. Total clients: %d", len(clients))
	}()
	
	c.conn.SetReadLimit(1024)
	c.conn.SetReadDeadline(time.Now().Add(300 * time.Second)) // 5 minutes
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(300 * time.Second))
		return nil
	})
	
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		// Process the received Strudel pattern
		pattern := string(message)
		log.Printf("Received pattern: %s", pattern)
		
		// Evaluate the pattern
		result := evaluateStrudelPattern(pattern)
		
		// Send result back to client
		response := WSMessage{
			Type:    MsgResult,
			Content: result,
			Time:    time.Now(),
		}
		
		select {
		case c.send <- response:
		default:
			close(c.send)
			return
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			// Send JSON message
			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		for client := range clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

