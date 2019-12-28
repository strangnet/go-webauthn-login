package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := createRouter()
	http.ListenAndServe(":4711", r)
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(Hello()))
	})
	return r
}

func Hello() string {
	return "Hello, World!"
}
