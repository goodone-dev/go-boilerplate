package health

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

type HealthResponse struct {
	Status Status `json:"status"`
}
