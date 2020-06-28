package services

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_domain"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_error"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_domain"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_error"
	"shakespearing-pokemon/api/domains/translation/translation_domain"
	"shakespearing-pokemon/api/domains/translation/translation_error"
	"shakespearing-pokemon/api/providers/pokemon_provider"
	"shakespearing-pokemon/api/providers/translation_provider"
	"testing"
)

var (
	getPokemonInfo              func(request pokemon_domain.PokemonInfoRequest) (*pokemon_domain.PokemonInfoResponse, *pokemon_error.PokemonError)
	getShakespeareanTranslation func(request translation_domain.TranslationRequest) (*translation_domain.TranslationResponse, *translation_error.TranslationError)
)

type getPokemonProviderMock struct{}
type getTranslationProviderMock struct{}

func (p *getPokemonProviderMock) GetPokemonInfo(request pokemon_domain.PokemonInfoRequest) (*pokemon_domain.PokemonInfoResponse, *pokemon_error.PokemonError) {
	return getPokemonInfo(request)
}

func (s *getTranslationProviderMock) GetShakespeareanTranslation(request translation_domain.TranslationRequest) (*translation_domain.TranslationResponse, *translation_error.TranslationError) {
	return getShakespeareanTranslation(request)
}

func TestGetShakespeareanPokemonTranslationSuccess(t *testing.T) {
	languageField := pokemon_domain.LanguageFields{Name: "en"}

	descriptionField := pokemon_domain.FlavourTextList{
		{Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Language: languageField,
		},
		{Text: "Quisque cursus, metus vitae pharetra auctor, sem massa mattis sem, at interdum magna augue eget diam.",
			Language: languageField,
		},
		{Text: "Morbi lectus risus, iaculis vel, suscipit quis, luctus non, massa.",
			Language: languageField,
		},
		{Text: "Vestibulum lacinia arcu eget nulla.",
			Language: pokemon_domain.LanguageFields{Name: "jp"},
		},
	}

	mockPokemonInfoResp := pokemon_domain.PokemonInfoResponse{
		Name:        "charizard",
		Description: descriptionField,
	}

	mockTranslationResp := translation_domain.TranslationResponse{
		Content: translation_domain.ContentFields{
			Translation: "Vestibulum lacinia arcu eget nulla.",
		},
	}

	expectedResponse := shksprean_pokemon_domain.ShakespeareanPokemonResponse{
		Name:        "charizard",
		Translation: "Vestibulum lacinia arcu eget nulla.",
	}

	getPokemonInfo = func(request pokemon_domain.PokemonInfoRequest) (*pokemon_domain.PokemonInfoResponse, *pokemon_error.PokemonError) {
		return &mockPokemonInfoResp, nil
	}

	getShakespeareanTranslation = func(request translation_domain.TranslationRequest) (*translation_domain.TranslationResponse, *translation_error.TranslationError) {
		return &mockTranslationResp, nil
	}

	translation_provider.TranslationProvider = &getTranslationProviderMock{}
	pokemon_provider.PokemonProvider = &getPokemonProviderMock{}

	request := shksprean_pokemon_domain.ShakespeareanPokemonRequest{Name: "charizard"}
	actualResponse, err := TranslationService.GetShakespeareanPokemonTranslation(request)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse.Name, actualResponse.Name)
	assert.EqualValues(t, expectedResponse.Translation, actualResponse.Translation)
}

func TestGetShakespeareanPokemonTranslationSuccessIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	request := shksprean_pokemon_domain.ShakespeareanPokemonRequest{
		Name: "charizard",
	}

	actualResponse, err := TranslationService.GetShakespeareanPokemonTranslation(request)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	if actualResponse.Translation == "" {
		assert.Fail(t, "the translation field of the response is empty")
	}
	assert.EqualValues(t, request.Name, actualResponse.Name)
}

func TestGetShakespeareanPokemonTranslationWithEmptyName(t *testing.T) {
	expectedError := shksprean_pokemon_error.ShkspreanPokemonError{
		Error: shksprean_pokemon_error.ErrorFields{
			Code:    http.StatusBadRequest,
			Message: "name field cannot be empty",
		}}
	request := shksprean_pokemon_domain.ShakespeareanPokemonRequest{Name: ""}

	actualResponse, err := TranslationService.GetShakespeareanPokemonTranslation(request)
	assert.Nil(t, actualResponse)
	assert.NotNil(t, err)
	assert.EqualValues(t, expectedError.Status(), err.Status())
	assert.EqualValues(t, expectedError.Message(), err.Message())
}
