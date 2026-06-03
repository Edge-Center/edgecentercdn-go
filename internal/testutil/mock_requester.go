package testutil

import (
	"context"
	"encoding/json"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
)

var _ edgecenter.Requester = (*MockRequester)(nil)

type Call struct {
	Method, Path string
	Payload      interface{}
}

type MockRequester struct {
	RequestFunc func(ctx context.Context, method, path string, payload, result interface{}) error
	Calls       []Call
}

func (m *MockRequester) Request(ctx context.Context, method, path string, payload, result interface{}) error {
	m.Calls = append(m.Calls, Call{Method: method, Path: path, Payload: payload})
	if m.RequestFunc != nil {
		return m.RequestFunc(ctx, method, path, payload, result)
	}
	return nil
}

func (m *MockRequester) RespondWith(v interface{}) {
	m.RequestFunc = func(_ context.Context, _, _ string, _, result interface{}) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, result)
	}
}

func (m *MockRequester) RespondWithError(err error) {
	m.RequestFunc = func(context.Context, string, string, interface{}, interface{}) error {
		return err
	}
}
