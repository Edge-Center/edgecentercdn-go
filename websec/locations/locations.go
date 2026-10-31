package locations

import (
	"context"
	"fmt"

	"github.com/Edge-Center/edgecentercdn-go/websec/options"
)

type Service interface {
	List(ctx context.Context, resourceID int64) ([]Location, error)
	Get(ctx context.Context, resourceID, id int64) (*Location, error)
	Create(ctx context.Context, resourceID int64, req *CreateRequest) (*Location, error)
	Update(ctx context.Context, resourceID, id int64, req *UpdateRequest) (*Location, error)
	Delete(ctx context.Context, resourceID, id int64) error
}

type CreateRequest struct {
	Name        string                   `json:"name,omitempty"`
	Location    string                   `json:"location"`
	Active      *bool                    `json:"active,omitempty"`
	SourceGroup *int64                   `json:"source_group,omitempty"`
	Weight      *int                     `json:"weight,omitempty"`
	Options     *options.LocationOptions `json:"options,omitempty"`
}

type UpdateRequest struct {
	Name        string                   `json:"name,omitempty"`
	Location    string                   `json:"location"`
	Active      bool                     `json:"active"`
	SourceGroup *int64                   `json:"source_group"`
	Weight      *int                     `json:"weight,omitempty"`
	Options     *options.LocationOptions `json:"options,omitempty"`
}

type Location struct {
	ID          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Location    string                   `json:"location"`
	Active      bool                     `json:"active"`
	SourceGroup *int64                   `json:"source_group"`
	Weight      int                      `json:"weight"`
	Deleted     bool                     `json:"deleted"`
	Options     *options.LocationOptions `json:"options"`
}

func (r *CreateRequest) Validate() error {
	if r.Location == "" {
		return fmt.Errorf("location is required")
	}

	return nil
}

func (r *UpdateRequest) Validate() error {
	if r.Location == "" {
		return fmt.Errorf("location is required")
	}

	return nil
}
