package routes

import "net/http"

// Index is for /
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("don't panic"))
}
