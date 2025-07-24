package lecerts

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

func (s *Service) GetLECert(ctx context.Context, resourceID int64) (*LECertStatus, error) {
	var status LECertStatus
	path := fmt.Sprintf("/cdn/resources/%d/ssl/le/status", resourceID)
	if err := s.r.Request(ctx, http.MethodGet, path, nil, &status); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	return &status, nil
}

func (s *Service) CreateLECert(ctx context.Context, resourceID int64) error {
	path := fmt.Sprintf("/cdn/resources/%d/ssl/le/issue", resourceID)
	if err := s.r.Request(ctx, http.MethodPost, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}
	return nil
}

func (s *Service) UpdateLECert(ctx context.Context, resourceID int64) error {
	path := fmt.Sprintf("/cdn/resources/%d/ssl/le/renew", resourceID)
	if err := s.r.Request(ctx, http.MethodPost, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}
	return nil
}

func (s *Service) DeleteLECert(ctx context.Context, resourceID int64, force bool) error {
	path := fmt.Sprintf("/cdn/resources/%d/ssl/le/revoke?force=%t", resourceID, force)
	if err := s.r.Request(ctx, http.MethodPost, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}
	return nil
}
