package tools

import "context"

type ResourceToolsService interface {
	Purge(ctx context.Context, id int64, req *PurgeRequest) (*PurgeResponse, error)
	Whoami(ctx context.Context) (string, error)
}

type PurgeRequest struct {
	Paths []string `json:"paths"`
}

type PurgeResponse struct {
	Paths []string `json:"paths"`
}

type WhoamiResponse struct {
	ID     int64  `json:"id"`
	Client int64  `json:"client"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
