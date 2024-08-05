package shielding

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

var _ ShieldingService = (*Service)(nil)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Get(ctx context.Context, resourceID int64) (*ShieldingData, error) {
	var shielding ShieldingData

	path := fmt.Sprintf("/cdn/resources/%d/shielding", resourceID)
	if err := s.r.Request(ctx, http.MethodGet, path, nil, &shielding); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &shielding, nil
}

func (s *Service) Update(ctx context.Context, resourceID int64, req *UpdateShieldingData) (*ShieldingData, error) {
	var shielding ShieldingData

	path := fmt.Sprintf("/cdn/resources/%d/shielding", resourceID)

	if err := s.r.Request(ctx, http.MethodPut, path, req, &shielding); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	return &shielding, nil
}

func (s *Service) GetShieldingLocations(ctx context.Context) (*[]ShieldingLocations, error) {
	var shieldingLocations []ShieldingLocations

	if err := s.r.Request(ctx, http.MethodPut, "/cdn/shielding_pop", nil, &shieldingLocations); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	return &shieldingLocations, nil
}
