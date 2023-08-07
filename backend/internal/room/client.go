package room

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	UserID   string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

// Sends messages from the room to the websocket connection.
func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

// Sends messages from the websocket connection to the room.
func (c *Client) readMessage(room *Room) {
	defer func() {
		room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(data),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		room.Broadcast <- msg
	}
}
