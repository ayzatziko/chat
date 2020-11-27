package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ayzatziko/chat/repository"
)

func (chat *Chat) HandleConnections(w http.ResponseWriter, r *http.Request) {
	nickname := getNickname(r)
	conn, err := chat.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	for _, msg := range chat.Messages.Messages() {
		if err := conn.WriteJSON(msg); err != nil {
			return
		}
	}

	var message repository.Message

	withLock(&chat.cliMu, func() {
		chat.Clients[nickname] = conn
	})

	for {
		if err = conn.ReadJSON(&message); err != nil {
			message = repository.Message{Author: "System", Body: fmt.Sprintf("%s disconnected due to bad connection", nickname)}
			chat.Broadcast <- &message

			withLock(&chat.cliMu, func() {
				delete(chat.Clients, nickname)
			})
			break
		}
		chat.Broadcast <- &message
	}

	http.Error(w, "connection closed", 400)
}

func withLock(l sync.Locker, f func()) {
	l.Lock()
	f()
	l.Unlock()
}
