package health

import "context"

type HealthService interface {
	Ping(ctx context.Context) error
}

type HealthStatus struct {
	Status Status `json:"status"`
}

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)
