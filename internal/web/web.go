package web

import (
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// Define headers para HTML
	headers(w)

	http.ServeFile(w, r, "./pages/web.html")
}

func headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
