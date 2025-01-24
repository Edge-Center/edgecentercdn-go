package tools

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

var _ ResourceToolsService = (*Service)(nil)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) Purge(ctx context.Context, id int64, req *PurgeRequest) (*PurgeResponse, error) {
	var result PurgeResponse

	if req == nil {
		req = &PurgeRequest{Paths: []string{}}
	}

	if err := s.r.Request(ctx, http.MethodPost, fmt.Sprintf("/cdn/resources/%d/purge", id), req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Service) Whoami(ctx context.Context) (string, error) {
	var result WhoamiResponse

	if err := s.r.Request(ctx, http.MethodGet, "/iam/users/me", nil, &result); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %d", result.Email, result.Client), nil
}
