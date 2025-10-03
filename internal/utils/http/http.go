package http

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker/v2"
)

var httpClient *resty.Client
var breakerMap map[string]*gobreaker.CircuitBreaker[*resty.Response] = make(map[string]*gobreaker.CircuitBreaker[*resty.Response])

type customHttp struct {
	Request customRequest
}

type customRequest struct {
	*resty.Request
	breaker *gobreaker.CircuitBreaker[*resty.Response]
}

func NewClient() *customHttp {
	if httpClient != nil {
		return &customHttp{
			Request: customRequest{
				Request: httpClient.NewRequest(),
			},
		}
	}

	httpClient = resty.New()
	httpClient.SetDebug(false)
	// httpClient.SetRetryCount(1)
	// httpClient.SetRetryWaitTime(1 * time.Second)
	// httpClient.AddRetryCondition(
	// 	func(r *resty.Response, err error) bool {
	// 		return r.StatusCode() >= 500 && r.StatusCode() <= 599
	// 	},
	// )

	return &customHttp{
		Request: customRequest{
			Request: httpClient.NewRequest(),
		},
	}
}

func (c *customHttp) WithBreaker() (*customHttp, error) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	methodName := parseMethodName(funcName)

	if _, ok := breakerMap[methodName]; !ok {
		breakerMap[methodName] = createBreaker(methodName)
	}

	c.Request.breaker = breakerMap[methodName]
	if c.Request.breaker.State() == gobreaker.StateOpen {
		return nil, errors.New("circuit breaker is open")
	}

	return c, nil
}

func (r *customRequest) Get(url string) (*resty.Response, error) {
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

func (r *customRequest) get(url string) (*resty.Response, error) {
	res, err := r.Request.Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Post(url string) (*resty.Response, error) {
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

func (r *customRequest) post(url string) (*resty.Response, error) {
	res, err := r.Request.Post(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Put(url string) (*resty.Response, error) {
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

func (r *customRequest) put(url string) (*resty.Response, error) {
	res, err := r.Request.Put(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Patch(url string) (*resty.Response, error) {
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

func (r *customRequest) patch(url string) (*resty.Response, error) {
	res, err := r.Request.Patch(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Delete(url string) (*resty.Response, error) {
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

func (r *customRequest) delete(url string) (*resty.Response, error) {
	res, err := r.Request.Delete(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Head(url string) (*resty.Response, error) {
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

func (r *customRequest) head(url string) (*resty.Response, error) {
	res, err := r.Request.Head(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func (r *customRequest) Options(url string) (*resty.Response, error) {
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

func (r *customRequest) options(url string) (*resty.Response, error) {
	res, err := r.Request.Options(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("request failed with status: %s", res.Status())
	}

	return res, nil
}

func createBreaker(name string) *gobreaker.CircuitBreaker[*resty.Response] {
	var setting gobreaker.Settings
	setting.Name = name
	setting.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.5
	}

	return gobreaker.NewCircuitBreaker[*resty.Response](setting)
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
