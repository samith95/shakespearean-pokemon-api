package translation_provider

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"shakespearing-pokemon/api/clients/restclient"
	"shakespearing-pokemon/api/domains/translation/translation_domain"
	"shakespearing-pokemon/api/domains/translation/translation_error"
	"strings"
	"testing"
)

var getRequestFunc func(url string) (*http.Response, error)

type getClientMock struct{}

func (c *getClientMock) Get(request string) (*http.Response, error) {
	return getRequestFunc(request)
}

func TestGetShakespeareanTranslation(t *testing.T) {
	ContentFields := translation_domain.ContentFields{
		Translation: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	}

	expectedResponse := translation_domain.TranslationResponse{
		Content: ContentFields,
	}

	b, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	r := bytes.NewReader(b)

	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(r),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(translation_domain.TranslationRequest{Text: "deosn't matter"})
	assert.Nil(t, errorResponse)
	assert.NotNil(t, actualResponse)
	assert.EqualValues(t, expectedResponse.Content.Translation, actualResponse.Content.Translation)
}

func TestGetShakespeareanTranslationInvalidErrorFormatting(t *testing.T) {
	expectedResponse := translation_error.TranslationError{
		Error: translation_error.ErrorFields{
			Code:    http.StatusForbidden,
			Message: "error occurred whilst un-marshaling error expectedResponse from api",
		},
	}

	b, err := json.Marshal(expectedResponse)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	r := bytes.NewReader(b)

	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       ioutil.NopCloser(r),
		}, nil
	}
	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(translation_domain.TranslationRequest{Text: "deosn't matter"})
	assert.Nil(t, actualResponse)
	assert.EqualValues(t, expectedResponse.Error.Code, errorResponse.Error.Code)
	assert.EqualValues(t, expectedResponse.Error.Message, errorResponse.Error.Message)
}

func TestGetShakespeareanTranslationInvalidBodyArguments(t *testing.T) {
	t.Parallel()
	getRequestFunc = func(url string) (*http.Response, error) {
		invalidCloser, _ := os.Open("-asf3")
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       invalidCloser,
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(translation_domain.TranslationRequest{Text: "deosn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t, "error when parsing the translation response body: invalid argument", errorResponse.Error.Message)
	assert.EqualValues(t, http.StatusBadRequest, errorResponse.Error.Code)
}

//Checks in case the error response is invalid, it can happen when the external api changes their error response's data types
func TestGetShakespeareanTranslationInvalidErrorInterface(t *testing.T) {
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"random'": 2020}`)),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(translation_domain.TranslationRequest{Text: "deosn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t, "error from external api", errorResponse.Error.Message)
	assert.EqualValues(t, http.StatusInternalServerError, errorResponse.Error.Code)
}

//Checks whether even though we are getting a positive response, the data types of the response might not be the same
//if in that case, the external api changed the data types, the unmarshaling will fail
func TestGetShakespeareanTranslationInvalidResponseInterface(t *testing.T) {
	getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`{
										"success": {
											"total": 1
										},
										"contents": {
											"translated": ["wrong","type",],
											"text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
											"translation": "shakespeare"
										}
									}`)),
		}, nil
	}

	restclient.ClientStruct = &getClientMock{}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(translation_domain.TranslationRequest{Text: "deosn't matter"})
	assert.Nil(t, actualResponse)
	assert.NotNil(t, errorResponse)
	assert.EqualValues(t,
		"error when trying to unmarshal translation response from API: invalid character ']' looking for beginning of value",
		errorResponse.Error.Message)
	assert.EqualValues(t, http.StatusBadRequest, errorResponse.Error.Code)
}

func TestGetShakespeareanTranslationIntegration(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	request := translation_domain.TranslationRequest{
		Text: "CHARIZARD flies around the sky in\nsearch of powerful opponents.\nIt breathes fire of such great heat\f" +
			"that it melts anything. However, it\nnever turns its fiery breath on any\nopponent weaker than itself.",
	}

	actualResponse, errorResponse := TranslationProvider.GetShakespeareanTranslation(request)
	assert.Nil(t, errorResponse)
	if actualResponse.Content.Translation == "" {
		assert.Fail(t, "translation from API is empty")
	}
}
