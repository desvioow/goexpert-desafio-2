package webserver

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/desvioow/goexpert-desafio-2/internal/dto"
	"github.com/go-chi/chi/v5"
)

const (
	VIACEP_URL           = "http://viacep.com.br/ws/"
	VIACEP_IDENTIFIER    = "VIA_CEP"
	BRASILAPI_URL        = "https://brasilapi.com.br/api/cep/v1/"
	BRASILAPI_IDENTIFIER = "BRASIL_API"
)

type apiResponse struct {
	Response      *http.Response
	ApiIdentifier string
}

func FastestCepHandler(w http.ResponseWriter, r *http.Request) {

	cep := chi.URLParam(r, "cep")
	if cep == "" {
		writeError(w, http.StatusBadRequest, "cep não informado")
		return
	}

	viacepURI := VIACEP_URL + cep + "/json/"
	brasilapiURI := BRASILAPI_URL + cep

	responsesChannel := make(chan apiResponse, 2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go externalCepApiRequest(ctx, brasilapiURI, BRASILAPI_IDENTIFIER, responsesChannel)
	go externalCepApiRequest(ctx, viacepURI, VIACEP_IDENTIFIER, responsesChannel)

	resp := <-responsesChannel
	defer resp.Response.Body.Close()

	body, err := io.ReadAll(resp.Response.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao ler resposta da API")
		return
	}

	processFastestCepApiResponse(w, body, resp.ApiIdentifier)

}

func externalCepApiRequest(ctx context.Context, url string, apiIdentifier string, ch chan<- apiResponse) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	ch <- apiResponse{Response: resp, ApiIdentifier: apiIdentifier}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {

	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func processFastestCepApiResponse(w http.ResponseWriter, body []byte, apiIdentifier string) {

	if apiIdentifier == BRASILAPI_IDENTIFIER {
		var dto dto.BrasilApiOutput
		json.Unmarshal(body, &dto)
		processJSON(w, dto)
	} else if apiIdentifier == VIACEP_IDENTIFIER {
		var dto dto.ViaCepOutput
		json.Unmarshal(body, &dto)
		processJSON(w, dto)
	}
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
