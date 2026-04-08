package sslcerts

import (
	"context"
	"fmt"
	"time"
)

type SSLCertService interface {
	Create(ctx context.Context, req *CreateRequest) (*Cert, error)
	Get(ctx context.Context, id int64) (*Cert, error)
	Delete(ctx context.Context, id int64) error
}

type Cert struct {
	ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Deleted             bool      `json:"deleted"`
	CertIssuer          string    `json:"cert_issuer"`
	CertSubjectCN       string    `json:"cert_subjeck_cn"`
	ValidityNotBefore   time.Time `json:"validity_not_before"`
	ValidityNotAfter    time.Time `json:"validity_not_after"`
	HasRelatedResources bool      `json:"hasRelatedResources"`
	Automated           bool      `json:"automated"`
}

type CreateRequest struct {
	Name       string `json:"name"`
	Cert       string `json:"sslCertificate"`
	PrivateKey string `json:"sslPrivateKey"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Cert == "" {
		return fmt.Errorf("sslCertificate is required")
	}
	if r.PrivateKey == "" {
		return fmt.Errorf("sslPrivateKey is required")
	}

	return nil
}
