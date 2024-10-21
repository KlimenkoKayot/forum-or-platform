package main

import (
	"net/http"
	"time"
)

func (h *Handler) NewsAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	date := time.Now().Format("2006-January-02")
	author := r.FormValue("author")
	text := r.FormValue("text")
	h.DB.Exec("INSERT INTO news (`title`, `date`, `author`, `text`) VALUES (?, ?, ?, ?);", title, date, author, text)
	http.Redirect(w, r, "/news", http.StatusFound)
}

func (h *Handler) PublicationsAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	date := time.Now().Format("2006-January-02")
	author := r.FormValue("author")
	text := r.FormValue("text")
	h.DB.Exec("INSERT INTO publications (`title`, `date`, `author`, `text`) VALUES (?, ?, ?, ?);", title, date, author, text)
	http.Redirect(w, r, "/publications", http.StatusFound)
}

func (h *Handler) IdeasAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.FormValue("title")
	date := time.Now().Format("2006-January-02")
	text := r.FormValue("text")
	h.DB.Exec("INSERT INTO ideas (`title`, `date`, `text`) VALUES (?, ?, ?);", title, date, text)
	http.Redirect(w, r, "/ideas", http.StatusFound)
}
