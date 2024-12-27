# AWS AppConfig Agent Client for Go

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Arthur1/aws-appconfig-agent-client-go) [![Go Reference](https://pkg.go.dev/badge/github.com/Arthur1/aws-appconfig-agent-client-go.svg)](https://pkg.go.dev/github.com/Arthur1/aws-appconfig-agent-client-go) [![CI for main branch](https://github.com/Arthur1/aws-appconfig-agent-client-go/actions/workflows/on-push-main.yml/badge.svg)](https://github.com/Arthur1/aws-appconfig-agent-client-go/actions/workflows/on-push-main.yml)

Go client for AWS AppConfig Agent.

## Requirements

Go 1.22 or higher is required. We support latest two major releases of Go.

Currently, the major version of AppConfig Agent must be v2.

## Usage

### Create a client

```go
app := "app" // Application name or ID of AppConfig
env := "env" // Environment name or ID of AppConfig
client := appconfigagentv2.NewClient(app, env)
```

### Get configuration data

```go
ctx := context.Background()
cfg := "cfg" // Configuration name or ID of AppConfig
res, err := client.GetConfiguration(ctx, cfg)
if err != nil {
    log.Fatal(err)
}
data, err := io.ReadAll(res.ConfigurationBody)
```

### Evaluate a flag in feature flag configuration

This is only available for configurations of feature flags type.

```go
ctx := context.Background()
cfg := "cfg" // Configuration name or ID of AppConfig
flag := "feature1" // Flag key of AppConfig feature flags
res, err := client.EvaluateFeatureFlag(ctx, cfg, flag, nil)
```

You can also pass evaluation contexts for multi-variant flags:

```go
ctx := context.Background()
cfg := "cfg" // Configuration name or ID of AppConfig
flag := "feature1" // Flag key of AppConfig feature flags
evalCtx := map[string]any{"orgID": "awesomeOrg", "userID": 123}
res, err := client.EvaluateFeatureFlag(ctx, cfg, flag, evalCtx)
```

### Evaluate multiple flags in feature flag configuration

This is only available for configurations of feature flags type.

```go
ctx := context.Background()
cfg := "cfg" // Configuration name or ID of AppConfig
flag1 := "feature1" // Flag key of AppConfig feature flags
flag2 := "feature2"
res, err := client.BulkEvaluateFeatureFlag(ctx, cfg, []string{flag1, flag2}, nil)
if err != nil {
    log.Fatal(err)
}
flag1res, ok := res.Evaluations[flag1]
flag2res, ok := res.Evaluations[flag2]
```

## License

MIT License

## Contact

Please contact me in GitHub issues or [`@Arthur1__` on X](https://x.com/arthur1__).
