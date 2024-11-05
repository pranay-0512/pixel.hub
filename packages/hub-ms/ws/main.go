package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MessageTypePosition = "position"
	MessageTypeChat     = "chat"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	roomManager = &RoomManager{
		rooms: make(map[string]map[*Client]bool),
	}

	chatRoomManager = &ChatRoomManager{
		rooms: make(map[string]map[*Client]bool),
	}
)

type Client struct {
	conn   *websocket.Conn
	room   string
	userId string
	mu     sync.Mutex
}

type Message struct {
	Type    string      `json:"type"`
	RoomId  string      `json:"room_id"`
	UserId  string      `json:"user_id"`
	Content interface{} `json:"content"`
}

type Position struct {
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`
}

type ChatMessage struct {
	Text string `json:"text"`
}

type RoomManager struct {
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
}

func (rm *RoomManager) addClient(room string, client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.rooms[room] == nil {
		rm.rooms[room] = make(map[*Client]bool)
	}
	rm.rooms[room][client] = true
}

func (rm *RoomManager) removeClient(room string, client *Client) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if clients, exists := rm.rooms[room]; exists {
		delete(clients, client)
		if len(clients) == 0 {
			delete(rm.rooms, room)
		}
	}
}

func (rm *RoomManager) broadcast(room string, message Message, sender *Client) {
	rm.mu.RLock()
	clients := rm.rooms[room]
	rm.mu.RUnlock()

	for client := range clients {
		if client != sender {
			client.mu.Lock()
			err := client.conn.WriteJSON(message)
			client.mu.Unlock()

			if err != nil {
				log.Printf("Error broadcasting to client: %v \n", err)
			}
		}
	}
}

type ChatRoomManager struct {
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
}

func (crm *ChatRoomManager) addClient(chatRoom string, client *Client) {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	if crm.rooms[chatRoom] == nil {
		crm.rooms[chatRoom] = make(map[*Client]bool)
	}
	crm.rooms[chatRoom][client] = true
}

func (crm *ChatRoomManager) removeClient(chatRoom string, client *Client) {
	crm.mu.Lock()
	defer crm.mu.Unlock()

	if clients, exists := crm.rooms[chatRoom]; exists {
		delete(clients, client)
		if len(clients) == 0 {
			delete(crm.rooms, chatRoom)
		}
	}
}

func (crm *ChatRoomManager) broadcast(chatRoom string, message Message) {
	crm.mu.Lock()
	clients := crm.rooms[chatRoom]
	crm.mu.RUnlock()

	for client := range clients {
		client.mu.Lock()
		err := client.conn.WriteJSON(message)
		client.mu.Unlock()

		if err != nil {
			log.Printf("Error broadcasting to client: %v \n", err)
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsPosition)

	log.Println("WebSocket server starting on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func wsPosition(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to ws %v \n", err)
		return
	}

	userId := r.URL.Query().Get("userId")
	roomId := r.URL.Query().Get("roomId")

	if userId == "" || roomId == "" {
		log.Println("Missing clientId or roomId")
		conn.Close()
		return
	}
	client := &Client{
		conn:   conn,
		room:   roomId,
		userId: userId,
	}
	roomManager.addClient(roomId, client)
	defer func() {
		roomManager.removeClient(roomId, client)
		log.Println("client removed: ", client.userId)
		conn.Close()
	}()
	for {
		var message Message

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message: %v \n", err)
		}

		message.UserId = userId
		message.RoomId = roomId

		if message.Type == MessageTypePosition {
			roomManager.broadcast(roomId, message, client)
		}
	}
}

var UserPos = make(map[string]string) // key: value --> "x_y": "user1"

func updateUserPos(userId string, x int, y int) {
	strPos := strconv.Itoa(x) + "_" + strconv.Itoa(y) // "x_y"

	if user, ok := UserPos[strPos]; ok {
		// if it is, delete the user from the old position
		delete(UserPos, user)
	}
	UserPos[strPos] = userId
}

/*
	func handleWsPos () {
		maintain a userPosMap for each room.
		it will contain the positions of all users or rather the key here will be the coordinates and the value corresponding to that key will be userId

		that way we will know what all coordinates are invalid in that moment (if the key is present it means it is invalid)

		if a user wants to move, we will store it's old position in a variable
		then we will check if there is a key (new position) present in the map (is valid or not)
			if it is invalid (present) - don't do anything return (no websocket broadcast)
			if it is valid (not present) - delete the old key and add a new key with new pos and that userid

		for chat option
			on every move of a user, check the key corresponding to that user with all keys. If user's pos is in proximity to any other user,
				fire up a goroutine, wait for 1 sec -
					if still in proximity
					start a new chatroom between the 2 people

	}
*/
