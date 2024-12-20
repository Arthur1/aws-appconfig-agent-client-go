package appconfigagentv2

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var appConfigAgentBaseURL string

func TestMain(m *testing.M) {
	code := testMain(m)
	os.Exit(code)
}

func testMain(m *testing.M) int {
	ctx := context.Background()
	appConfigAgentCtr, err := setupAppConfigAgentTestcontainers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer appConfigAgentCtr.Terminate(ctx) // nolint
	appConfigAgentBaseURL, err = appConfigAgentCtr.Endpoint(ctx, "http")
	if err != nil {
		log.Fatal(err)
	}
	return m.Run()
}

func TestClient_GetConfiguration_WithTestcontainers(t *testing.T) {
	t.Parallel()
	t.Run("get the full text of the configuration", func(t *testing.T) {
		t.Parallel()
		cases := []struct{ app, env, cfg, want string }{
			{app: "app1", env: "env1", cfg: "cfg1", want: "apple"},
			{app: "app1", env: "env1", cfg: "cfg2", want: "banana"},
			{app: "app1", env: "env2", cfg: "cfg1", want: "apples"},
			{app: "app2", env: "env1", cfg: "cfg1", want: "Alice"},
		}
		for _, tt := range cases {
			t.Run(fmt.Sprintf("%s:%s:%s should be %s", tt.app, tt.env, tt.cfg, tt.want), func(t *testing.T) {
				t.Parallel()
				client, err := NewClient(tt.app, tt.env, WithBaseURL(appConfigAgentBaseURL))
				require.NoError(t, err)
				ctx := context.Background()
				res, err := client.GetConfiguration(ctx, tt.cfg)
				require.NoError(t, err)
				b, err := io.ReadAll(res.ConfigurationBody)
				require.NoError(t, err)
				assert.Equal(t, tt.want, string(b))
			})
		}
	})

	t.Run("return error when resource is not found", func(t *testing.T) {
		t.Parallel()
		cases := []struct{ app, env, cfg string }{
			{app: "app1000000", env: "env1", cfg: "cfg1"},
			{app: "app1", env: "env1000000", cfg: "cfg1"},
			{app: "app1", env: "env1", cfg: "cfg1000000"},
		}
		for _, tt := range cases {
			t.Run(fmt.Sprintf("%s:%s:%s should be not found", tt.app, tt.env, tt.cfg), func(t *testing.T) {
				t.Parallel()
				client, err := NewClient(tt.app, tt.env, WithBaseURL(appConfigAgentBaseURL))
				require.NoError(t, err)
				ctx := context.Background()
				_, err = client.GetConfiguration(ctx, "cfg100000")
				require.Error(t, err)
			})
		}
	})
}

func TestClient_EvaluateFeatureFlag_WithTestcontainers(t *testing.T) {
	t.Parallel()
	t.Run("evaluate a feature flag", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			inFlag string
			want   *EvaluateFeatureFlagResult
		}{
			{
				inFlag: "feature1",
				want: &EvaluateFeatureFlagResult{
					Evaluation: &FeatureFlagEvaluation{
						Enabled:    false,
						Attributes: map[string]any{},
					},
				},
			},
			{
				inFlag: "feature2",
				want: &EvaluateFeatureFlagResult{
					Evaluation: &FeatureFlagEvaluation{
						Enabled: true,
						Attributes: map[string]any{
							"max_items": float64(200),
							"campaign":  "birthday",
						},
					},
				},
			},
			{
				inFlag: "feature3",
				want: &EvaluateFeatureFlagResult{
					Evaluation: &FeatureFlagEvaluation{
						Enabled:    false,
						Variant:    "default",
						Attributes: map[string]any{},
					},
				},
			},
		}
		for _, tt := range cases {
			t.Run(tt.inFlag, func(t *testing.T) {
				t.Parallel()
				client, err := NewClient("app3", "env1", WithBaseURL(appConfigAgentBaseURL))
				require.NoError(t, err)
				ctx := context.Background()
				res, err := client.EvaluateFeatureFlag(ctx, "cfgfeatureflags", tt.inFlag, nil)
				require.NoError(t, err)
				assert.Equal(t, tt.want, res)
			})
		}
	})

	t.Run("return error when resource is not found", func(t *testing.T) {
		t.Parallel()
		cases := []struct{ app, env, cfg, flag string }{
			{app: "app1000000", env: "env1", cfg: "cfgfeatureflags", flag: "feature1"},
			{app: "app3", env: "env000000", cfg: "cfgfeatureflags", flag: "feature1"},
			{app: "app3", env: "env1", cfg: "cfgfeatureflags1000000", flag: "feature1"},
			{app: "app3", env: "env1", cfg: "cfgfeatureflags", flag: "feature1000000"},
		}
		for _, tt := range cases {
			t.Run(fmt.Sprintf("%s:%s:%s:%s should be not found", tt.app, tt.env, tt.cfg, tt.flag), func(t *testing.T) {
				t.Parallel()
				client, err := NewClient(tt.app, tt.env, WithBaseURL(appConfigAgentBaseURL))
				require.NoError(t, err)
				ctx := context.Background()
				_, err = client.EvaluateFeatureFlag(ctx, tt.cfg, tt.flag, nil)
				require.Error(t, err)
			})
		}
	})
}

func TestClient_BulkEvaluateFeatureFlag_WithTestcontainers(t *testing.T) {
	t.Parallel()
	t.Run("evaluate specified multiple feature flags", func(t *testing.T) {
		t.Parallel()
		client, err := NewClient("app3", "env1", WithBaseURL(appConfigAgentBaseURL))
		require.NoError(t, err)
		ctx := context.Background()
		res, err := client.BulkEvaluateFeatureFlag(ctx, "cfgfeatureflags", []string{"feature1", "feature2"}, nil)
		require.NoError(t, err)
		want := &BulkEvaluateFeatureFlagResult{
			Evaluations: map[string]*FeatureFlagEvaluation{
				"feature1": {
					Enabled:    false,
					Attributes: map[string]any{},
				},
				"feature2": {
					Enabled: true,
					Attributes: map[string]any{
						"max_items": float64(200),
						"campaign":  "birthday",
					},
				},
			},
		}
		assert.Equal(t, want, res)
	})

	t.Run("evaluate all feature flags if specified flag keys are empty", func(t *testing.T) {
		t.Parallel()
		client, err := NewClient("app3", "env1", WithBaseURL(appConfigAgentBaseURL))
		require.NoError(t, err)
		ctx := context.Background()
		res, err := client.BulkEvaluateFeatureFlag(ctx, "cfgfeatureflags", nil, nil)
		require.NoError(t, err)
		want := &BulkEvaluateFeatureFlagResult{
			Evaluations: map[string]*FeatureFlagEvaluation{
				"feature1": {
					Enabled:    false,
					Attributes: map[string]any{},
				},
				"feature2": {
					Enabled: true,
					Attributes: map[string]any{
						"max_items": float64(200),
						"campaign":  "birthday",
					},
				},
				"feature3": {
					Enabled:    false,
					Variant:    "default",
					Attributes: map[string]any{},
				},
			},
		}
		assert.Equal(t, want, res)
	})

	t.Run("return error when resource is not found", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			app, env, cfg string
			flags         []string
		}{
			{app: "app1000000", env: "env1", cfg: "cfgfeatureflags", flags: []string{"feature1"}},
			{app: "app3", env: "env000000", cfg: "cfgfeatureflags", flags: []string{"feature1"}},
			{app: "app3", env: "env1", cfg: "cfgfeatureflags1000000", flags: []string{"feature1"}},
			{app: "app3", env: "env1", cfg: "cfgfeatureflags", flags: []string{"feature1", "feature4"}},
		}
		for _, tt := range cases {
			t.Run(fmt.Sprintf("%s:%s:%s:%v should be not found", tt.app, tt.env, tt.cfg, tt.flags), func(t *testing.T) {
				t.Parallel()
				client, err := NewClient(tt.app, tt.env, WithBaseURL(appConfigAgentBaseURL))
				require.NoError(t, err)
				ctx := context.Background()
				_, err = client.BulkEvaluateFeatureFlag(ctx, tt.cfg, tt.flags, nil)
				require.Error(t, err)
			})
		}
	})
}
