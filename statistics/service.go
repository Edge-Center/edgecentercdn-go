package statistics

import (
	"context"
	"fmt"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"net/http"
)

type Service struct {
	r edgecenter.Requester
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{r: r}
}

func (s *Service) GetTimeSeriesData(ctx context.Context, req *ResourceStatisticsTimeSeriesRequest) (*ResourceStatisticsTimeSeriesResponse, error) {
	var response ResourceStatisticsTimeSeriesResponse

	if err := s.r.Request(ctx, http.MethodGet, req.ToPath(), nil, &response); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &response, nil
}

func (s *Service) GetTableData(ctx context.Context, req *ResourceStatisticsTableRequest) (*ResourceStatisticsTableResponse, error) {
	var response ResourceStatisticsTableResponse

	if err := s.r.Request(ctx, http.MethodGet, req.ToPath(), nil, &response); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &response, nil
}

func (r *ResourceStatisticsTableRequest) ToPath() string {
	return buildStatisticsPath(
		r.From,
		r.To,
		"",
		r.GroupBy,
		r.Metrics,
		r.Regions,
		r.Hosts,
		r.Resources,
		r.Clients,
		r.Countries)
}

func (r *ResourceStatisticsTimeSeriesRequest) ToPath() string {
	return buildStatisticsPath(
		r.From,
		r.To,
		r.Granularity,
		r.GroupBy,
		r.Metrics,
		r.Regions,
		r.Hosts,
		r.Resources,
		r.Clients,
		r.Countries)
}
