package sslcerts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
)

type service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) Service {
	return &service{r: r}
}

func (s *service) List(ctx context.Context) ([]Cert, error) {
	var certs []Cert

	if err := s.r.Request(ctx, http.MethodGet, "/ssl/certificates", nil, &certs); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return certs, nil
}

func (s *service) Get(ctx context.Context, id int64) (*Cert, error) {
	var cert Cert

	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/ssl/certificates/%d", id), nil, &cert); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &cert, nil
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*Cert, error) {
	var cert Cert

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate ssl cert create request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPost, "/ssl/certificates", req, &cert); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &cert, nil
}

func (s *service) Update(ctx context.Context, id int64, req *UpdateRequest) (*Cert, error) {
	var cert Cert

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate ssl cert update request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPut, fmt.Sprintf("/ssl/certificates/%d", id), req, &cert); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &cert, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/ssl/certificates/%d", id), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
