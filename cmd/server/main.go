package main

import (
	"net/http"

	"github.com/desvioow/goexpert-desafio-2/internal/infra/webserver"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/fastestcep/{cep}", webserver.FastestCepHandler)
	http.ListenAndServe(":8080", r)
}
