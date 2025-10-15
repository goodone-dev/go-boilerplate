package utils

import (
	"context"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
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
				log.Printf("❌ %s forced to shutdown: %v", packageName, err)
				return
			}

			log.Printf("✅ %s shutdown gracefully.\n", packageName)
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
		return strings.ToUpper(string(name[0])) + name[1:]
	}

	return ""
}
