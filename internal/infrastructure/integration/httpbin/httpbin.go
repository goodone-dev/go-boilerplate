package httpbin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	httpclient "github.com/goodone-dev/go-boilerplate/internal/utils/http_client"
)

const url = "https://httpbin.org"

type HttpbinIntegration interface {
	GetErrorStatus(ctx context.Context) (data any, err error)
	GetSuccessStatus(ctx context.Context) (data any, err error)
}

type httpbinIntegration struct {
	http *httpclient.CustomHttpClient
}

func NewHttpBinIntegration() HttpbinIntegration {
	return &httpbinIntegration{
		http: httpclient.NewHttpClient(),
	}
}

func (i *httpbinIntegration) GetErrorStatus(ctx context.Context) (body any, err error) {
	ctx, span := tracer.Start(ctx)
	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"body": body,
		}).End(err)
	}()

	http, err := i.http.WithBreaker()
	if err != nil {
		return nil, err
	}

	res, err := http.Request.
		Get(ctx, fmt.Sprintf("%s/status/500", url))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &body)
	if err != nil {
		return
	}

	return
}

func (i *httpbinIntegration) GetSuccessStatus(ctx context.Context) (body any, err error) {
	ctx, span := tracer.Start(ctx)
	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"body": body,
		}).End(err)
	}()

	http, err := i.http.WithBreaker()
	if err != nil {
		return nil, err
	}

	res, err := http.Request.
		Get(ctx, fmt.Sprintf("%s/headers", url))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &body)
	if err != nil {
		return
	}

	return
}
