package pokemon_provider

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"shakespearing-pokemon/api/clients/restclient"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_domain"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_error"
	"strings"
	"testing"
)

var getRequestFunc func(url string) (*http.Response, error)

type getClientMock struct{}

func (c *getClientMock) Get(request string) (*http.Response, error) {
	return getRequestFunc(request)
}

func TestGetPokemonInfo(t *testing.T) {
	languageFields := pokemon_domain.LanguageFields{
		Name: "en",
	}

	flavourTextList := pokemon_domain.FlavourTextList{
		{Text: "Spits fire that\nis hot enough to\nmelt boulders.\fKnown to cause\nforest fires\nunintentionally.",
			Language: languageFields,
		},
		{Text: "When expelling a\nblast of super\nhot fire, the red\fflame at the tip\nof its tail burns\nmore intensely.",
			Language: languageFields,
		},
		{Text: "If CHARIZARD be­\ncomes furious, the\nflame at the tip\fof its tail flares\nup in a whitish-\nblue color.",
			Language: languageFields,
		},
		{Text: "CHARIZARD flies around the sky in\nsearch of powerful opponents.\nIt breathes fire of such great heat\f" +
			"that it melts anything. However, it\nnever turns its fiery breath on any\nopponent weaker than itself.",
			Language: languageFields,
		},
	}

	expectedResponse := pokemon_domain.PokemonInfoResponse{
		Name:        "charizard",
		Description: flavourTextList,
	}

	b, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	r := bytes.NewReader(b)

	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(r),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(pokemon_domain.PokemonInfoRequest{Name: "doesn't matter"})
	assert.Nil(t, errorResponse)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse.Name, actualResponse.Name)
	for textIndex := range flavourTextList {
		assert.EqualValues(t, flavourTextList[textIndex].Text, actualResponse.Description[textIndex].Text)
		assert.EqualValues(t, flavourTextList[textIndex].Language.Name, actualResponse.Description[textIndex].Language.Name)
	}
}

func TestGetPokemonInfoInvalidErrorFormatting(t *testing.T) {
	expectedResponse := pokemon_error.PokemonError{
		Code:         http.StatusForbidden,
		ErrorMessage: "error occurred whilst un-marshaling error expectedResponse from api",
	}

	b, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	r := bytes.NewReader(b)

	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       ioutil.NopCloser(r),
		}, nil
	}
	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(pokemon_domain.PokemonInfoRequest{Name: "doesn't matter"})
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, expectedResponse.Code, errorResponse.Code)
	assert.EqualValues(t, expectedResponse.ErrorMessage, errorResponse.ErrorMessage)
}

func TestGetPokemonInfoInvalidBodyArguments(t *testing.T) {
	t.Parallel()
	getRequestFunc = func(url string) (*http.Response, error) {
		invalidCloser, _ := os.Open("-asf3")
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       invalidCloser,
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(pokemon_domain.PokemonInfoRequest{Name: "doesn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t, "error when parsing the pokemon info response body: invalid argument", errorResponse.ErrorMessage)
	assert.EqualValues(t, http.StatusBadRequest, errorResponse.Code)
}

//Checks in case the error response is invalid, it can happen when the external api changes their error response's data types
func TestGetPokemonInfosInvalidErrorInterface(t *testing.T) {
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"random'": 2020}`)),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(pokemon_domain.PokemonInfoRequest{Name: "doesn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t, "error from external api", errorResponse.ErrorMessage)
	assert.EqualValues(t, http.StatusInternalServerError, errorResponse.Code)
}

//Checks whether even though we are getting a positive response, the data types of the response might not be the same
//if in that case, the external api changed the data types, the unmarshaling will fail
func TestGetPokemonInfoInvalidResponseInterface(t *testing.T) {
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`{
						"flavor_text_entries": [
							{
								"flavor_text": "Spits fire that\nis hot enough to\nmelt boulders.\fKnown to cause\nforest fires\nunintentionally.",
								"language": {
									"name": THIS WILL NEVER PASS,
									"url": "https://pokeapi.co/api/v2/language/9/"
								},
							},
							{
								"flavor_text": "Spits fire that\nis hot enough to\nmelt boulders.\fKnown to cause\nforest fires\nunintentionally.",
								"language": {
									"name": "en",
									"url": "https://pokeapi.co/api/v2/language/9/"
								},
							}
						],
						"name": "charizard",
					}`)),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(pokemon_domain.PokemonInfoRequest{Name: "doesn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t,
		"error when trying to unmarshal pokemon information response from API: invalid character 'T' looking for beginning of value",
		errorResponse.ErrorMessage)
	assert.EqualValues(t, http.StatusBadRequest, errorResponse.Code)
}

func TestGetPokemonInfoIntegration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	request := pokemon_domain.PokemonInfoRequest{
		Name: "charizard",
	}

	actualResponse, errorResponse := PokemonProvider.GetPokemonInfo(request)
	assert.Nil(t, errorResponse)
	if len(actualResponse.Description) == 0 {
		assert.Fail(t, "pokemon info from API is empty")
	}
	assert.EqualValues(t, request.Name, actualResponse.Name)
}
