package shksprean_pokemon_error

import "encoding/json"

type ShkspreanPokemonErrorInterface interface {
	Status() int
	Message() string
}

type ShkspreanPokemonError struct {
	Error ErrorFields `json:"error"`
}

type ErrorFields struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ShkspreanPokemonError) Status() int {
	return e.Error.Code
}

func (e ShkspreanPokemonError) Message() string {
	return e.Error.Message
}

func New(statusCode int, message string) ShkspreanPokemonErrorInterface {
	return &ShkspreanPokemonError{
		ErrorFields{
			Code:    statusCode,
			Message: message,
		},
	}
}

func NewApiErrorFromBytes(body []byte) (ShkspreanPokemonErrorInterface, error) {
	var result ShkspreanPokemonError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
