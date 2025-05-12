package lecerts

import (
	"context"
)

type LECertService interface {
	CreateLECert(ctx context.Context, resourceID int64) error
	UpdateLECert(ctx context.Context, resourceID int64) error
	DeleteLECert(ctx context.Context, resourceID int64, force bool) error
}
