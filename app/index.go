package app

import "net/http"

func (e *Chat) Index(w http.ResponseWriter, r *http.Request) {
	v := IndexView{Port: e.Port, Nickname: getNickname(r)}
	if err := e.Templates.ExecuteTemplate(w, "index.html", v); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type IndexView struct {
	Port     string
	Nickname string
}
