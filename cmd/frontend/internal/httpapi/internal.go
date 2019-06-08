package httpapi

import "net/http"

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
