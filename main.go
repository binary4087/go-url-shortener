package main

import (
	"fmt"
	"net/http"
	"github.com/binary4087/go-url-shortener/shortener"
)

func main() {
	svc := shortener.New()

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "Missing url parameter", http.StatusBadRequest)
			return
		}
		short := svc.Shorten(url)
		fmt.Fprintf(w, "Short URL: /r/%s", short)
	})

	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		short := r.URL.Path[len("/r/"):]
		long, exists := svc.Resolve(short)
		if !exists {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, long, http.StatusMovedPermanently)
	})

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}