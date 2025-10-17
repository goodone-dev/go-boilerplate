package httpclient

import (
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

func (r *customHttpRequest) Get(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.get(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.get(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) get(url string) (*resty.Response, error) {
	res, err := r.Request.Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Post(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.post(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.post(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) post(url string) (*resty.Response, error) {
	res, err := r.Request.Post(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Put(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.put(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.put(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) put(url string) (*resty.Response, error) {
	res, err := r.Request.Put(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Patch(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.patch(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.patch(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) patch(url string) (*resty.Response, error) {
	res, err := r.Request.Patch(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Delete(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.delete(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.delete(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) delete(url string) (*resty.Response, error) {
	res, err := r.Request.Delete(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Head(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.head(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.head(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) head(url string) (*resty.Response, error) {
	res, err := r.Request.Head(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to request %s %s: %s", r.Method, url, res.Error())
	}

	return res, nil
}

func (r *customHttpRequest) Options(url string) (*resty.Response, error) {
	if r.breaker == nil {
		return r.options(url)
	}

	res, err := r.breaker.Execute(func() (*resty.Response, error) {
		return r.options(url)
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *customHttpRequest) options(url string) (*resty.Response, error) {
	res, err := r.Request.Options(url)
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
