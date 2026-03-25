package edgecenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIError_Error(t *testing.T) {
	err := NewAPIError(404, ErrNotFound)
	err.Message = "Group not found"

	assert.Equal(t, "Group not found", err.Error())
}

func TestAPIError_Error_FallbackToSentinel(t *testing.T) {
	err := NewAPIError(404, ErrNotFound)

	assert.Equal(t, "resource not found", err.Error())
}

func TestAPIError_Unwrap(t *testing.T) {
	err := NewAPIError(404, ErrNotFound)

	assert.Equal(t, ErrNotFound, err.Unwrap())
}

func TestAPIError_Is(t *testing.T) {
	err := NewAPIError(404, ErrNotFound)

	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrUnauthorized))
}

func TestAPIError_As(t *testing.T) {
	src := NewAPIError(404, ErrNotFound)
	src.Message = "Group not found"

	err := fmt.Errorf("get resource: %w", src)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Equal(t, "Group not found", apiErr.Message)
}

func TestAPIError_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"Message": "validation failed",
		"Errors": {
			"group": ["Group not found"],
			"source": ["Invalid source"]
		}
	}`)

	var apiErr APIError
	err := json.Unmarshal(data, &apiErr)

	require.NoError(t, err)
	assert.Equal(t, "validation failed", apiErr.Message)
	assert.Len(t, apiErr.Details, 2)

	details := map[string][]string{}
	for _, detail := range apiErr.Details {
		details[detail.Field] = detail.Messages
	}

	assert.Equal(t, []string{"Group not found"}, details["group"])
	assert.Equal(t, []string{"Invalid source"}, details["source"])
}

func TestAPIError_UnmarshalJSON_LowercaseKeys(t *testing.T) {
	data := []byte(`{
		"message": "validation failed",
		"errors": {
			"group": ["Group not found"],
			"source": ["Invalid source"]
		}
	}`)

	var apiErr APIError
	err := json.Unmarshal(data, &apiErr)

	require.NoError(t, err)
	assert.Equal(t, "validation failed", apiErr.Message)
	assert.Len(t, apiErr.Details, 2)

	details := map[string][]string{}
	for _, detail := range apiErr.Details {
		details[detail.Field] = detail.Messages
	}

	assert.Equal(t, []string{"Group not found"}, details["group"])
	assert.Equal(t, []string{"Invalid source"}, details["source"])
}
