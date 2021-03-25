package websocket

import (
	"log"
	"net/http"

	"github.com/niklod/highload-social-network/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	User *user.User
	Body MessageBody
}

type MessageBody struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	User *user.User
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) SendMessage(msg MessageBody) {
	err := c.Conn.WriteJSON(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := Message{
			User: c.User,
			Body: MessageBody{
				Type: messageType,
				Data: string(p),
			},
		}

		c.Pool.Messages <- message
	}
}

func upgrade(c *gin.Context) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}

	return ws, nil
}
