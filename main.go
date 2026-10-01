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
		if short == "" {
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}
		long, exists := svc.Resolve(short)
		if !exists {
			http.Error(w, fmt.Sprintf("Short URL '%s' not found", short), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, long, http.StatusMovedPermanently)
	})

	http.HandleFunc("/stats/", func(w http.ResponseWriter, r *http.Request) {
		short := r.URL.Path[len("/stats/"):]
		if short == "" {
			http.Redirect(w, r, "/stats", http.StatusMovedPermanently)
			return
		}
		hits := svc.GetHits(short)
		fmt.Fprintf(w, "URL %s has been visited %d times", short, hits)
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := svc.ListAll()
		fmt.Fprintf(w, "All Shortened URLs:\n")
		for _, s := range stats {
			fmt.Fprintf(w, "%s -> %s (%d hits)\n", s.ShortURL, s.LongURL, s.Hits)
		}
	})

	fmt.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}