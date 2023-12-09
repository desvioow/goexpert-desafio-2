package main

import (
	"net/http"

	"github.com/desvioow/goexpert-desafio-2/internal/infra/webserver"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Get("/fastestcep/{cep}", webserver.FastestCepHandler)
	http.ListenAndServe(":8080", r)
}

/*

 Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:

https://brasilapi.com.br/api/cep/v1/01153000 + cep

http://viacep.com.br/ws/" + cep + "/json/

Os requisitos para este desafio são:

- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.

- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
*/
