package lecerts

import (
	"context"
)

type LECertService interface {
	GetLECert(ctx context.Context, resourceID int64) (*LECertStatus, error)
	CreateLECert(ctx context.Context, resourceID int64) error
	UpdateLECert(ctx context.Context, resourceID int64) error
	DeleteLECert(ctx context.Context, resourceID int64, force bool) error
}

type LECertStatus struct {
	ID       int              `json:"id"`
	Statuses []LEStatusDetail `json:"statuses"`
	Started  string           `json:"started"`
	Finished *string          `json:"finished"`
	Active   bool             `json:"active"`
	Resource int              `json:"resource"`
}

type LEStatusDetail struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Error   string `json:"error"`
	Details string `json:"details"`
}
