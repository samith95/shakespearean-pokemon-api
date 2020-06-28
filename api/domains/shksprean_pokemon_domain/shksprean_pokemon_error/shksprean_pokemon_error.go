package shksprean_pokemon_error

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
