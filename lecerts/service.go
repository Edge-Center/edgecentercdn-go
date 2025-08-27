package lecerts

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
	"net/url"
)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) GetLECert(ctx context.Context, resourceID int64) (*LECertStatus, error) {
	var status LECertStatus

	u := url.URL{
		Path: fmt.Sprintf("/cdn/resources/%d/ssl/le/status", resourceID),
	}

	if err := s.r.Request(ctx, http.MethodGet, u.String(), nil, &status); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &status, nil
}

func (s *Service) CreateLECert(ctx context.Context, resourceID int64) error {
	u := url.URL{
		Path: fmt.Sprintf("/cdn/resources/%d/ssl/le/issue", resourceID),
	}

	if err := s.r.Request(ctx, http.MethodPost, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *Service) UpdateLECert(ctx context.Context, resourceID int64) error {
	u := url.URL{
		Path: fmt.Sprintf("/cdn/resources/%d/ssl/le/renew", resourceID),
	}

	if err := s.r.Request(ctx, http.MethodPost, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *Service) DeleteLECert(ctx context.Context, resourceID int64, force bool) error {
	u := url.URL{
		Path: fmt.Sprintf("/cdn/resources/%d/ssl/le/revoke", resourceID),
	}
	q := u.Query()
	q.Set("force", fmt.Sprintf("%t", force))
	u.RawQuery = q.Encode()

	if err := s.r.Request(ctx, http.MethodPost, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *Service) CancelLECert(ctx context.Context, resourceID int64, active bool) error {
	u := url.URL{
		Path: fmt.Sprintf("/cdn/resources/%d/ssl/le/status", resourceID),
	}
	body := map[string]bool{
		"active": active,
	}
	if err := s.r.Request(ctx, http.MethodPut, u.String(), body, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
