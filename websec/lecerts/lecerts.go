package lecerts

import (
	"context"
)

type Service interface {
	Issue(ctx context.Context, resourceID int64, req *IssueRequest) error
	Renew(ctx context.Context, resourceID int64) error
	Revoke(ctx context.Context, resourceID int64, force bool) error
	Get(ctx context.Context, resourceID int64) (*Status, error)
	Cancel(ctx context.Context, resourceID int64) error
}

type CertType string

const (
	CertTypeLE   CertType = "LE"
	CertTypeMDDC CertType = "MDDC"
	CertTypeEC   CertType = "EC"
)

type IssueRequest struct {
	CertType CertType `json:"cert_type,omitempty"`
}

type Status struct {
	ID           int64          `json:"id"`
	Resource     int64          `json:"resource"`
	Active       bool           `json:"active"`
	Started      *string        `json:"started"`
	Finished     *string        `json:"finished"`
	ACMECertType CertType       `json:"acme_cert_type"`
	SSLID        *int64         `json:"ssl_id"`
	CertType     *CertType      `json:"cert_type"`
	Statuses     []StatusDetail `json:"statuses"`
}

type StatusDetail struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Error   string `json:"error"`
	Details string `json:"details"`
	Created string `json:"created"`
}
