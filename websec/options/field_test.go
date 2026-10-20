package options

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fieldHolder struct {
	ID Field[int64] `json:"id,omitzero"`
}

func TestField_States(t *testing.T) {
	tests := []struct {
		name      string
		field     Field[int64]
		wantZero  bool
		wantNull  bool
		wantValue int64
		wantSet   bool
	}{
		{name: "zero value is absent", field: Field[int64]{}, wantZero: true},
		{name: "null", field: Null[int64](), wantNull: true},
		{name: "value", field: Value[int64](42), wantValue: 42, wantSet: true},
		{name: "zero is still a value", field: Value[int64](0), wantSet: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantZero, tt.field.IsZero())
			assert.Equal(t, tt.wantNull, tt.field.IsNull())
			v, ok := tt.field.Get()
			assert.Equal(t, tt.wantSet, ok)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}

func TestField_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		holder fieldHolder
		want   string
	}{
		{name: "absent is omitted", holder: fieldHolder{}, want: `{}`},
		{name: "null is sent as null", holder: fieldHolder{ID: Null[int64]()}, want: `{"id":null}`},
		{name: "value is sent", holder: fieldHolder{ID: Value[int64](7)}, want: `{"id":7}`},
		{name: "zero value is sent", holder: fieldHolder{ID: Value[int64](0)}, want: `{"id":0}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.holder)

			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(b))
		})
	}
}

func TestField_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
		want Field[int64]
	}{
		{name: "missing key stays absent", body: `{}`, want: Field[int64]{}},
		{name: "null key becomes null", body: `{"id": null}`, want: Null[int64]()},
		{name: "value key becomes value", body: `{"id": 7}`, want: Value[int64](7)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var holder fieldHolder

			require.NoError(t, json.Unmarshal([]byte(tt.body), &holder))
			assert.Equal(t, tt.want, holder.ID)
		})
	}
}

func TestField_UnmarshalJSON_ValueOfWrongTypeFails(t *testing.T) {
	var holder fieldHolder

	err := json.Unmarshal([]byte(`{"id": "seven"}`), &holder)

	require.Error(t, err)
	assert.True(t, holder.ID.IsZero())
}
