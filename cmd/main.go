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
	dsn := "r91807sb_data:Gk15092006@tcp(r91807sb.beget.tech)/r91807sb_data?"
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

	adminRouter.Use(AdminAuthMiddleware)

	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/", handlers.Index).Methods("GET")
	mainRouter.HandleFunc("/publications", handlers.Publications).Methods("GET")
	mainRouter.HandleFunc("/ideas", handlers.Ideas).Methods("GET")
	mainRouter.HandleFunc("/smth", handlers.Smth).Methods("GET")
	mainRouter.HandleFunc("/profile", handlers.Persone).Methods("GET")
	mainRouter.PathPrefix("/admin").Handler(adminRouter)

	fileServer := http.FileServer(http.Dir("/web/static/"))

	mainRouter.Handle("/web/static/", http.StripPrefix("/web/static", fileServer))

	fmt.Println("starting server at :9990")
	http.ListenAndServe(":9990", mainRouter)
}
