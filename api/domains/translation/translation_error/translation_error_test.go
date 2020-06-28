package translation_error

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	expectedError := TranslationError{
		ErrorFields{
			Code:    400,
			Message: "Bad Request Error",
		},
	}
	actualError := New(expectedError.Error.Code, expectedError.Error.Message)
	assert.EqualValues(t, &expectedError, actualError)
}

func TestTranslationError(t *testing.T) {
	expectedError := TranslationError{
		ErrorFields{
			Code:    400,
			Message: "Bad Request Error",
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
