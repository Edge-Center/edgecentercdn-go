package origingroups

import (
	"context"
)

type OriginGroupService interface {
	Create(ctx context.Context, req *GroupRequest) (*OriginGroup, error)
	Get(ctx context.Context, id int64) (*OriginGroup, error)
	Update(ctx context.Context, id int64, req *GroupRequest) (*OriginGroup, error)
	Delete(ctx context.Context, id int64) error
}

type GroupRequest struct {
	Name                string          `json:"name"`
	UseNext             bool            `json:"useNext"`
	Origins             []OriginRequest `json:"origins"`
	Authorization       *Authorization  `json:"authorization"`
	ConsistentBalancing bool            `json:"consistent_balancing"`
}

type OriginRequest struct {
	Source  string `json:"source"`
	Backup  bool   `json:"backup"`
	Enabled bool   `json:"enabled"`
}

type OriginGroup struct {
	ID                  int64          `json:"id"`
	Name                string         `json:"name"`
	UseNext             bool           `json:"useNext"`
	Origins             []Origin       `json:"origin_ids"`
	Authorization       *Authorization `json:"authorization"`
	ConsistentBalancing bool           `json:"consistent_balancing"`
}

type Origin struct {
	ID      int64  `json:"id"`
	Source  string `json:"source"`
	Backup  bool   `json:"backup"`
	Enabled bool   `json:"enabled"`
}

type Authorization struct {
	AuthType        string `json:"auth_type"`
	AccessKeyID     string `json:"access_key_id"`
	AddressingStyle string `json:"addressing_style"`
	AwsRegion       string `json:"aws_region"`
	SecretKey       string `json:"secret_key"`
	BucketName      string `json:"bucket_name"`
}
