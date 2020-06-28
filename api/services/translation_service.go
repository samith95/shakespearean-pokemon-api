package services

import (
	"errors"
	"net/http"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_domain"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_domain"
	"shakespearing-pokemon/api/domains/shksprean_pokemon_domain/shksprean_pokemon_error"
	"shakespearing-pokemon/api/domains/translation/translation_domain"
	"shakespearing-pokemon/api/providers/pokemon_provider"
	"shakespearing-pokemon/api/providers/translation_provider"
)

type translationService struct{}

type translationServiceInterface interface {
	GetShakespeareanPokemonTranslation(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface)
}

var (
	TranslationService translationServiceInterface = &translationService{}
)

func (t *translationService) GetShakespeareanPokemonTranslation(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (*shksprean_pokemon_domain.ShakespeareanPokemonResponse, shksprean_pokemon_error.ShkspreanPokemonErrorInterface) {
	request, err := validateRequestFields(request)
	if err != nil {
		return nil, shksprean_pokemon_error.New(http.StatusBadRequest, err.Error())
	}

	pokemonInfoReq := pokemon_domain.PokemonInfoRequest{Name: request.Name}

	//get description from pokemon provider
	pokemonInfoResp, pokemonErrorResp := pokemon_provider.PokemonProvider.GetPokemonInfo(pokemonInfoReq)
	if pokemonErrorResp != nil {
		return nil, shksprean_pokemon_error.New(pokemonErrorResp.Status(), pokemonErrorResp.Message())
	}

	translationRequest := translation_domain.TranslationRequest{}

	//get the most recent description form the pokemon info response
	for i := len(pokemonInfoResp.Description) - 1; i >= 0; i-- {
		if pokemonInfoResp.Description[i].Language.Name == "en" {
			translationRequest.Text = pokemonInfoResp.Description[i].Text
			break
		}
	}

	//get translation from Shakespearean translation provider
	translationResp, translationErrorResp := translation_provider.TranslationProvider.GetShakespeareanTranslation(translationRequest)
	if translationErrorResp != nil {
		return nil, shksprean_pokemon_error.New(translationErrorResp.Status(), translationErrorResp.Message())
	}

	response := &shksprean_pokemon_domain.ShakespeareanPokemonResponse{
		Name:        request.Name,
		Translation: translationResp.Content.Translation,
	}

	//generate the client response
	return response, nil
}

func validateRequestFields(request shksprean_pokemon_domain.ShakespeareanPokemonRequest) (shksprean_pokemon_domain.ShakespeareanPokemonRequest, error) {
	if request.Name == "" {
		return shksprean_pokemon_domain.ShakespeareanPokemonRequest{}, errors.New("name field cannot be empty")
	}
	return request, nil
}
