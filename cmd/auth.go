package main

import (
	"fmt"
	"net/http"
)

var (
	checker map[string]bool = map[string]bool{
		"admin:admin": true,
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
		fmt.Println("adminAuthMiddleware", r.URL.Path)
		val, err := r.Cookie("admin_session_id")
		if err != nil || !checkAdminSessionId(val.Value) {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/admin/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	val, err := r.Cookie("admin_session_id")
	if err == nil && checkAdminSessionId(val.Value) {
		fmt.Println("habe admin auth", r.URL.Path)
		http.Redirect(w, r, "/admin", http.StatusFound)
		return
	}

	user, pass, ok := r.BasicAuth()
	if ok && checkAdminAuth(user, pass) {
		cookie := &http.Cookie{
			Name:  "admin_session_id",
			Value: getAdminSession(user, pass),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		fmt.Println("try admin login", r.URL.Path)
		w.Header().Set("WWW-Authenticate", `Basic readlm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
