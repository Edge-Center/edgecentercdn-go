package edgecenter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
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
	head := e.Message
	if head == "" {
		head = e.statusText()
	}

	if details := e.detailsText(); details != "" {
		return head + ": " + details
	}

	return head
}

func (e *APIError) statusText() string {
	if e.sentinel != nil {
		return e.sentinel.Error()
	}

	if e.StatusCode != 0 {
		return fmt.Sprintf("api error (HTTP %d)", e.StatusCode)
	}

	return "api error"
}

func (e *APIError) detailsText() string {
	parts := make([]string, 0, len(e.Details))
	for _, detail := range e.Details {
		for _, message := range detail.Messages {
			if detail.Field == "" {
				parts = append(parts, message)
				continue
			}

			parts = append(parts, detail.Field+": "+message)
		}
	}

	return strings.Join(parts, "; ")
}

func (e *APIError) Is(target error) bool {
	return errors.Is(e.sentinel, target)
}

func (e *APIError) Unwrap() error {
	return e.sentinel
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		return err
	}

	e.Message = ""
	e.Details = e.Details[:0]

	known := false
	var fieldErrors []APIErrorDetail

	for _, key := range slices.Sorted(maps.Keys(raw)) {
		value := raw[key]

		switch strings.ToLower(key) {
		case "message", "detail":
			known = true
			if message, ok := value.(string); ok {
				if e.Message == "" {
					e.Message = message
				}

				continue
			}

			e.Details = append(e.Details, collectDetails(nil, value)...)
		case "errors":
			known = true
			if len(fieldErrors) == 0 {
				fieldErrors = collectDetails(nil, value)
			}
		}
	}

	e.Details = append(e.Details, fieldErrors...)

	if !known && len(raw) > 0 {
		e.Details = append(e.Details, APIErrorDetail{Messages: []string{strings.TrimSpace(string(data))}})
	}

	return nil
}

func collectDetails(path []string, value interface{}) []APIErrorDetail {
	switch value := value.(type) {
	case nil:
		return nil
	case []interface{}:
		return collectListDetails(path, value)
	case map[string]interface{}:
		details := make([]APIErrorDetail, 0, len(value))
		for _, key := range slices.Sorted(maps.Keys(value)) {
			details = append(details, collectDetails(append(path, key), value[key])...)
		}

		return details
	default:
		return []APIErrorDetail{{Field: strings.Join(path, "."), Messages: []string{fmt.Sprint(value)}}}
	}
}

func collectListDetails(path []string, items []interface{}) []APIErrorDetail {
	var messages []string
	var nested []APIErrorDetail

	for i, item := range items {
		switch item.(type) {
		case nil:
		case []interface{}, map[string]interface{}:
			nested = append(nested, collectDetails(append(path, strconv.Itoa(i)), item)...)
		default:
			messages = append(messages, fmt.Sprint(item))
		}
	}

	if len(messages) == 0 {
		return nested
	}

	return append([]APIErrorDetail{{Field: strings.Join(path, "."), Messages: messages}}, nested...)
}
