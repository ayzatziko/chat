package app

import (
	"net/http"
)

func (chat *Chat) ViewMessages(w http.ResponseWriter, r *http.Request) {
	messages := chat.Messages.Messages()
	chat.Templates.ExecuteTemplate(w, "messages.html", messages)
}
