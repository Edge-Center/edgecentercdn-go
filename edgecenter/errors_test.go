package edgecenter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

func TestAPIError_ErrorText(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		sentinel error
		body     string
		want     string
	}{
		{
			name:     "field errors are appended to the status text",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"weight":["Ensure this value is less than or equal to 2147483647."]}}`,
			want:     "bad request: weight: Ensure this value is less than or equal to 2147483647.",
		},
		{
			name:     "several fields are listed in stable order",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"useNext":["Must be a valid boolean."],"name":["This field is required."],"origins":["This field is required."]}}`,
			want:     "bad request: name: This field is required.; origins: This field is required.; useNext: Must be a valid boolean.",
		},
		{
			name:     "nested object errors keep the full path",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"origins":{"nonexistent.invalid":{"source":["The domain name cannot be resolved. Please specify a valid domain name."]}}}}`,
			want:     "bad request: origins.nonexistent.invalid.source: The domain name cannot be resolved. Please specify a valid domain name.",
		},
		{
			name:     "list item errors are addressed by index",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"origins":[{},{"source":["Invalid source"]}]}}`,
			want:     "bad request: origins.1.source: Invalid source",
		},
		{
			name:     "plain string field error",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"resource":"Resource doesn't have any certificates."}}`,
			want:     "bad request: resource: Resource doesn't have any certificates.",
		},
		{
			name:     "errors without field names",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":["first problem","second problem"]}`,
			want:     "bad request: first problem; second problem",
		},
		{
			name:     "message only",
			status:   403,
			sentinel: ErrForbidden,
			body:     `{"message":"You do not have permission to perform this action."}`,
			want:     "You do not have permission to perform this action.",
		},
		{
			name:     "message with field errors",
			status:   403,
			sentinel: ErrForbidden,
			body:     `{"message":"access denied","errors":{"token":["invalid token"]}}`,
			want:     "access denied: token: invalid token",
		},
		{
			name:     "structured message is shown as field errors",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"message":{"resource":["Resource is suspended."]}}`,
			want:     "bad request: resource: Resource is suspended.",
		},
		{
			name:     "detail is used as message",
			status:   404,
			sentinel: ErrNotFound,
			body:     `{"detail":"Not found."}`,
			want:     "Not found.",
		},
		{
			name:     "empty body object keeps the status text",
			status:   404,
			sentinel: ErrNotFound,
			body:     `{}`,
			want:     "resource not found",
		},
		{
			name:     "unknown body object follows the status text",
			status:   401,
			sentinel: ErrUnauthorized,
			body:     `{"error":"invalid token"}`,
			want:     `unauthorized: {"error":"invalid token"}`,
		},
		{
			name:   "field errors without a known status",
			status: 500,
			body:   `{"errors":{"field":["broken"]}}`,
			want:   "api error (HTTP 500): field: broken",
		},
		{
			name:   "status without a sentinel keeps the code",
			status: 503,
			body:   `{"error":"Error occurred."}`,
			want:   `api error (HTTP 503): {"error":"Error occurred."}`,
		},
		{
			name:     "empty field errors keep the status text",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"weight":[]}}`,
			want:     "bad request",
		},
		{
			name:     "null message keeps the status text",
			status:   404,
			sentinel: ErrNotFound,
			body:     `{"message":null}`,
			want:     "resource not found",
		},
		{
			name:     "keys are matched regardless of case",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"MESSAGE":"Quota exceeded","ERRORS":{"name":["required"]}}`,
			want:     "Quota exceeded: name: required",
		},
		{
			name:     "the same errors under two spellings are not repeated",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"Errors":{"a":["x"]},"errors":{"a":["x"]}}`,
			want:     "bad request: a: x",
		},
		{
			name:     "numbers keep their digits",
			status:   400,
			sentinel: ErrBadRequest,
			body:     `{"errors":{"weight":[2147483648],"id":9007199254740993}}`,
			want:     "bad request: id: 9007199254740993; weight: 2147483648",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := NewAPIError(tt.status, tt.sentinel)

			require.NoError(t, json.Unmarshal([]byte(tt.body), apiErr))
			assert.Equal(t, tt.want, apiErr.Error())

			if tt.sentinel != nil {
				assert.True(t, errors.Is(apiErr, tt.sentinel))
			}
		})
	}
}

func TestAPIError_UnmarshalJSON_DeepNesting(t *testing.T) {
	const depth = 1000

	body := `{"errors":` + strings.Repeat(`{"a":{"x":["sibling"]},"b":`, depth) + `["leaf"]` + strings.Repeat("}", depth) + "}"

	var apiErr APIError
	require.NoError(t, json.Unmarshal([]byte(body), &apiErr))

	require.Len(t, apiErr.Details, depth+1)
	assert.Equal(t, "a.x", apiErr.Details[0].Field)
	assert.Equal(t, strings.Repeat("b.", depth-1)+"b", apiErr.Details[depth].Field)
	assert.Equal(t, []string{"leaf"}, apiErr.Details[depth].Messages)
	assert.Equal(t, strings.Repeat("b.", depth-1)+"a.x", apiErr.Details[depth-1].Field)
}

func TestAPIError_UnmarshalJSON_NestedDetails(t *testing.T) {
	data := []byte(`{"errors":{"origins":[{},{"source":["Invalid source"],"backup":["Must be a valid boolean."]}],"name":["This field is required."]}}`)

	var apiErr APIError
	require.NoError(t, json.Unmarshal(data, &apiErr))

	assert.Empty(t, apiErr.Message)
	assert.Equal(t, []APIErrorDetail{
		{Field: "name", Messages: []string{"This field is required."}},
		{Field: "origins.1.backup", Messages: []string{"Must be a valid boolean."}},
		{Field: "origins.1.source", Messages: []string{"Invalid source"}},
	}, apiErr.Details)
}
