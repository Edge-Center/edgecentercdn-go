package resources

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

var _ ResourceService = (*Service)(nil)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*Resource, error) {
	var resource Resource
	if err := s.r.Request(ctx, http.MethodPost, "/cdn/resources", req, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Resource, error) {
	var resource Resource
	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/cdn/resources/%d", id), nil, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateRequest) (*Resource, error) {
	var resource Resource
	if err := s.r.Request(ctx, http.MethodPut, fmt.Sprintf("/cdn/resources/%d", id), req, &resource); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resource, nil
}

func (s *Service) Delete(ctx context.Context, resourceID int64) error {
	path := fmt.Sprintf("/cdn/resources/%d", resourceID)
	if err := s.r.Request(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *Service) Page(ctx context.Context, offset uint, size uint, filter *ListFilterRequest) (*PaginatedResource, error) {
	var response PaginatedResource

	path := buildListPath(offset, size, filter)

	if err := s.r.Request(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) List(ctx context.Context, filter *ListFilterRequest) ([]Resource, error) {
	var response []Resource

	path := buildListPath(0, 0, filter)

	if err := s.r.Request(ctx, http.MethodGet, path, nil, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) Count(ctx context.Context, filter *ListFilterRequest) (uint, error) {
	page, err := s.Page(ctx, 0, 1, filter)
	if err != nil {
		return 0, err
	}
	return page.Count, nil
}
