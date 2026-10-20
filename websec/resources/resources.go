package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/Edge-Center/edgecentercdn-go/websec/locations"
	"github.com/Edge-Center/edgecentercdn-go/websec/options"
)

type Service interface {
	List(ctx context.Context) ([]Resource, error)
	Get(ctx context.Context, id int64) (*Resource, error)
	Create(ctx context.Context, req *CreateRequest) (*Resource, error)
	Update(ctx context.Context, id int64, req *UpdateRequest) (*Resource, error)
	Delete(ctx context.Context, id int64) error
}

type Status string

const (
	StatusProcessed Status = "processed"
	StatusActive    Status = "active"
	StatusError     Status = "error"
	StatusSuspended Status = "suspended"
)

type CreateRequest struct {
	Domain         string                   `json:"domain"`
	Subdomains     []string                 `json:"subdomains,omitempty"`
	Active         *bool                    `json:"active,omitempty"`
	SSLCertificate *int64                   `json:"ssl_certificate,omitempty"`
	SSLAutomated   bool                     `json:"ssl_automated,omitempty"`
	SourceGroup    int64                    `json:"source_group"`
	Options        *options.ResourceOptions `json:"options,omitempty"`
}

type UpdateRequest struct {
	Subdomains     options.Field[[]string]  `json:"subdomains,omitzero"`
	Active         bool                     `json:"active"`
	SSLCertificate options.Field[int64]     `json:"ssl_certificate,omitzero"`
	SSLAutomated   *bool                    `json:"ssl_automated,omitempty"`
	SourceGroup    int64                    `json:"source_group"`
	Options        *options.ResourceOptions `json:"options,omitempty"`
}

type Resource struct {
	ID             int64                    `json:"id"`
	Domain         string                   `json:"domain"`
	Subdomains     []string                 `json:"subdomains"`
	Active         bool                     `json:"active"`
	Enabled        bool                     `json:"enabled"`
	Status         Status                   `json:"status"`
	ProtectedIP    *string                  `json:"protected_ip"`
	SSLCertificate *int64                   `json:"ssl_certificate"`
	SSLAutomated   bool                     `json:"ssl_automated"`
	SourceGroup    *int64                   `json:"source_group"`
	Client         int64                    `json:"client"`
	Locations      []locations.Location     `json:"locations"`
	Options        *options.ResourceOptions `json:"options"`
	Created        time.Time                `json:"created"`
	Updated        time.Time                `json:"updated"`
	Deleted        bool                     `json:"deleted"`
}

func (r *CreateRequest) Validate() error {
	if r.Domain == "" {
		return fmt.Errorf("domain is required")
	}
	if r.SourceGroup == 0 {
		return fmt.Errorf("source_group is required")
	}

	return nil
}

func (r *UpdateRequest) Validate() error {
	if r.SourceGroup == 0 {
		return fmt.Errorf("source_group is required")
	}

	return nil
}
