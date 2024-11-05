package main

// import (
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// const (
// 	MessageTypePosition = "position"
// 	MessageTypeChat     = "chat"
// )

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	roomManager = &RoomManager{
// 		rooms: make(map[string]map[*Client]bool),
// 	}

// 	chatRoomManager = &ChatRoomManager{
// 		rooms: make(map[string]map[*Client]bool),
// 	}
// )

// type Client struct {
// 	conn   *websocket.Conn
// 	room   string
// 	userId string
// 	mu     sync.Mutex
// }

// type Message struct {
// 	Type    string      `json:"type"`
// 	RoomId  string      `json:"room_id"`
// 	UserId  string      `json:"user_id"`
// 	Content interface{} `json:"content"`
// }

// type Position struct {
// 	PosX int `json:"pos_x"`
// 	PosY int `json:"pos_y"`
// }

// type ChatMessage struct {
// 	Text string `json:"text"`
// }

// type RoomManager struct {
// 	rooms map[string]map[*Client]bool
// 	mu    sync.RWMutex
// }

// func (rm *RoomManager) addClient(room string, client *Client) {
// 	rm.mu.Lock()
// 	defer rm.mu.Unlock()

// 	if rm.rooms[room] == nil {
// 		rm.rooms[room] = make(map[*Client]bool)
// 	}
// 	rm.rooms[room][client] = true
// }

// func (rm *RoomManager) removeClient(room string, client *Client) {
// 	rm.mu.Lock()
// 	defer rm.mu.Unlock()

// 	if clients, exists := rm.rooms[room]; exists {
// 		delete(clients, client)
// 		if len(clients) == 0 {
// 			delete(rm.rooms, room)
// 		}
// 	}
// }

// func (rm *RoomManager) broadcast(room string, message Message, sender *Client) {
// 	rm.mu.RLock()
// 	clients := rm.rooms[room]
// 	rm.mu.RUnlock()

// 	for client := range clients {
// 		if client != sender {
// 			client.mu.Lock()
// 			err := client.conn.WriteJSON(message)
// 			client.mu.Unlock()

// 			if err != nil {
// 				log.Printf("Error broadcasting to client: %v \n", err)
// 			}
// 		}
// 	}
// }

// type ChatRoomManager struct {
// 	rooms map[string]map[*Client]bool
// 	mu    sync.RWMutex
// }

// func (crm *ChatRoomManager) addClient(chatRoom string, client *Client) {
// 	crm.mu.Lock()
// 	defer crm.mu.Unlock()

// 	if crm.rooms[chatRoom] == nil {
// 		crm.rooms[chatRoom] = make(map[*Client]bool)
// 	}
// 	crm.rooms[chatRoom][client] = true
// }

// func (crm *ChatRoomManager) removeClient(chatRoom string, client *Client) {
// 	crm.mu.Lock()
// 	defer crm.mu.Unlock()

// 	if clients, exists := crm.rooms[chatRoom]; exists {
// 		delete(clients, client)
// 		if len(clients) == 0 {
// 			delete(crm.rooms, chatRoom)
// 		}
// 	}
// }

// func (crm *ChatRoomManager) broadcast(chatRoom string, message Message) {
// 	crm.mu.Lock()
// 	clients := crm.rooms[chatRoom]
// 	crm.mu.RUnlock()

// 	for client := range clients {
// 		client.mu.Lock()
// 		err := client.conn.WriteJSON(message)
// 		client.mu.Unlock()

// 		if err != nil {
// 			log.Printf("Error broadcasting to client: %v \n", err)
// 		}
// 	}
// }
