package webserver

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/desvioow/goexpert-desafio-2/internal/dto"
	"github.com/go-chi/chi/v5"
)

// Constantes para URLs e identicadfores das APIs externas
const (
	VIACEP_URL           = "http://viacep.com.br/ws/"
	VIACEP_IDENTIFIER    = "VIA_CEP"
	BRASILAPI_URL        = "https://brasilapi.com.br/api/cep/v1/"
	BRASILAPI_IDENTIFIER = "BRASIL_API"
)

// Regex para validar cep
var re = regexp.MustCompile("^[0-9]{8}$")

// Representação das respostas das APIs externas
type ExternalApiResponse struct {
	Response      *http.Response
	ApiIdentifier string
}

// Repesentação do output para o console
type ConsoleOutput struct {
	Api  string      `json:"api"`
	Data interface{} `json:"data"`
}

func FastestCepHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	cep := chi.URLParam(r, "cep")
	if !re.MatchString(cep) {
		writeError(w, http.StatusBadRequest, "cep inválido")
		return
	}

	viacepURI := VIACEP_URL + cep + "/json/"
	brasilapiURI := BRASILAPI_URL + cep
	responsesChannel := make(chan ExternalApiResponse)

	go externalCepApiRequest(ctx, brasilapiURI, BRASILAPI_IDENTIFIER, responsesChannel)
	go externalCepApiRequest(ctx, viacepURI, VIACEP_IDENTIFIER, responsesChannel)

	select {
	case <-ctx.Done():
		writeError(w, http.StatusRequestTimeout, "tempo de resposta excedido")
		return
	case resp := <-responsesChannel:
		defer resp.Response.Body.Close()

		body, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro ao ler resposta da API")
			return
		}

		processFastestCepApiResponse(w, body, resp.ApiIdentifier)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {

	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func externalCepApiRequest(ctx context.Context, url string, apiIdentifier string, ch chan<- ExternalApiResponse) {

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

	select {
	case ch <- ExternalApiResponse{Response: resp, ApiIdentifier: apiIdentifier}:
	case <-ctx.Done():
		return
	}
}

func processFastestCepApiResponse(w http.ResponseWriter, body []byte, apiIdentifier string) {

	type ApiDtoMap struct {
		Identifier string
		Dto        interface{}
	}

	var apiDtoMap = map[string]ApiDtoMap{
		BRASILAPI_IDENTIFIER: {Identifier: BRASILAPI_IDENTIFIER, Dto: dto.BrasilApiOutput{}},
		VIACEP_IDENTIFIER:    {Identifier: VIACEP_IDENTIFIER, Dto: dto.ViaCepOutput{}},
	}

	if apiDto, ok := apiDtoMap[apiIdentifier]; ok {
		err := json.Unmarshal(body, &apiDto.Dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Erro ao processar JSON para dados")
			return
		}
		processAndSendJSON(w, apiDto.Dto, apiDto.Identifier)
	}
}

func processAndSendJSON(w http.ResponseWriter, data interface{}, apiIdentifier string) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Erro ao processar dados para JSON")
		return
	}

	encoder := json.NewEncoder(os.Stdout)
	output := ConsoleOutput{
		Api:  apiIdentifier,
		Data: data,
	}

	err = encoder.Encode(output)
	if err != nil {
		log.Println("Erro ao encodar JSON para console:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
