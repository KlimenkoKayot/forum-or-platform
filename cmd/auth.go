package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	checker map[string]bool = map[string]bool{
		"admin:0000": true,
	}
)

func checkAdminAuth(user string, pass string) bool {
	session_id := user + ":" + pass
	return checker[session_id]
}

func checkAdminSessionId(session_id string) bool {
	return checker[session_id]
}

func getAdminSession(user string, pass string) string {
	session_id := user + ":" + pass
	return session_id
}

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val, err := r.Cookie("admin_session_id")
		fmt.Println(val.Value)
		fmt.Println(checkAdminSessionId(val.Value))
		if err == http.ErrNoCookie || (err == nil && !checkAdminSessionId(val.Value)) {
			fmt.Println("["+r.URL.Path+"] "+"[adminAuthMiddleware] no auth... redirect /admin/login", r.URL.Path)
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("admin_session_id")
	if err == nil && checkAdminSessionId(session.Value) {
		fmt.Println(r.URL.Path, "have authorization... redirect /admin")
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	user, pass, ok := r.BasicAuth()
	if ok && checkAdminAuth(user, pass) {
		expiration := time.Now().Add(10 * time.Hour)
		cookie := http.Cookie{
			Name:    "admin_session_id",
			Value:   getAdminSession(user, pass),
			Expires: expiration,
		}
		http.SetCookie(w, &cookie)
		fmt.Println("admin logged in... redirect /admin")
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		fmt.Println("admin bad login...")
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (h *Handler) AdminExit(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("admin_session_id")
	if err == http.ErrNoCookie || !checkAdminSessionId(session.Value) {
		fmt.Println("admin bad login... redirect /admin/login")
		http.Redirect(w, r, "/admin/login", http.StatusFound)
		return
	}
	session.Value = ""
	session.Expires = time.Unix(0, 0)

	http.SetCookie(w, session)

	time.Sleep(3 * time.Second)
	_, err = r.Cookie("admin_session_id")
	if err == nil {
		fmt.Println("CANT DELETE COOKIE!!!")
		fmt.Printf("%#v\n", session)
	}
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}
