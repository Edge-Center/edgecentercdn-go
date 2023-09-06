package rules

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

var _ RulesService = (*Service)(nil)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, resourceID int64, req *CreateRequest) (*Rule, error) {
	var rule Rule

	path := fmt.Sprintf("/cdn/resources/%d/locations", resourceID)
	if err := s.r.Request(ctx, http.MethodPost, path, req, &rule); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &rule, nil
}

func (s *Service) Get(ctx context.Context, resourceID, ruleID int64) (*Rule, error) {
	var rule Rule

	path := fmt.Sprintf("/cdn/resources/%d/locations/%d", resourceID, ruleID)
	if err := s.r.Request(ctx, http.MethodGet, path, nil, &rule); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &rule, nil
}

func (s *Service) Update(ctx context.Context, resourceID, ruleID int64, req *UpdateRequest) (*Rule, error) {
	var rule Rule

	path := fmt.Sprintf("/cdn/resources/%d/locations/%d", resourceID, ruleID)
	if err := s.r.Request(ctx, http.MethodPut, path, req, &rule); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &rule, nil
}

func (s *Service) Delete(ctx context.Context, resourceID, ruleID int64) error {
	path := fmt.Sprintf("/cdn/resources/%d/locations/%d", resourceID, ruleID)
	if err := s.r.Request(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
