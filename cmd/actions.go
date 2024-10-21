package main

import (
	"net/http"
	"time"
)

func (h *Handler) NewsAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	date := time.Now().String()
	author := r.FormValue("author")
	text := r.FormValue("text")
	h.DB.Exec("INSERT INTO news (`title`, `date`, `author`, `text`) VALUES (?, ?, ?, ?);", title, date, author, text)
	http.Redirect(w, r, "/news", http.StatusFound)
}
