package app

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (chat *Chat) ChangeNickname(w http.ResponseWriter, r *http.Request) {
	newNickname := r.FormValue("new_nickname")
	password := r.FormValue("password")
	if newNickname == "" {
		http.Error(w, "bad nickname", 400)
		return
	}

	u, err := chat.Users.UserByNickname(getNickname(r))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	u.Nickname = newNickname
	if err = chat.Users.UpdateUser(u); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "nickname", Value: newNickname})
	w.WriteHeader(200)
}

func (chat *Chat) EditNickname(w http.ResponseWriter, r *http.Request) {
	if err := chat.Templates.ExecuteTemplate(w, "profile.html", getNickname(r)); err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(200)
}
