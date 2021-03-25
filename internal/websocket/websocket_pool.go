package websocket

import (
	"fmt"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
	Messages   chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		Messages:   make(chan Message),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client.User.Login] = client
			fmt.Printf("Добавлен клиент %s\n", client.User.Login)
			break

		case client := <-p.Unregister:
			delete(p.Clients, client.User.Login)
			fmt.Printf("Удален клиент %s\n", client.User.Login)
			break

		case message := <-p.Messages:
			fmt.Printf("Поступило сообщение от пользователя %s: %s\n", message.User.Login, message.Body.Data)
		}
	}
}
