package edgecentercdn_go

import (
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"github.com/Edge-Center/edgecentercdn-go/origingroups"
	"github.com/Edge-Center/edgecentercdn-go/resources"
	"github.com/Edge-Center/edgecentercdn-go/rules"
	"github.com/Edge-Center/edgecentercdn-go/shielding"
	"github.com/Edge-Center/edgecentercdn-go/sslcerts"
	"github.com/Edge-Center/edgecentercdn-go/statistics"
	"github.com/Edge-Center/edgecentercdn-go/tools"
)

type ClientService interface {
	Resources() resources.ResourceService
	Rules() rules.RulesService
	OriginGroups() origingroups.OriginGroupService
	Shielding() shielding.ShieldingService
	SSLCerts() sslcerts.SSLCertService
	Statistics() statistics.ResourceStatisticsService
	Tools() tools.ResourceToolsService
}

var _ ClientService = (*Service)(nil)

type Service struct {
	r                   edgecenter.Requester
	resourcesService    resources.ResourceService
	statisticsService   statistics.ResourceStatisticsService
	tools               tools.ResourceToolsService
	rulesService        rules.RulesService
	originGroupsService origingroups.OriginGroupService
	shieldingService    shielding.ShieldingService
	sslCertsService     sslcerts.SSLCertService
}

func NewService(r edgecenter.Requester) *Service {
	return &Service{
		r:                   r,
		resourcesService:    resources.NewService(r),
		rulesService:        rules.NewService(r),
		tools:               tools.NewService(r),
		statisticsService:   statistics.NewService(r),
		originGroupsService: origingroups.NewService(r),
		shieldingService:    shielding.NewService(r),
		sslCertsService:     sslcerts.NewService(r),
	}
}

func (s *Service) Resources() resources.ResourceService {
	return s.resourcesService
}

func (s *Service) Rules() rules.RulesService {
	return s.rulesService
}

func (s *Service) OriginGroups() origingroups.OriginGroupService {
	return s.originGroupsService
}

func (s *Service) Shielding() shielding.ShieldingService {
	return s.shieldingService
}

func (s *Service) SSLCerts() sslcerts.SSLCertService {
	return s.sslCertsService
}

func (s *Service) Statistics() statistics.ResourceStatisticsService {
	return s.statisticsService
}

func (s *Service) Tools() tools.ResourceToolsService {
	return s.tools
}
