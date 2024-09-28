package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type tpl struct {
	Posts []*Post
	Post  *Post
}

type Post struct {
	Id      int
	Title   string
	Author  sql.NullString
	Text    string
	Updated sql.NullString
}

type Handler struct {
	DB        *sql.DB
	AdminTmpl *template.Template
	Tmpl      *template.Template
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	posts := []*Post{}

	rows, err := h.DB.Query("SELECT id, title, author, text, updated FROM posts")
	Check(err)
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Author, &post.Text, &post.Updated)
		Check(err)
		posts = append(posts, post)
	}
	rows.Close()

	err = h.AdminTmpl.ExecuteTemplate(w, "index.html", tpl{
		Posts: posts,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AdminAdd(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	author := r.FormValue("author")
	text := r.FormValue("text")

	if title == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `title` is requeired! </p>")
		return
	}
	if author == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `author` is requeired! </p>")
		return
	}
	if text == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `text` is requeired! </p>")
		return
	}

	result, err := h.DB.Exec("INSERT INTO posts (`title`, `author`, `text`) VALUES (?, ?, ?)", title, author, text)
	Check(err)
	newid, _ := result.LastInsertId()
	addcnt, _ := result.RowsAffected()
	fmt.Printf("Created!\n\tLast id: %v\n\tAdded rows: %v", newid, addcnt)
	fmt.Println(result.RowsAffected())

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handler) AdminAddPost(w http.ResponseWriter, r *http.Request) {
	err := h.AdminTmpl.ExecuteTemplate(w, "add.html", nil)
	Check(err)
}

func (h *Handler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TRY DELETE")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	Check(err)

	res, err := h.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	Check(err)

	affected, err := res.RowsAffected()
	Check(err)

	fmt.Println("DELETE by [" + r.UserAgent() + "]")
	fmt.Printf("\tid: %v\n\taffected: %v\n", id, affected)

	fmt.Println("DELETED ", id)
}

func (h *Handler) AdminEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	Check(err)

	rows, err := h.DB.Query("SELECT id, title, author, text, updated FROM posts WHERE id = ?", id)
	Check(err)
	post := &Post{}
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Author, &post.Text, &post.Updated)
		Check(err)
	}
	rows.Close()

	h.AdminTmpl.ExecuteTemplate(w, "edit.html", post)
}

func (h *Handler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	Check(err)

	title := r.FormValue("title")
	text := r.FormValue("text")
	updated := r.FormValue("updated")
	if title == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `title` is requeired! </p>")
		return
	}
	if updated == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `updated` is requeired! </p>")
		return
	}
	if text == "" {
		fmt.Println(r.UserAgent() + " badrequest")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "<p> param `text` is requeired! </p>")
		return
	}

	res, err := h.DB.Exec("UPDATE posts SET"+
		"`title` = ?,"+
		"`text` = ?,"+
		"`updated` = ?"+
		" WHERE id = ?",
		title,
		text,
		updated,
		id,
	)
	Check(err)

	affected, err := res.RowsAffected()
	Check(err)
	fmt.Println("UPDATED BY [" + r.UserAgent() + "]")
	fmt.Printf("\tid: %v\n\taffected: %v\n", id, affected)

	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "main.html", nil)
}

func (h *Handler) Publications(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "publications.html", nil)
}

func (h *Handler) Ideas(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "ideas.html", nil)
}

func (h *Handler) Smth(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "smth.html", nil)
}

func (h *Handler) Persone(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "persone.html", nil)
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("TEST ADMIN AUTH MIDDLEWARE")
				log.Fatal(err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
