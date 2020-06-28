package shksprean_pokemon_domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShakespeareanPokemonResponse(t *testing.T) {
	expectedTranslation := ShakespeareanPokemonResponse{
		Name:        "charizard",
		Translation: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}

	bytes, err := json.Marshal(expectedTranslation)

	var actualResponse ShakespeareanPokemonResponse

	err = json.Unmarshal(bytes, &actualResponse)
	assert.Nil(t, err)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedTranslation.Name, actualResponse.Name)
	assert.EqualValues(t, expectedTranslation.Translation, actualResponse.Translation)
}
