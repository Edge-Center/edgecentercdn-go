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

func (s *Service) CreateLECert(ctx context.Context, resourceID int64) error {
	if err := s.r.Request(ctx, http.MethodPost, fmt.Sprintf("/cdn/resources/%d/ssl/le/issue", resourceID), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}
	return nil
}

func (s *Service) UpdateLECert(ctx context.Context, resourceID int64) error {
	if err := s.r.Request(ctx, http.MethodPost, fmt.Sprintf("/cdn/resources/%d/ssl/le/renew", resourceID), nil, nil); err != nil {
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
