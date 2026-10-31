package sourcegroups

import (
	"context"
	"fmt"
)

type Service interface {
	List(ctx context.Context) ([]SourceGroup, error)
	Get(ctx context.Context, id int64) (*SourceGroup, error)
	Create(ctx context.Context, req *GroupRequest) (*SourceGroup, error)
	Update(ctx context.Context, id int64, req *GroupRequest) (*SourceGroup, error)
	Delete(ctx context.Context, id int64) error
}

type SourceMode string

const (
	SourceModePrimary SourceMode = "primary"
	SourceModeBackup  SourceMode = "backup"
	SourceModeDown    SourceMode = "down"
)

type GroupRequest struct {
	Name        string   `json:"name"`
	UseNext     *bool    `json:"use_next,omitempty"`
	Description string   `json:"description"`
	Sources     []Source `json:"sources"`
}

type Source struct {
	Source      string     `json:"source"`
	Mode        SourceMode `json:"mode,omitempty"`
	Weight      int        `json:"weight,omitempty"`
	MaxFails    int        `json:"max_fails,omitempty"`
	FailTimeout int        `json:"fail_timeout,omitempty"`
	Comment     string     `json:"comment,omitempty"`
}

type SourceGroup struct {
	ID               int64             `json:"id"`
	Name             string            `json:"name"`
	UseNext          bool              `json:"use_next"`
	Description      string            `json:"description"`
	Sources          []Source          `json:"sources"`
	RelatedResources []RelatedResource `json:"related_resources"`
}

type RelatedResource struct {
	ID        int64             `json:"id"`
	Domain    string            `json:"domain"`
	Locations []RelatedLocation `json:"locations"`
}

type RelatedLocation struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (r *GroupRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(r.Sources) == 0 {
		return fmt.Errorf("sources is required")
	}
	for i, source := range r.Sources {
		if source.Source == "" {
			return fmt.Errorf("sources[%d].source is required", i)
		}
	}

	return nil
}
