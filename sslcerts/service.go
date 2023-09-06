package sslcerts

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*Cert, error) {
	var cert Cert

	if err := s.r.Request(ctx, http.MethodPost, "/cdn/ssl/certificates", req, &cert); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &cert, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Cert, error) {
	var cert Cert
	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/cdn/ssl/certificates/%d", id), nil, &cert); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &cert, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/cdn/ssl/certificates/%d", id), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
