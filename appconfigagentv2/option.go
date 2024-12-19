package appconfigagentv2

import (
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type clientOptions struct {
	baseURL        string
	httpClient     *http.Client
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
}

// Option for NewClient
type ClientOption interface {
	apply(opts *clientOptions)
}

var (
	_ ClientOption = baseURLOption("")
	_ ClientOption = (*httpClientOption)(nil)
	_ ClientOption = (*tracerProviderOption)(nil)
	_ ClientOption = (*meterProviderOption)(nil)
)

type baseURLOption string

func (o baseURLOption) apply(opts *clientOptions) {
	opts.baseURL = string(o)
}

// WithBaseURL sets a base URL of AppConfig Agent Server.
// Default value is http://localhost:2772.
func WithBaseURL(baseURL string) ClientOption {
	return baseURLOption(baseURL)
}

type httpClientOption struct {
	httpClient *http.Client
}

func (o httpClientOption) apply(opts *clientOptions) {
	opts.httpClient = o.httpClient
}

// WithHTTPClient sets a HTTP client.
// Default value is http.DefaultClient.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return httpClientOption{httpClient}
}

type tracerProviderOption struct {
	tracerProvider trace.TracerProvider
}

func (o tracerProviderOption) apply(opts *clientOptions) {
	opts.tracerProvider = o.tracerProvider
}

// WithTracerProvider sets a tracer provider of OpenTelemetry.
func WithTracerProvider(tracerProvider trace.TracerProvider) ClientOption {
	return tracerProviderOption{tracerProvider}
}

type meterProviderOption struct {
	meterProvider metric.MeterProvider
}

func (o meterProviderOption) apply(opts *clientOptions) {
	opts.meterProvider = o.meterProvider
}

// WithMeterProvider sets a meter provider of OpenTelemetry.
func WithMeterProvider(meterProvider metric.MeterProvider) ClientOption {
	return meterProviderOption{meterProvider}
}
