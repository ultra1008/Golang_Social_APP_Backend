package websocket

import (
	"log"

	"github.com/niklod/highload-social-network/internal/user"

	"github.com/gin-gonic/gin"
)

type WebsocketHandler struct {
	pool        *Pool
	userService *user.Service
}

func NewWebsocketHandler(pool *Pool, userService *user.Service) *WebsocketHandler {
	return &WebsocketHandler{
		pool:        pool,
		userService: userService,
	}
}

func (w *WebsocketHandler) HandleWS(c *gin.Context) {
	login := c.Param("login")

	user, err := w.userService.GetUserByLogin(login)
	if err != nil {
		log.Printf("getting user for ws connection: %v\n", err)
		return
	}

	ws, err := upgrade(c)
	if err != nil {
		log.Println(err)
		return
	}

	if user == nil {
		return
	}

	client := &Client{
		User: user,
		Conn: ws,
		Pool: w.pool,
	}

	w.pool.Register <- client
	client.Read()
}
