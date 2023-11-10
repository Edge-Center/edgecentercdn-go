package rules

import (
	"context"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
)

type RulesService interface {
	Create(ctx context.Context, resourceID int64, req *CreateRequest) (*Rule, error)
	Get(ctx context.Context, resourceID, ruleID int64) (*Rule, error)
	Update(ctx context.Context, resourceID, ruleID int64, req *UpdateRequest) (*Rule, error)
	Delete(ctx context.Context, resourceID, ruleID int64) error
}

type CreateRequest struct {
    Active				   bool						   `json:"active,omitempty"`
	Name                   string                      `json:"name,omitempty"`
	Rule                   string                      `json:"rule,omitempty"`
	OriginGroup            *int                        `json:"originGroup"`
	OverrideOriginProtocol *string                     `json:"overrideOriginProtocol"`
	Weight                 int            	           `json:"weight,omitempty"`
	Options                *edgecenter.LocationOptions `json:"options,omitempty"`
}

type UpdateRequest struct {
    Active				   bool						   `json:"active,omitempty"`
	Name                   string                      `json:"name,omitempty"`
	Rule                   string                      `json:"rule,omitempty"`
	OriginGroup            *int                        `json:"originGroup"`
	OverrideOriginProtocol *string                     `json:"overrideOriginProtocol"`
	Weight                 int            	           `json:"weight,omitempty"`
	Options                *edgecenter.LocationOptions `json:"options,omitempty"`
}

type Rule struct {
	ID                     int64                       `json:"id"`
	Name                   string                      `json:"name"`
	Active				   bool						   `json:"active"`
	Deleted                bool                        `json:"deleted"`
	OriginGroup            *int                        `json:"originGroup"`
	OriginProtocol         string                      `json:"originProtocol"`
	OverrideOriginProtocol *string                     `json:"overrideOriginProtocol"`
	Pattern                string                      `json:"rule"`
	Weight                 int            	           `json:"weight"`
	Options                *edgecenter.LocationOptions `json:"options,omitempty"`
}
