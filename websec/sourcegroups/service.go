package sourcegroups

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

func (s *service) List(ctx context.Context) ([]SourceGroup, error) {
	var groups []SourceGroup

	if err := s.r.Request(ctx, http.MethodGet, "/source_groups", nil, &groups); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return groups, nil
}

func (s *service) Get(ctx context.Context, id int64) (*SourceGroup, error) {
	var group SourceGroup

	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/source_groups/%d", id), nil, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &group, nil
}

func (s *service) Create(ctx context.Context, req *GroupRequest) (*SourceGroup, error) {
	var group SourceGroup

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate source group create request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPost, "/source_groups", req, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &group, nil
}

func (s *service) Update(ctx context.Context, id int64, req *GroupRequest) (*SourceGroup, error) {
	var group SourceGroup

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate source group update request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPut, fmt.Sprintf("/source_groups/%d", id), req, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &group, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/source_groups/%d", id), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
