package app

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ayzatziko/chat/repository"
)

func (chat *Chat) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	nickname := r.FormValue("nickname")
	password := r.FormValue("password")

	if username == "" || nickname == "" || password == "" {
		http.Error(w, "invalid form", 400)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = chat.Users.InsertUser(&repository.User{
		Username:     username,
		Nickname:     nickname,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, "/authorization", 301)
}

func (e *Chat) ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := e.Templates.ExecuteTemplate(w, "registration.html", nil); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
