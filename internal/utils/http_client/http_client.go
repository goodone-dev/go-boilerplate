package httpclient

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/goodone-dev/go-boilerplate/internal/utils/breaker"
	"github.com/sony/gobreaker/v2"
)

var httpClient = resty.New().
	SetDebug(false).
	SetRetryCount(1).
	SetRetryWaitTime(1 * time.Second).
	AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500 && r.StatusCode() <= 599
		},
	)

var breakerMap = make(map[string]*gobreaker.CircuitBreaker[*resty.Response])

type customHttpClient struct {
	Request customHttpRequest
}

type customHttpRequest struct {
	*resty.Request
	breaker *gobreaker.CircuitBreaker[*resty.Response]
}

func NewHttpClient() *customHttpClient {
	return &customHttpClient{
		Request: customHttpRequest{
			Request: httpClient.NewRequest(),
		},
	}
}

func (c *customHttpClient) WithBreaker() (*customHttpClient, error) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	methodName := parseMethodName(funcName)

	if _, ok := breakerMap[methodName]; !ok {
		breakerMap[methodName] = breaker.NewHttpBreaker(methodName)
	}

	c.Request.breaker = breakerMap[methodName]
	if c.Request.breaker.State() == gobreaker.StateOpen {
		return nil, fmt.Errorf("circuit breaker is open for %s", methodName)
	}

	return c, nil
}

func (r *customHttpRequest) Get(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.get(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.get(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) get(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Post(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.post(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.post(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) post(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Post(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Put(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.put(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.put(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) put(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Put(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Patch(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.patch(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.patch(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) patch(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Patch(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Delete(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.delete(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.delete(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) delete(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Delete(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Head(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.head(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.head(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) head(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Head(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Options(ctx context.Context, url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.options(ctx, url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.options(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) options(ctx context.Context, url string) (*resty.Response, error) {
	res, err := r.Request.SetContext(ctx).Options(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func parseMethodName(funcName string) string {
	re := regexp.MustCompile(`\(\*?([^)]+)\)\.([^.]+)$`)
	matches := re.FindStringSubmatch(funcName)

	if len(matches) == 3 {
		typeName := matches[1]
		methodName := matches[2]

		re = regexp.MustCompile(`([^.]+)$`)
		typeName = re.FindString(typeName)

		return typeName + "." + methodName
	}

	return funcName
}
