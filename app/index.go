package app

import "net/http"

func (e *Chat) Index(w http.ResponseWriter, r *http.Request) {
	nickname := getNickname(r)
	if err := e.Templates.ExecuteTemplate(w, "index.html", nickname); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
