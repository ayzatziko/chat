package app

import (
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/ayzatziko/chat/repository"
)

type Chat struct {
	Port      string
	Users     repository.Users
	Messages  repository.Messages
	Broadcast chan *repository.Message

	cliMu   sync.Mutex
	Clients map[string]*websocket.Conn

	Templates *template.Template
	Upgrader  *websocket.Upgrader
	Router    *mux.Router
}

func CreateChat(users repository.Users, messages repository.Messages, port string) *Chat {
	tplPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ayzatziko", "chat", "templates")
	chat := Chat{
		Port:      port,
		Users:     users,
		Messages:  messages,
		Clients:   make(map[string]*websocket.Conn),
		Broadcast: make(chan *repository.Message),
		Templates: template.Must(template.ParseFiles(
			filepath.Join(tplPath, "index.html"),
			filepath.Join(tplPath, "registration.html"),
			filepath.Join(tplPath, "authorization.html"),
			filepath.Join(tplPath, "profile.html"),
			filepath.Join(tplPath, "messages.html"),
		)),
		Upgrader: &websocket.Upgrader{},
	}

	chat.Router = Routing(&chat)

	go chat.HandleMessages()

	return &chat
}
