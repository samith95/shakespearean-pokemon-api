package translation_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shakespearing-pokemon/api/clients/restclient"
	"shakespearing-pokemon/api/domains/translation/translation_domain"
	"shakespearing-pokemon/api/domains/translation/translation_error"
)

const (
	shakeSpeareTranslationUrl = "https://api.funtranslations.com/translate/shakespeare.json?text=\"%s\""
)

type translationProvider struct{}

type translationProviderInterface interface {
	GetShakespeareanTranslation(request translation_domain.TranslationRequest) (*translation_domain.TranslationResponse, *translation_error.TranslationError)
}

var (
	//TranslationProvider is used to mock teh provider in test
	TranslationProvider translationProviderInterface = &translationProvider{}
)

func (t *translationProvider) GetShakespeareanTranslation(request translation_domain.TranslationRequest) (*translation_domain.TranslationResponse, *translation_error.TranslationError) {
	url := fmt.Sprintf(shakeSpeareTranslationUrl, request.Text)
	bytes, errorResponse := getResults(url)
	if errorResponse != nil {
		return nil, errorResponse
	}

	var result translation_domain.TranslationResponse
	err := json.Unmarshal(bytes, &result)
	errorResponse = createErrorResponse(err, "error when trying to unmarshal translation response from API")
	if errorResponse != nil {
		return nil, errorResponse
	}

	return &result, nil
}

func getResults(url string) ([]byte, *translation_error.TranslationError) {
	response, err := restclient.ClientStruct.Get(url)
	errorResponse := createErrorResponse(err, "error when trying to get translation results")
	if errorResponse != nil {
		return []byte{}, errorResponse
	}

	bytes, errorResponse := checkResponseBody(response)
	if errorResponse != nil {
		return []byte{}, errorResponse
	}
	return bytes, nil
}

func checkResponseBody(response *http.Response) ([]byte, *translation_error.TranslationError) {
	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	errorResponse := createErrorResponse(err, "error when parsing the translation response body")
	if errorResponse != nil {
		return []byte{}, errorResponse
	}

	if response.StatusCode > 299 {
		errorResponse := &translation_error.TranslationError{}
		errorResponse.Error.Code = http.StatusInternalServerError
		errorResponse.Error.Message = "error from external api"
		if err := json.Unmarshal(bytes, &errorResponse); err != nil {
			errorResponse.Error.Message = fmt.Sprintf("error when unmarshaling request from external api: %s", err.Error())
		}
		return nil, errorResponse
	}

	return bytes, nil
}

func createErrorResponse(err error, errorMsg string) *translation_error.TranslationError {
	if err != nil {
		errorMsg := fmt.Sprintf(errorMsg+": %s", err.Error())
		log.Println(errorMsg)
		return &translation_error.TranslationError{
			Error: translation_error.ErrorFields{
				Code:    http.StatusBadRequest,
				Message: errorMsg,
			},
		}
	}
	return nil
}
