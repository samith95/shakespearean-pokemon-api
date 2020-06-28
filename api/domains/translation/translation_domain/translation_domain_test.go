package translation_domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslationResponse(t *testing.T) {
	content := ContentFields{
		Translation: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}
	expectedResponse := TranslationResponse{
		Content: content,
	}

	bytes, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var actualResponse TranslationResponse

	err = json.Unmarshal(bytes, &actualResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, expectedResponse.Content.Translation, actualResponse.Content.Translation)
}
