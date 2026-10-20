package options

import (
	"bytes"
	"encoding/json"
)

type fieldState uint8

const (
	fieldAbsent fieldState = iota
	fieldNull
	fieldValue
)

type Field[T any] struct {
	value T
	state fieldState
}

func Value[T any](v T) Field[T] {
	return Field[T]{value: v, state: fieldValue}
}

func Null[T any]() Field[T] {
	return Field[T]{state: fieldNull}
}

func (f Field[T]) IsZero() bool {
	return f.state == fieldAbsent
}

func (f Field[T]) IsNull() bool {
	return f.state == fieldNull
}

func (f Field[T]) Get() (T, bool) {
	return f.value, f.state == fieldValue
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	if f.state != fieldValue {
		return []byte("null"), nil
	}

	return json.Marshal(f.value)
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*f = Null[T]()
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*f = Value(v)

	return nil
}
