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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) AdminIndex(w http.ResponseWriter, r *http.Request) {
	posts := []*Post{}

	rows, err := h.DB.Query("SELECT id, title, author, text, updated FROM posts")
	check(err)
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Author, &post.Text, &post.Updated)
		check(err)
		posts = append(posts, post)
	}
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "index.html", tpl{
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
	check(err)
	newid, _ := result.LastInsertId()
	addcnt, _ := result.RowsAffected()
	fmt.Printf("Created!\n\tLast id: %v\n\tAdded rows: %v", newid, addcnt)
	fmt.Println(result.RowsAffected())

	http.Redirect(w, r, "/posts", http.StatusFound)
}

func (h *Handler) AdminAddPost(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "add.html", nil)
	check(err)
}

func (h *Handler) AdminDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TRY DELETE")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	check(err)

	res, err := h.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	check(err)

	affected, err := res.RowsAffected()
	check(err)

	fmt.Println("DELETE by [" + r.UserAgent() + "]")
	fmt.Printf("\tid: %v\n\taffected: %v\n", id, affected)

	fmt.Println("DELETED ", id)
}

func (h *Handler) AdminEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	check(err)

	rows, err := h.DB.Query("SELECT id, title, author, text, updated FROM posts WHERE id = ?", id)
	check(err)
	post := &Post{}
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Author, &post.Text, &post.Updated)
		check(err)
	}
	rows.Close()

	h.Tmpl.ExecuteTemplate(w, "edit.html", post)
}

func (h *Handler) AdminUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	check(err)

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
	check(err)

	affected, err := res.RowsAffected()
	check(err)
	fmt.Println("UPDATED BY [" + r.UserAgent() + "]")
	fmt.Printf("\tid: %v\n\taffected: %v\n", id, affected)

	http.Redirect(w, r, "/posts", http.StatusFound)
}

func adminAuthMiddleware(next http.Handler) http.Handler {
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

func main() {
	fmt.Println("Connecting to database...")
	dsn := "r91807sb_data:Gk15092006@tcp(r91807sb.beget.tech)/r91807sb_data?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	check(err)

	fmt.Println("Ping!")
	err = db.Ping()
	check(err)
	fmt.Println("     Pong!")

	handlers := &Handler{
		DB:        db,
		AdminTmpl: template.Must(template.ParseGlob("templates/*")),
		Tmpl:      template.Must(template.ParseGlob("html/*")),
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin", handlers.AdminIndex).Methods("GET")
	adminRouter.HandleFunc("/admin/add", handlers.AdminAddPost).Methods("GET")
	adminRouter.HandleFunc("/admin/add", handlers.AdminAdd).Methods("POST")
	adminRouter.HandleFunc("/admin/edit/{id}", handlers.AdminEdit).Methods("GET")
	adminRouter.HandleFunc("/admin/edit/{id}", handlers.AdminUpdate).Methods("POST")
	adminRouter.HandleFunc("/admin/delete/{id}", handlers.AdminDelete).Methods("DELETE")

	adminHandler := adminAuthMiddleware(adminRouter)

	mainRouter := mux.NewRouter()
	mainRouter.Handle("/admin", adminHandler)
	mainRouter.HandleFunc("/", handlers.Index).Methods("GET")
	mainRouter.HandleFunc("/publications", handlers.Publications).Methods("GET")
	mainRouter.HandleFunc("/ideas", handlers.Ideas).Methods("GET")
	mainRouter.HandleFunc("/smth", handlers.Smth).Methods("GET")
	mainRouter.HandleFunc("/user", handlers.Persone).Methods("GET")

	fmt.Println("starting server at :9990")
	http.ListenAndServe(":9990", mainRouter)
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
