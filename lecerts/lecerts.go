package lecerts

import (
	"context"
	"encoding/json"
)

type LECertService interface {
	GetLECert(ctx context.Context, resourceID int64) (*LECertStatus, error)
	CreateLECert(ctx context.Context, resourceID int64) error
	UpdateLECert(ctx context.Context, resourceID int64) error
	DeleteLECert(ctx context.Context, resourceID int64, force bool) error
	CancelLECert(ctx context.Context, resourceID int64, active bool) error
}

type LECertIssueService interface {
	IssueLECert(ctx context.Context, resourceID int64, req *IssueRequest) error
}

type CertType string

const (
	CertTypeLE   CertType = "LE"
	CertTypeMDDC CertType = "MDDC"
)

type IssueRequest struct {
	CertType CertType `json:"cert_type,omitempty"`
}

type LECertStatus struct {
	ID       int              `json:"ssl_id"`
	Statuses []LEStatusDetail `json:"statuses"`
	Started  string           `json:"started"`
	Finished *string          `json:"finished"`
	Active   bool             `json:"active"`
	Resource int              `json:"resource"`
	CertType CertType         `json:"cert_type,omitempty"`
}

type LEStatusDetail struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Details string `json:"details"`
	Created string `json:"created"`
}

func (s *LECertStatus) UnmarshalJSON(data []byte) error {
	type alias LECertStatus

	aux := struct {
		alias
		CertTypeLegacy CertType `json:"cet_type"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*s = LECertStatus(aux.alias)
	if s.CertType == "" {
		s.CertType = aux.CertTypeLegacy
	}

	return nil
}
