package locations

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

func (s *service) List(ctx context.Context, resourceID int64) ([]Location, error) {
	var locations []Location

	path := fmt.Sprintf("/resources/%d/locations", resourceID)
	if err := s.r.Request(ctx, http.MethodGet, path, nil, &locations); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return locations, nil
}

func (s *service) Get(ctx context.Context, resourceID, id int64) (*Location, error) {
	var location Location

	path := fmt.Sprintf("/resources/%d/locations/%d", resourceID, id)
	if err := s.r.Request(ctx, http.MethodGet, path, nil, &location); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &location, nil
}

func (s *service) Create(ctx context.Context, resourceID int64, req *CreateRequest) (*Location, error) {
	var location Location

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate location create request: %w", err)
	}

	path := fmt.Sprintf("/resources/%d/locations", resourceID)
	if err := s.r.Request(ctx, http.MethodPost, path, req, &location); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &location, nil
}

func (s *service) Update(ctx context.Context, resourceID, id int64, req *UpdateRequest) (*Location, error) {
	var location Location

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate location update request: %w", err)
	}

	path := fmt.Sprintf("/resources/%d/locations/%d", resourceID, id)
	if err := s.r.Request(ctx, http.MethodPut, path, req, &location); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &location, nil
}

func (s *service) Delete(ctx context.Context, resourceID, id int64) error {
	path := fmt.Sprintf("/resources/%d/locations/%d", resourceID, id)
	if err := s.r.Request(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
