package origingroups

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

var _ OriginGroupService = (*Service)(nil)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Create(ctx context.Context, req *GroupRequest) (*OriginGroup, error) {
	var group OriginGroup

	if err := s.r.Request(ctx, http.MethodPost, "/cdn/source_groups", req, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	if req.Authorization != nil {
		auth, err := s.manageAuth(ctx, group.ID, false, req.Authorization)
		if err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}
		group.Authorization = auth
	}

	return &group, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*OriginGroup, error) {
	var group OriginGroup

	if err := s.r.Request(ctx, http.MethodGet, fmt.Sprintf("/cdn/source_groups/%d", id), nil, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &group, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *GroupRequest) (*OriginGroup, error) {
	var group OriginGroup

	if err := s.r.Request(ctx, http.MethodPut, fmt.Sprintf("/cdn/source_groups/%d", id), req, &group); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	isUpdate := true
	if group.Authorization == nil {
		isUpdate = false
	}
	auth, err := s.manageAuth(ctx, id, isUpdate, req.Authorization)
	if err != nil {
		return nil, fmt.Errorf("request: %w, %v", err, auth)
	}
	group.Authorization = auth

	return &group, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/cdn/source_groups/%d", id), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	if _, err := s.manageAuth(ctx, id, false, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *Service) manageAuth(ctx context.Context, groupID int64, isUpdate bool, reqAuth *Authorization) (*Authorization, error) {
	if reqAuth == nil {
		if err := s.r.Request(ctx, http.MethodDelete, fmt.Sprintf("/cdn/source_groups/%d/authorization", groupID), nil, nil); err != nil {
			return nil, fmt.Errorf("request: %w", err)
		}
	}

	method := http.MethodPost
	if isUpdate {
		method = http.MethodPut
	}

	authRes := &Authorization{}

	if err := s.r.Request(ctx, method, fmt.Sprintf("/cdn/source_groups/%d/authorization", groupID), reqAuth, authRes); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return authRes, nil
}
