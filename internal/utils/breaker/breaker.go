package breaker

import (
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/sony/gobreaker/v2"
)

func NewCircuitBreaker[T any](name string) *gobreaker.CircuitBreaker[T] {
	setting := gobreaker.Settings{
		Name:        name,
		MaxRequests: uint32(config.CircuitBreaker.MaxRequests),
		Timeout:     config.CircuitBreaker.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= uint32(config.CircuitBreaker.MinRequests) && failureRatio >= config.CircuitBreaker.FailureRatio
		},
	}

	return gobreaker.NewCircuitBreaker[T](setting)
}
