package translation_controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_domain"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_error"
	"shakespearing-pokemon/api/services"
	"testing"
)

var (
	getShakespeareanPokemonTranslationFunc func(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface)
)

type translationServiceMock struct{}

func (t *translationServiceMock) GetShakespeareanPokemonTranslation(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface) {
	return getShakespeareanPokemonTranslationFunc(request)
}

func TestGetShakespeareanPokemonTranslationSuccess(t *testing.T) {
	expectedTranslation := shksprean_pokemon_domain.ShakespeareanPokemonResponse{
		Name:        "charizard",
		Translation: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}

	getShakespeareanPokemonTranslationFunc = func(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface) {
		return &expectedTranslation, nil
	}

	services.TranslationService = &translationServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{Key: "name", Value: "charizard"},
	}
	HandleShakespeareanPokemonTranslationRequest(c)
	var actualResponse shksprean_pokemon_domain.ShakespeareanPokemonResponse
	err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, expectedTranslation, actualResponse)
}

func TestGetShakespeareanPokemonTranslationInvalidName(t *testing.T) {
	expectedError := fmt.Sprintf("wrong name")

	getShakespeareanPokemonTranslationFunc = func(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface) {
		return nil, shksprean_pokemon_error.New(http.StatusInternalServerError, expectedError)
	}

	services.TranslationService = &translationServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{Key: "name", Value: "invalid name"},
	}
	HandleShakespeareanPokemonTranslationRequest(c)
	assert.EqualValues(t, http.StatusInternalServerError, response.Code)
	apiErr, err := shksprean_pokemon_error.NewApiErrorFromBytes(response.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, expectedError, apiErr.Message())
}

func TestGetShakespeareanPokemonTranslationSuccessSuccessIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	expectedTranslation := shksprean_pokemon_domain.ShakespeareanPokemonResponse{
		Name:        "charizard",
		Translation: "charizard flies 'round the sky in search of powerful opponents.'t breathes fire of such most wondrous heat yond 't melts aught. However,  't nev'r turns its fiery breath on any opponentweaker than itself.",
	}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{Key: "pokemonName", Value: "charizard"},
	}
	HandleShakespeareanPokemonTranslationRequest(c)
	var actualResponse shksprean_pokemon_domain.ShakespeareanPokemonResponse
	err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, expectedTranslation, actualResponse)
}
