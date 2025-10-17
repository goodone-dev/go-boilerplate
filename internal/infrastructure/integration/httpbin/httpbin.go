package httpbin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	httpclient "github.com/goodone-dev/go-boilerplate/internal/utils/http_client"
)

const url = "https://httpbin.org"

type IHttpbinIntegration interface {
	GetErrorStatus(ctx context.Context) (data any, err error)
	GetSuccessStatus(ctx context.Context) (data any, err error)
}

type HttpbinIntegration struct{}

func NewHttpBinIntegration() IHttpbinIntegration {
	return &HttpbinIntegration{}
}

func (*HttpbinIntegration) GetErrorStatus(ctx context.Context) (body any, err error) {
	_, span := tracer.StartSpan(ctx)
	defer func() {
		span.EndSpan(err, body)
	}()

	http, err := httpclient.NewHttpClient().WithBreaker()
	if err != nil {
		return nil, err
	}

	res, err := http.Request.
		Get(fmt.Sprintf("%s/status/500", url))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &body)
	if err != nil {
		return
	}

	return
}

func (*HttpbinIntegration) GetSuccessStatus(ctx context.Context) (body any, err error) {
	_, span := tracer.StartSpan(ctx)
	defer func() {
		span.EndSpan(err, body)
	}()

	http, err := httpclient.NewHttpClient().WithBreaker()
	if err != nil {
		return nil, err
	}

	res, err := http.Request.
		Get(fmt.Sprintf("%s/headers", url))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &body)
	if err != nil {
		return
	}

	return
}
