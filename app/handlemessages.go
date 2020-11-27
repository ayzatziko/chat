package app

import (
	"fmt"

	"github.com/ayzatziko/chat/repository"
)

func (chat *Chat) HandleMessages() {
	for {
		message := <-chat.Broadcast
		if message.Recipient == "" {
			chat.Messages.Add(message)

			for nick := range chat.Clients {
				withLock(&chat.cliMu, func() {
					chat.writeLocked(nick, message)
				})
			}
		} else {
			withLock(&chat.cliMu, func() {
				if _, ok := chat.Clients[message.Recipient]; !ok {
					chat.writeLocked(message.Author, &repository.Message{Author: "System", Body: fmt.Sprintf("%s is offline", message.Recipient)})
				} else {
					chat.writeLocked(message.Author, message)
					chat.writeLocked(message.Recipient, message)
				}
			})
		}
	}
}

func (chat *Chat) writeLocked(nick string, m *repository.Message) {
	cli, ok := chat.Clients[nick]
	if !ok {
		return
	}

	if err := cli.WriteJSON(m); err != nil {
		cli.Close()
		delete(chat.Clients, nick)
	}
}
