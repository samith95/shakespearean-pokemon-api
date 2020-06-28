package shksprean_pokemon_error

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	expectedError := ShkspreanPokemonError{
		ErrorFields{
			Code:    400,
			Message: "Bad Request Error",
		},
	}
	actualError := New(expectedError.Error.Code, expectedError.Error.Message)
	assert.EqualValues(t, &expectedError, actualError)
}
