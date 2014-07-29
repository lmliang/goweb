package main

import (
	"session"
)

var globalSessions *session.Manager

func init() {
	globalSessions = session.NewManager("memory", "globalSessionId", 3600)

	go globalSessions.GC()
}
