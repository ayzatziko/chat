package app

import (
	"github.com/gorilla/mux"
)

func Routing(chat *Chat) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/registration", chat.ShowRegisterPage).Methods("GET")
	r.HandleFunc("/registration", chat.Register).Methods("POST")

	r.HandleFunc("/authorization", chat.Authorization).Methods("GET")
	r.HandleFunc("/authorization", chat.Authorize).Methods("POST")

	r.HandleFunc("/logout", chat.Logout).Methods("GET")

	r.HandleFunc("/", chat.Auth(chat.Index)).Methods("GET")

	r.HandleFunc("/profile", chat.Auth(chat.EditNickname)).Methods("GET")
	r.HandleFunc("/profile", chat.Auth(chat.ChangeNickname)).Methods("POST")

	r.HandleFunc("/ws", chat.Auth(chat.HandleConnections))
	r.HandleFunc("/messages", chat.ViewMessages).Methods("GET")

	return r
}
