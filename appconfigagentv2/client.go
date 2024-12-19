package appconfigagentv2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Arthur1/aws-appconfig-agent-client-go/internal/apiv2"
)

// Client for AWS AppConfig Agent
type Client struct {
	apiClient   *apiv2.Client
	application string
	environment string
}

// NewClient initialize new Client of AppConfig Agent.
//
// application argument is the name or ID of AWS AppConfig Application.
// environment argument is the name or ID of AWS AppConfig Environment.
func NewClient(application, environment string, opts ...ClientOption) (*Client, error) {
	options := clientOptions{
		baseURL:    "http://localhost:2772",
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}

	apiClient, err := apiv2.NewClient(
		options.baseURL,
		apiv2.WithClient(options.httpClient),
		apiv2.WithTracerProvider(options.tracerProvider),
		apiv2.WithMeterProvider(options.meterProvider),
	)
	if err != nil {
		return nil, err
	}
	return &Client{
		apiClient:   apiClient,
		application: application,
		environment: environment,
	}, nil
}

// Result of GetConfiguration
type GetConfigurationResult struct {
	ConfigurationBody io.Reader
}

// GetConfiguration gets a configuration value.
//
// configuration argument is a name or ID of AWS AppConfig Configuration.
func (c *Client) GetConfiguration(ctx context.Context, configuration string) (*GetConfigurationResult, error) {
	res, err := c.apiClient.GetConfiguration(ctx, apiv2.GetConfigurationParams{
		Application:   c.application,
		Environment:   c.environment,
		Configuration: configuration,
	})
	if err != nil {
		return nil, err
	}
	body, _, err := getBodyAndConfigurationVersionFromResponse(res)
	if err != nil {
		return nil, err
	}
	return &GetConfigurationResult{
		ConfigurationBody: body,
	}, nil
}

// Evaluated feacture flag value
type FeatureFlagEvaluation struct {
	// Feature flag status
	Enabled bool
	// Feature flag attributes
	Attributes map[string]any
	// Feature flag variant in multi-variant feature flags
	Variant string
}

// Result of EvaluateFeatureFlagResult
type EvaluateFeatureFlagResult struct {
	Evaluation *FeatureFlagEvaluation
}

// EvaluateFeatureFlag evaluates a single flag in feature flag configuration.
//
// configuration argument is a name or ID of AWS AppConfig Configuration.
// evalCtx argument is a context value for evaluating multi-variant flags.
func (c *Client) EvaluateFeatureFlag(ctx context.Context, configuration string, flagKey string, evalCtx map[string]any) (*EvaluateFeatureFlagResult, error) {
	evalCtxHeaders := make([]string, 0, len(evalCtx))
	for k, v := range evalCtx {
		vs, err := evalCtxValueToString(v)
		if err != nil {
			continue
		}
		evalCtxHeaders = append(evalCtxHeaders, fmt.Sprintf("%s=%s", k, vs))
	}

	res, err := c.apiClient.GetConfiguration(ctx, apiv2.GetConfigurationParams{
		Application:   c.application,
		Environment:   c.environment,
		Configuration: configuration,
		Flag:          []string{flagKey},
		Context:       evalCtxHeaders,
	})
	if err != nil {
		return nil, err
	}
	body, _, err := getBodyAndConfigurationVersionFromResponse(res)
	if err != nil {
		return nil, err
	}

	result := &EvaluateFeatureFlagResult{}

	var jsonm map[string]any
	if err = json.NewDecoder(body).Decode(&jsonm); err != nil {
		return nil, err
	}
	result.Evaluation, err = parseFeatureFlagEvaluation(jsonm)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Result of BulkEvaluateFeatureFlag
//
// Evaluations is a map whose key is the flag key.
type BulkEvaluateFeatureFlagResult struct {
	Evaluations map[string]*FeatureFlagEvaluation
}

// BulkEvaluateFeatureFlag evaluates multiple flags in feature flag configuration.
//
// configuration argument is the name or ID of AWS AppConfig Configuration.
// evalCtx argument is the context value for evaluating multi-variant flags.
func (c *Client) BulkEvaluateFeatureFlag(ctx context.Context, configuration string, flagKeys []string, evalCtx map[string]any) (*BulkEvaluateFeatureFlagResult, error) {
	evalCtxHeaders := make([]string, 0, len(evalCtx))
	for k, v := range evalCtx {
		vs, err := evalCtxValueToString(v)
		if err != nil {
			continue
		}
		evalCtxHeaders = append(evalCtxHeaders, fmt.Sprintf("%s=%s", k, vs))
	}

	res, err := c.apiClient.GetConfiguration(ctx, apiv2.GetConfigurationParams{
		Application:   c.application,
		Environment:   c.environment,
		Configuration: configuration,
		Flag:          flagKeys,
		Context:       evalCtxHeaders,
	})
	if err != nil {
		return nil, err
	}
	body, _, err := getBodyAndConfigurationVersionFromResponse(res)
	if err != nil {
		return nil, err
	}

	var jsonm map[string]map[string]any
	if err = json.NewDecoder(body).Decode(&jsonm); err != nil {
		return nil, err
	}

	evaluations := make(map[string]*FeatureFlagEvaluation, len(jsonm))
	for flagk, flagev := range jsonm {
		evaluation, err := parseFeatureFlagEvaluation(flagev)
		if err != nil {
			return nil, err
		}
		evaluations[flagk] = evaluation
	}

	result := &BulkEvaluateFeatureFlagResult{
		Evaluations: evaluations,
	}
	return result, nil
}

func parseFeatureFlagEvaluation(m map[string]any) (*FeatureFlagEvaluation, error) {
	evaluation := &FeatureFlagEvaluation{
		Attributes: map[string]any{},
	}
	for k, v := range m {
		switch k {
		case "enabled":
			enabled, ok := v.(bool)
			if !ok {
				return nil, fmt.Errorf("enabled must be bool")
			}
			evaluation.Enabled = enabled
		case "_variant":
			variant, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("_variant must be string")
			}
			evaluation.Variant = variant
		default:
			evaluation.Attributes[k] = v
		}
	}
	return evaluation, nil
}

func getBodyAndConfigurationVersionFromResponse(res apiv2.GetConfigurationRes) (body io.Reader, version string, err error) {
	switch tres := res.(type) {
	case *apiv2.GetConfigurationOKHeaders:
		body = tres.GetResponse().Data
		version, _ = tres.GetConfigurationVersion().Get()
		return
	case *apiv2.GetConfigurationBadRequest:
		b, _ := tres.MarshalJSON()
		err = fmt.Errorf("BadRequestException: %s", string(b))
		return
	case *apiv2.GetConfigurationNotFound:
		b, _ := tres.MarshalJSON()
		err = fmt.Errorf("ResourceNotFoundException: %s", string(b))
		return
	case *apiv2.GetConfigurationInternalServerErrorHeaders:
		err = fmt.Errorf("InternalServerException")
		return
	case *apiv2.GetConfigurationBadGatewayHeaders:
		err = fmt.Errorf("BadGatewayException")
		return
	case *apiv2.GetConfigurationGatewayTimeoutHeaders:
		err = fmt.Errorf("GatewayTimeoutException")
		return
	}
	err = fmt.Errorf("UnexpectedResponseStatusException")
	return
}

func evalCtxValueToString(v any) (string, error) {
	switch vt := v.(type) {
	case string:
		return vt, nil
	case bool:
		return strconv.FormatBool(vt), nil
	case int64:
		return strconv.FormatInt(vt, 10), nil
	case float64:
		return strconv.FormatFloat(vt, 'E', -1, 64), nil
	case time.Time:
		return vt.String(), nil
	default:
		j, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(j), nil
	}
}
