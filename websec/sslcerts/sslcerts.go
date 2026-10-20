package sslcerts

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	List(ctx context.Context) ([]Cert, error)
	Get(ctx context.Context, id int64) (*Cert, error)
	Create(ctx context.Context, req *CreateRequest) (*Cert, error)
	Update(ctx context.Context, id int64, req *UpdateRequest) (*Cert, error)
	Delete(ctx context.Context, id int64) error
}

type Cert struct {
	ID                  int64             `json:"id"`
	Name                string            `json:"name"`
	Deleted             bool              `json:"deleted"`
	Automated           bool              `json:"automated"`
	CertIssuer          string            `json:"cert_issuer"`
	CertSubjectCN       string            `json:"cert_subject_cn"`
	CertSubjectAlt      string            `json:"cert_subject_alt"`
	ValidityNotBefore   *time.Time        `json:"validity_not_before"`
	ValidityNotAfter    *time.Time        `json:"validity_not_after"`
	RelatedResources    []RelatedResource `json:"related_resources"`
	HasRelatedResources bool              `json:"has_related_resources"`
}

type RelatedResource struct {
	ID     int64  `json:"id"`
	Domain string `json:"domain"`
}

type CreateRequest struct {
	Name              string `json:"name"`
	Certificate       string `json:"certificate"`
	PrivateKey        string `json:"private_key"`
	CertificatesChain string `json:"certificates_chain,omitempty"`
}

type UpdateRequest struct {
	Name string `json:"name"`
}

func (r *CreateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Certificate == "" {
		return fmt.Errorf("certificate is required")
	}
	if r.PrivateKey == "" {
		return fmt.Errorf("private_key is required")
	}

	return nil
}

func (r *UpdateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}
