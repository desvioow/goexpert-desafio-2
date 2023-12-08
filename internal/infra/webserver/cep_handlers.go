package webserver

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/desvioow/goexpert-desafio-2/internal/dto"
	"github.com/go-chi/chi/v5"
)

func BrasilApiCepHandler(w http.ResponseWriter, r *http.Request) {

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		writeError(w, http.StatusBadRequest, "cep não informado")
		return
	}

	var url string = "https://brasilapi.com.br/api/cep/v1/" + cep
	resp, err := http.Get(url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao consumir API "+url)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao ler resposta da API "+url)
		return
	}

	var dto dto.BrasilApiOutput
	json.Unmarshal(body, &dto)

	processJSON(w, dto)
}

func ViaCepCepHandler(w http.ResponseWriter, r *http.Request) {

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		writeError(w, http.StatusBadRequest, "cep não informado")
	}

	var url string = "http://viacep.com.br/ws/" + cep + "/json/"
	resp, err := http.Get(url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao consumir API"+url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao ler resposta da API "+url)
		return
	}

	var dto dto.ViaCepOutput
	json.Unmarshal(body, &dto)

	processJSON(w, dto)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func processJSON(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error processing JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
