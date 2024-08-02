package shielding

import (
	"context"
)

type ShieldingService interface {
	Get(ctx context.Context, resourceID int64) (*ShieldingData, error)
	Update(ctx context.Context, resourceID int64, req *UpdateShieldingData) (*ShieldingData, error)
	GetShieldingLocations(ctx context.Context) (*[]ShieldingLocations, error)
}

type ShieldingData struct {
	ShieldingPop *int `json:"shielding_pop"`
}

type UpdateShieldingData struct {
	ShieldingPop *int `json:"shielding_pop"`
}

type ShieldingLocations struct {
	ID         int    `json:"id"`
	Datacenter string `json:"datacenter"`
	Country    string `json:"country"`
	City       string `json:"city"`
}
