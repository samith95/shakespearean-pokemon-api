package translation_error

import "encoding/json"

type TranslationErrorInterface interface {
	Status() int
	Message() string
}

type TranslationError struct {
	Error ErrorFields `json:"error"`
}

type ErrorFields struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e TranslationError) Status() int {
	return e.Error.Code
}

func (e TranslationError) Message() string {
	return e.Error.Message
}

func New(statusCode int, message string) TranslationErrorInterface {
	return &TranslationError{
		ErrorFields{
			Code:    statusCode,
			Message: message,
		},
	}
}

func NewApiErrorFromBytes(body []byte) (TranslationErrorInterface, error) {
	var result TranslationError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
