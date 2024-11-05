package main

import (
	"log"
	"net/http"
)

func WsChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to ws %v \n", err)
		return
	}

	userId := r.URL.Query().Get("userId")
	roomId := r.URL.Query().Get("roomId")

	if userId == "" || roomId == "" {
		log.Println("Missing clientId, and/or roomId")
		conn.Close()
		return
	}
	client := &Client{
		conn:   conn,
		room:   roomId,
		userId: userId,
	}
	chatRoomManager.addClient(roomId, client)
	defer func() {
		chatRoomManager.removeClient(roomId, client)
		conn.Close()
	}()

	for {
		var message Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading the message: %v \n", err)
		}

		message.UserId = userId
		message.RoomId = roomId

		if message.Type == MessageTypeChat {
			chatRoomManager.broadcast(roomId, message)
		} else {
			log.Printf("Unknown message type: %s \n", message.Type)
		}
	}

}
