package client

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	// channel that olds the data to be sent
	send     chan *message
	room     *Room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			avatarURL, err := c.room.avatar.GetAvatarURL(c)
			if err != nil {
				log.Println("No avatar url found")
			}
			msg.AvatarURL = avatarURL
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
