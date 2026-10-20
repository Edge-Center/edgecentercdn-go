package websec

import (
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"github.com/Edge-Center/edgecentercdn-go/websec/lecerts"
	"github.com/Edge-Center/edgecentercdn-go/websec/locations"
	"github.com/Edge-Center/edgecentercdn-go/websec/resources"
	"github.com/Edge-Center/edgecentercdn-go/websec/sourcegroups"
	"github.com/Edge-Center/edgecentercdn-go/websec/sslcerts"
)

type ClientService interface {
	Resources() resources.Service
	Locations() locations.Service
	LECerts() lecerts.Service
	SourceGroups() sourcegroups.Service
	SSLCerts() sslcerts.Service
}

var _ ClientService = (*Service)(nil)

type Service struct {
	r                   edgecenter.Requester
	resourcesService    resources.Service
	locationsService    locations.Service
	leCertsService      lecerts.Service
	sourceGroupsService sourcegroups.Service
	sslCertsService     sslcerts.Service
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{
		r:                   r,
		resourcesService:    resources.NewService(r),
		locationsService:    locations.NewService(r),
		leCertsService:      lecerts.NewService(r),
		sourceGroupsService: sourcegroups.NewService(r),
		sslCertsService:     sslcerts.NewService(r),
	}
}

func (s *Service) Resources() resources.Service {
	return s.resourcesService
}

func (s *Service) Locations() locations.Service {
	return s.locationsService
}

func (s *Service) LECerts() lecerts.Service {
	return s.leCertsService
}

func (s *Service) SourceGroups() sourcegroups.Service {
	return s.sourceGroupsService
}

func (s *Service) SSLCerts() sslcerts.Service {
	return s.sslCertsService
}
