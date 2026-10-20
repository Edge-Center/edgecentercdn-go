package lecerts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
)

type service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) Service {
	return &service{r: r}
}

func (s *service) Issue(ctx context.Context, resourceID int64, req *IssueRequest) error {
	u := url.URL{
		Path: fmt.Sprintf("/resources/%d/ssl/le/issue", resourceID),
	}

	var body interface{}
	if req != nil && req.CertType != "" {
		body = req
	}

	if err := s.r.Request(ctx, http.MethodPost, u.String(), body, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *service) Renew(ctx context.Context, resourceID int64) error {
	u := url.URL{
		Path: fmt.Sprintf("/resources/%d/ssl/le/renew", resourceID),
	}

	if err := s.r.Request(ctx, http.MethodPost, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *service) Revoke(ctx context.Context, resourceID int64, force bool) error {
	u := url.URL{
		Path: fmt.Sprintf("/resources/%d/ssl/le/revoke", resourceID),
	}
	q := u.Query()
	q.Set("force", strconv.FormatBool(force))
	u.RawQuery = q.Encode()

	if err := s.r.Request(ctx, http.MethodPost, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *service) Get(ctx context.Context, resourceID int64) (*Status, error) {
	var status Status

	u := url.URL{
		Path: fmt.Sprintf("/resources/%d/ssl/le/status", resourceID),
	}

	if err := s.r.Request(ctx, http.MethodGet, u.String(), nil, &status); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &status, nil
}

func (s *service) Cancel(ctx context.Context, resourceID int64) error {
	u := url.URL{
		Path: fmt.Sprintf("/resources/%d/ssl/le/status", resourceID),
	}
	body := map[string]bool{
		"active": false,
	}

	if err := s.r.Request(ctx, http.MethodPut, u.String(), body, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
