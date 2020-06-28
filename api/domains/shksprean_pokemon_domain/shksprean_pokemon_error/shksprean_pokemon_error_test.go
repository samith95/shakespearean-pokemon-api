package shksprean_pokemon_error

import (
	"encoding/json"
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
	assert.EqualValues(t, expectedError.Message(), actualError.Message())
	assert.EqualValues(t, expectedError.Status(), actualError.Status())
}

func TestShkspreanPokemonError(t *testing.T) {
	expectedError := ShkspreanPokemonError{
		ErrorFields{
			Code:    400,
			Message: "error message",
		},
	}

	bytes, err := json.Marshal(expectedError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	actualError, err := NewApiErrorFromBytes(bytes)
	assert.Nil(t, err)

	assert.EqualValues(t, expectedError.Status(), actualError.Status())
	assert.EqualValues(t, expectedError.Message(), actualError.Message())
}
