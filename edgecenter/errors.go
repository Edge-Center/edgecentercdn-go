package edgecenter

import (
	"encoding/json"
	"errors"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
	ErrRateLimit    = errors.New("rate limit exceeded")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
)

type APIErrorDetail struct {
	Field    string
	Messages []string
}

type APIError struct {
	StatusCode int
	Message    string
	Details    []APIErrorDetail
	sentinel   error
}

// ErrorResponse kept for backward compatibility.
type ErrorResponse = APIError

func NewAPIError(statusCode int, sentinel error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		sentinel:   sentinel,
	}
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.sentinel != nil {
		return e.sentinel.Error()
	}

	return "api error"
}

func (e *APIError) Is(target error) bool {
	return errors.Is(e.sentinel, target)
}

func (e *APIError) Unwrap() error {
	return e.sentinel
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var raw struct {
		Message      string              `json:"Message"`
		MessageLower string              `json:"message"`
		Errors       map[string][]string `json:"Errors"`
		ErrorsLower  map[string][]string `json:"errors"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.Message = raw.Message
	if e.Message == "" {
		e.Message = raw.MessageLower
	}

	e.Details = e.Details[:0]

	errorsMap := raw.Errors
	if len(errorsMap) == 0 {
		errorsMap = raw.ErrorsLower
	}

	for field, messages := range errorsMap {
		e.Details = append(e.Details, APIErrorDetail{
			Field:    field,
			Messages: messages,
		})
	}

	return nil
}
