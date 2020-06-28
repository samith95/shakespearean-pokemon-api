package pokemon_domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPokemonInfoResponse(t *testing.T) {
	languageFields := LanguageFields{
		Name: "en",
	}

	flavourTextList := FlavourTextList{
		{Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Language: languageFields,
		},
		{Text: "Quisque cursus, metus vitae pharetra auctor, sem massa mattis sem, at interdum magna augue eget diam.",
			Language: languageFields,
		},
		{Text: "Morbi lectus risus, iaculis vel, suscipit quis, luctus non, massa.",
			Language: languageFields,
		},
		{Text: "Vestibulum lacinia arcu eget nulla.",
			Language: languageFields,
		},
	}

	expectedResponse := PokemonInfoResponse{
		Name:        "charizard",
		Description: flavourTextList,
	}

	bytes, err := json.Marshal(expectedResponse)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var actualResponse PokemonInfoResponse
	err = json.Unmarshal(bytes, &actualResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.Name, actualResponse.Name)
	for textIndex := range flavourTextList {
		assert.EqualValues(t, flavourTextList[textIndex].Text, actualResponse.Description[textIndex].Text)
		assert.EqualValues(t, flavourTextList[textIndex].Language.Name, actualResponse.Description[textIndex].Language.Name)
	}
}
