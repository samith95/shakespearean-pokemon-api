package pokemon_error

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	expectedError := PokemonError{
		Code:         400,
		ErrorMessage: "bad request",
	}
	actualError := New(expectedError.Code, expectedError.ErrorMessage)
	assert.EqualValues(t, &expectedError, actualError)
}

func TestPokemonError(t *testing.T) {
	expectedError := PokemonError{
		Code: 400,
	}

	bytes, err := json.Marshal(expectedError)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	actualError, err := NewApiErrorFromBytes(bytes)
	assert.Nil(t, err)

	assert.EqualValues(t, expectedError.Status(), actualError.Status())
	assert.EqualValues(t, expectedError.Message(), actualError.Message())

}
