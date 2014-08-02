package main

import (
	"fmt"
	"io"
	"net/http"
	"session"
	_ "session/Provider/memory"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "globalSessionId", 3600)

	go globalSessions.GC()
}

func login(w http.ResponseWriter, r *http.Request) {
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		str := "Welcome Back " + sess.SessionID()
		io.WriteString(w, str)
	} else {
		str := "Sadness " + sess.SessionID()
		io.WriteString(w, str)
	}
}

func main() {
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Listen failed")
	}
}
