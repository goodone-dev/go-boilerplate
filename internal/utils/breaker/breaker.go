package breaker

import (
	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker/v2"
)

func NewHttpBreaker(name string) *gobreaker.CircuitBreaker[*resty.Response] {
	setting := gobreaker.Settings{
		Name: name,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.5
		},
	}

	return gobreaker.NewCircuitBreaker[*resty.Response](setting)
}
