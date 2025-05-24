package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

type Client struct {
	nm   string
	conn *websocket.Conn
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		usrs:       make(map[*Client]string),
		srsu:       make(map[string]*Client),
		brdcast:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

type Hub struct {
	usrs       map[*Client]string
	srsu       map[string]*Client
	brdcast    chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver,omitempty"`
	Content  string `json:"content,omitempty"`
	IsTyping bool   `json:"isTyping,omitempty"`
}

type ChatMessage struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ChatHistoryRequest struct {
	Type   string `json:"type"`
	User   string `json:"user"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleClientRegister(client)

		case client := <-h.unregister:
			h.handleClientUnregister(client)
		case message := <-h.brdcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleClientRegister(client *Client) {
	h.mu.Lock()
	h.srsu[client.nm] = client
	h.usrs[client] = client.nm
	h.mu.Unlock()
	h.sendUserListUpdate()
}

func (h *Hub) handleClientUnregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.usrs[client]; ok {
		delete(h.srsu, client.nm)
		delete(h.usrs, client)
		close(client.send)
		// fmt.Println("yeeees")
	}
	h.mu.Unlock()
	h.sendUserListUpdate()
}

func (h *Hub) handleBroadcast(message []byte) {
	h.mu.Lock()
	clients := make([]*Client, 0, len(h.usrs))
	for client := range h.usrs {
		clients = append(clients, client)
	}
	h.mu.Unlock()

	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			go func(c *Client) {
				h.unregister <- c
				c.conn.Close()
			}(client)
		}
	}
}

func (h *Hub) sendUserListUpdate() {
	h.mu.Lock()
	userList := h.getOnlineUsersJSON()
	clients := make([]*Client, 0, len(h.usrs))
	for client := range h.usrs {
		clients = append(clients, client)
	}
	h.mu.Unlock()
	message := []byte(fmt.Sprintf(`{"type":"users_update","users":%s}`, userList))

	for _, client := range clients {
		select {
		case client.send <- message:
		default:
			go func(c *Client) {
				h.unregister <- c
				c.conn.Close()
			}(client)
		}
	}
}

func (h *Hub) getOnlineUsersJSON() string {

	users := make([]string, 0, len(h.srsu))
	for username := range h.srsu {
		users = append(users, fmt.Sprintf(`"%s"`, username))
	}
	return fmt.Sprintf(`[%s]`, strings.Join(users, ","))
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	cookie, err := r.Cookie("userId")
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(`{"error": "Not authenticated"}`))
		conn.Close()
		return
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE uuid = ?", cookie.Value).Scan(&username)
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(`{"error": "Invalid user"}`))
		conn.Close()
		return
	}

	client := &Client{
		nm:   username,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register <- client

	go client.writePump()
	go client.readPump(h)
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		switch msg["type"] {
		case "chat":
			var chatMsg Message
			if err := json.Unmarshal(message, &chatMsg); err == nil {
				h.handlePrivateMessage(chatMsg)
			}
		case "typing":
			var typingMsg Message
			if err := json.Unmarshal(message, &typingMsg); err == nil {
				h.handleTypingIndicator(typingMsg)
			}
		case "get_history":
			var req ChatHistoryRequest
			if err := json.Unmarshal(message, &req); err == nil {
				h.sendChatHistory(c, req.User, req.Offset, req.Limit)
			}
		}
	}
}

func (h *Hub) handleTypingIndicator(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if target, ok := h.srsu[msg.Receiver]; ok {
		typingMsg := fmt.Sprintf(`{
            "type": "typing",
            "sender": "%s",
            "isTyping": %t
        }`, msg.Sender, msg.IsTyping)
		target.send <- []byte(typingMsg)
	}
}

func (h *Hub) handlePrivateMessage(msg Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, content) VALUES ((SELECT id FROM users WHERE username = ?), (SELECT id FROM users WHERE username = ?), ?)",
		msg.Sender, msg.Receiver, msg.Content)
	if err != nil {
		log.Println("Error saving message:", err)
		return
	}

	messageData := fmt.Sprintf(`{
        "type": "chat",
        "sender": "%s",
        "receiver": "%s",
        "content": "%s",
        "created_at": "%s"
    }`, msg.Sender, msg.Receiver, msg.Content, time.Now().Format(time.RFC3339))

	if target, ok := h.srsu[msg.Receiver]; ok {
		target.send <- []byte(messageData)
	}

	if sender, ok := h.srsu[msg.Sender]; ok {
		sender.send <- []byte(messageData)
	}

	updateMsg := `{"type": "update_conversations"}`
	if target, ok := h.srsu[msg.Receiver]; ok {
		target.send <- []byte(updateMsg)
	}
	if sender, ok := h.srsu[msg.Sender]; ok {
		sender.send <- []byte(updateMsg)
	}
}

func (h *Hub) sendChatHistory(client *Client, otherUser string, offset, limit int) {
	var clientID, otherUserID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", client.nm).Scan(&clientID)
	if err != nil {
		log.Println("Error fetching client ID:", err)
		return
	}
	err = db.QueryRow("SELECT id FROM users WHERE username = ?", otherUser).Scan(&otherUserID)
	if err != nil {
		log.Println("Error fetching other user ID:", err)
		return
	}

	rows, err := db.Query(`
        SELECT DISTINCT  m.id, m.content, m.created_at, u.username 
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE (m.sender_id = ? AND m.receiver_id = ?)  
           OR (m.sender_id = ? AND m.receiver_id = ?)
        ORDER BY m.created_at DESC
        LIMIT ? OFFSET ?`,
		clientID, otherUserID, otherUserID, clientID, limit, offset)
	if err != nil {
		log.Println("Error fetching chat history:", err)
		return
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.CreatedAt, &msg.Sender); err != nil {
			log.Println("Error scanning message:", err)
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse the order so newest messages are at the bottom
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	historyJSON, err := json.Marshal(messages)
	if err != nil {
		log.Println("Error marshaling history:", err)
		return
	}
	// log.Println("Sending history:", len(messages), "messages")
	// fmt.Println(string(historyJSON))
	client.send <- []byte(fmt.Sprintf(`{
        "type": "chat_history",
        "user": "%s",
        "messages": %s
    }`, otherUser, string(historyJSON)))
}
