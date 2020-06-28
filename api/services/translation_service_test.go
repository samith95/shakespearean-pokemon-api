package services

import (
	"shakespearing-pokemon/api/domains/pokemon/pokemon_domain"
	"shakespearing-pokemon/api/domains/pokemon/pokemon_error"
	"shakespearing-pokemon/api/domains/translation/translation_domain"
	"shakespearing-pokemon/api/domains/translation/translation_error"
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
