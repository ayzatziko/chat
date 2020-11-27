package app

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (chat *Chat) Authorize(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "invalid login/password", 400)
		return
	}

	u, err := chat.Users.UserByUsername(username)
	if err != nil {
		http.Redirect(w, r, "/registration", 301)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "nickname", Value: u.Nickname})
	http.Redirect(w, r, "/", 301)
}

func (e *Chat) Authorization(w http.ResponseWriter, r *http.Request) {
	if err := e.Templates.ExecuteTemplate(w, "authorization.html", nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (chat *Chat) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("nickname")
		if err != nil {
			logout(w, r)
			http.Redirect(w, r, "/authorization", 301)
			return
		}

		exists, err := chat.Users.NicknameExists(c.Value)
		if err != nil {
			logout(w, r)
			http.Redirect(w, r, "/authorization", 301)
			return
		}
		if !exists {
			logout(w, r)
			http.Redirect(w, r, "/authorization", 301)
			return
		}

		next(w, r)
	}
}

func (*Chat) Logout(w http.ResponseWriter, r *http.Request) {
	logout(w, r)
}

func getNickname(r *http.Request) string {
	c, _ := r.Cookie("nickname")
	return c.Value
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "nickname", Value: "", MaxAge: -1})
}
