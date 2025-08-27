package resources

import (
	"context"
	"github.com/Edge-Center/edgecentercdn-go/edgecenter"
	"time"
)

type ResourceService interface {
	Create(ctx context.Context, req *CreateRequest) (*Resource, error)
	Get(ctx context.Context, id int64) (*Resource, error)
	Update(ctx context.Context, id int64, req *UpdateRequest) (*Resource, error)
	Delete(ctx context.Context, resourceID int64) error
	Page(ctx context.Context, offset uint, size uint, filter *ListFilterRequest) (*PaginatedResource, error)
	List(ctx context.Context, filter *ListFilterRequest) ([]Resource, error)
	Count(ctx context.Context, filter *ListFilterRequest) (uint, error)
}

type Protocol string

const (
	HTTPProtocol  Protocol = "HTTP"
	HTTPSProtocol Protocol = "HTTPS"
	MatchProtocol Protocol = "MATCH"
)

type ResourceStatus string

const (
	ActiveResourceStatus    ResourceStatus = "active"
	SuspendedResourceStatus ResourceStatus = "suspended"
	ProcessedResourceStatus ResourceStatus = "processed"
)

type CreateRequest struct {
	Cname              string                      `json:"cname,omitempty"`
	Description        string                      `json:"description"`
	OriginGroup        int                         `json:"originGroup,omitempty"`
	OriginProtocol     Protocol                    `json:"originProtocol,omitempty"`
	Origin             string                      `json:"origin,omitempty"`
	SecondaryHostnames []string                    `json:"secondaryHostnames,omitempty"`
	SSlEnabled         bool                        `json:"sslEnabled"`
	SSLData            int                         `json:"sslData,omitempty"`
	SSLAutomated       bool                        `json:"ssl_automated"`
	IssueLECert        bool                        `json:"le_issue,omitempty"`
	Options            *edgecenter.ResourceOptions `json:"options,omitempty"`
}

type UpdateRequest struct {
	Description        string                      `json:"description"`
	Active             bool                        `json:"active"`
	OriginGroup        int                         `json:"originGroup"`
	OriginProtocol     Protocol                    `json:"originProtocol,omitempty"`
	SecondaryHostnames []string                    `json:"secondaryHostnames,omitempty"`
	SSlEnabled         bool                        `json:"sslEnabled"`
	SSLData            *int                        `json:"sslData"`
	SSLAutomated       bool                        `json:"ssl_automated"`
	Options            *edgecenter.ResourceOptions `json:"options,omitempty"`
}

type ListFilterRequest struct {
	Ordering string
	Search   string
	Fields   []string
	Status   []ResourceStatus
	Cname    string
	Deleted  bool
}

type Resource struct {
	ID                 int64                       `json:"id,omitempty"`
	Name               string                      `json:"name,omitempty"`
	CreatedAt          *time.Time                  `json:"created,omitempty"`
	UpdatedAt          *time.Time                  `json:"updated,omitempty"`
	Status             ResourceStatus              `json:"status,omitempty"`
	Active             bool                        `json:"active,omitempty"`
	Client             int64                       `json:"client,omitempty"`
	OriginGroup        int64                       `json:"originGroup,omitempty"`
	Cname              string                      `json:"cname,omitempty"`
	Description        string                      `json:"description,omitempty"`
	SecondaryHostnames []string                    `json:"secondaryHostnames,omitempty"`
	Shielded           bool                        `json:"shielded,omitempty"`
	Deleted            bool                        `json:"deleted,omitempty"`
	SSlEnabled         bool                        `json:"sslEnabled,omitempty"`
	SSLData            int                         `json:"sslData,omitempty"`
	SSLAutomated       bool                        `json:"ssl_automated,omitempty"`
	SSLLEEnabled       bool                        `json:"ssl_le_enabled,omitempty"`
	OriginProtocol     Protocol                    `json:"originProtocol,omitempty"`
	PrimaryResource    int                         `json:"primary_resource,omitempty"`
	IsPrimary          bool                        `json:"is_primary,omitempty"`
	Options            *edgecenter.ResourceOptions `json:"options,omitempty"`
}

type PaginatedResource struct {
	Count   uint       `json:"count"`
	Results []Resource `json:"results"`
}
