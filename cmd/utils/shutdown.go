package utils

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
)

type Service interface {
	Shutdown(ctx context.Context) error
}

func GracefulShutdown(ctx context.Context, services ...Service) {
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)

		go func(s Service) {
			defer wg.Done()

			packageName := parsePackageName(s)

			if err := s.Shutdown(ctx); err != nil {
				logger.Errorf(ctx, err, "❌ %s forced to shutdown due to error", packageName)
			}

			logger.Infof(ctx, "✅ %s shutdown gracefully", packageName)
		}(service)
	}

	wg.Wait()
}

func parsePackageName(service Service) string {
	n := reflect.TypeOf(service).String()
	r := regexp.MustCompile(`\*?([^.]+)`)

	matches := r.FindStringSubmatch(n)
	if len(matches) > 1 {
		name := matches[1]
		if len(name) > 0 {
			return strings.ToUpper(string(name[0])) + name[1:]
		}
		return name
	}

	return ""
}
