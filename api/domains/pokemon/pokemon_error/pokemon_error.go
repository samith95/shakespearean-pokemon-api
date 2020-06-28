package pokemon_error

import (
	"encoding/json"
)

type PokemonErrorInterface interface {
	Status() int
	Message() string
}

type PokemonError struct {
	Code         int
	ErrorMessage string `json:",omitempty"`
}

func (p *PokemonError) Status() int {
	return p.Code
}

func (p *PokemonError) Message() string {
	return p.ErrorMessage
}

func New(statusCode int, errorMessage string) PokemonErrorInterface {
	return &PokemonError{
		Code:         statusCode,
		ErrorMessage: errorMessage,
	}
}

func NewApiErrorFromBytes(body []byte) (PokemonErrorInterface, error) {
	var result PokemonError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
