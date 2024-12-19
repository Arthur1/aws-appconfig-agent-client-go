# AWS AppConfig Agent Client for Go

Go client for AWS AppConfig Agent.

## Requirements

Go 1.22 or higher is required.

We support latest two major releases of Go.

## Usage

### Create a client

```go
app := "app" // Application name or ID of AppConfig
env := "env" // Environment name or ID of AppConfig
client := appconfigagent.NewClient(app, env)
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
