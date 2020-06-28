package pokemon_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shakespearing-pokemon/api/clients/restclient"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_domain"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_error"
)

const (
	pokemonInfoUrl = "https://pokeapi.co/api/v2/pokemon-species/%s"
)

type pokemonProvider struct{}

type pokemonProviderInterface interface {
	GetPokemonInfo(request pokemon_domain.PokemonInfoRequest) (*pokemon_domain.PokemonInfoResponse, *pokemon_error.PokemonError)
}

var (
	//PokemonProvider used to mock the provider in test
	PokemonProvider pokemonProviderInterface = &pokemonProvider{}
)

func (p *pokemonProvider) GetPokemonInfo(request pokemon_domain.PokemonInfoRequest) (*pokemon_domain.PokemonInfoResponse, *pokemon_error.PokemonError) {
	url := fmt.Sprintf(pokemonInfoUrl, request.Name)
	bytes, errorResponse := getResults(url)
	if errorResponse != nil {
		return nil, errorResponse
	}

	var result pokemon_domain.PokemonInfoResponse
	err := json.Unmarshal(bytes, &result)
	errorResponse = createErrorResponse(err, "error when trying to unmarshal pokemon information response from API")
	if errorResponse != nil {
		return nil, errorResponse
	}

	return &result, nil
}

func getResults(url string) ([]byte, *pokemon_error.PokemonError) {
	response, err := restclient.ClientStruct.Get(url)
	errorResponse := createErrorResponse(err, "error when trying to get pokemon info results")
	if errorResponse != nil {
		return []byte{}, errorResponse
	}

	bytes, errorResponse := checkResponseBody(response)
	if errorResponse != nil {
		return []byte{}, errorResponse
	}
	return bytes, nil
}

func checkResponseBody(response *http.Response) ([]byte, *pokemon_error.PokemonError) {
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	errorResponse := createErrorResponse(err, "error when parsing the pokemon info response body")
	if errorResponse != nil {
		return []byte{}, errorResponse
	}

	if response.StatusCode > 299 {
		errorResponse := &pokemon_error.PokemonError{}
		errorResponse.Code = http.StatusInternalServerError
		errorResponse.ErrorMessage = "error from external api"
		if err := json.Unmarshal(bytes, &errorResponse); err != nil {
			errorResponse.ErrorMessage = fmt.Sprintf("error when unmarshaling request from external api: %s", err.Error())
		}
		return nil, errorResponse
	}

	return bytes, nil
}

func createErrorResponse(err error, errorMsg string) *pokemon_error.PokemonError {
	if err != nil {
		errorMsg := fmt.Sprintf(errorMsg+": %s", err.Error())
		log.Println(errorMsg)
		return &pokemon_error.PokemonError{
			Code:         http.StatusBadRequest,
			ErrorMessage: errorMsg,
		}
	}
	return nil
}
