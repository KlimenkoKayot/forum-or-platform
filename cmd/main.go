package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Connecting to database...")
	dsn := "r91807sb_data:" + /* тут должен быть пароль */ "@tcp(r91807sb.beget.tech)/r91807sb_data?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	Check(err)

	fmt.Println("Ping!")
	err = db.Ping()
	Check(err)
	fmt.Println("     Pong!")

	handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("../web/html/*")),
	}

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/admin", handlers.AdminIndex).Methods("GET")
	adminRouter.HandleFunc("/admin/add", handlers.AdminAddPost).Methods("GET")
	adminRouter.HandleFunc("/admin/add", handlers.AdminAdd).Methods("POST")
	adminRouter.HandleFunc("/admin/edit/{id}", handlers.AdminEdit).Methods("GET")
	adminRouter.HandleFunc("/admin/edit/{id}", handlers.AdminUpdate).Methods("POST")
	adminRouter.HandleFunc("/admin/delete/{id}", handlers.AdminDelete).Methods("DELETE")
	adminRouter.Use(adminAuthMiddleware)

	actionRouter := mux.NewRouter()
	actionRouter.HandleFunc("/action/news/add", handlers.NewsAdd).Methods("POST")

	mainRouter := mux.NewRouter()
	// admin login
	mainRouter.HandleFunc("/admin/login", handlers.AdminLogin).Methods("GET")
	mainRouter.HandleFunc("/admin/exit", handlers.AdminExit).Methods("GET")

	staticDir := "/web/static/"
	mainRouter.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(".."+staticDir))))

	mainRouter.HandleFunc("/", handlers.Index).Methods("GET")
	mainRouter.HandleFunc("/publications", handlers.Publications).Methods("GET")
	mainRouter.HandleFunc("/ideas", handlers.Ideas).Methods("GET")
	mainRouter.HandleFunc("/news", handlers.News).Methods("GET")
	mainRouter.HandleFunc("/profile", handlers.Persone).Methods("GET")
	mainRouter.PathPrefix("/admin").Handler(adminRouter)
	mainRouter.PathPrefix("/action").Handler(actionRouter)

	fmt.Println("starting server at :9990")
	http.ListenAndServe(":9990", mainRouter)
}
