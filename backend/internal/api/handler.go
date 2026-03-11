package api

import "net/http"

func HandlerRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Go temperature broker 🌡️  from sunny side TUS 🏜️\n"))
}

func HandlerHealthChecker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Application health is available\n"))
}
