package resources

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

func (s *service) List(ctx context.Context) ([]Resource, error) {
	var resources []Resource

	if err := s.r.Request(ctx, http.MethodGet, "/resources", nil, &resources); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return resources, nil
}

func (s *service) Get(ctx context.Context, id int64) (*Resource, error) {
	var resource Resource

	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/resources/%d", id), nil, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*Resource, error) {
	var resource Resource

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate resource create request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPost, "/resources", req, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *service) Update(ctx context.Context, id int64, req *UpdateRequest) (*Resource, error) {
	var resource Resource

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate resource update request: %w", err)
	}

	if err := s.r.Request(ctx, http.MethodPut, fmt.Sprintf("/resources/%d", id), req, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/resources/%d", id), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
