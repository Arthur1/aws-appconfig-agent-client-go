package main

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/Arthur1/aws-appconfig-agent-client-go/appconfigagent"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env/v11"
)

type config struct {
	AppConfigApp            string `env:"APPCONFIG_APPLICATION"`
	AppConfigEnv            string `env:"APPCONFIG_ENVIRONMENT"`
	AppConfigCfg1           string `env:"APPCONFIG_CONFIGURATION1"`
	AppConfigCfg2           string `env:"APPCONFIG_CONFIGURATION2"`
	AppConfigCfgFeatureFlag string `env:"APPCONFIG_CONFIGURATION_FEATUREFLAG"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	client, err := appconfigagent.NewClient(cfg.AppConfigApp, cfg.AppConfigEnv)
	if err != nil {
		log.Fatal(err)
	}

	res1, err := client.GetConfiguration(ctx, cfg.AppConfigCfg1)
	if err != nil {
		log.Fatal(err)
	}

	res1b, err := io.ReadAll(res1.ConfigurationBody)
	if err != nil {
		log.Fatal(err)
	}
	type Content struct {
		Hello string
	}
	content := Content{}
	if err = json.Unmarshal(res1b, &content); err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration1 value: %+v\n", content)

	res2, err := client.GetConfiguration(ctx, cfg.AppConfigCfg2)
	if err != nil {
		log.Fatal(err)
	}
	res2b, err := io.ReadAll(res2.ConfigurationBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("configuration2 value: %+v\n", string(res2b))

	evalCtxA := map[string]any{"userId": "userA"}
	evalCtxB := map[string]any{"userId": "userB"}
	resf1, err := client.EvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, "feature1", evalCtxA)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resf1 enabled: expected=false, actual=%t\n", resf1.Evaluation.Enabled)
	log.Printf("resf1 variant: expected=, actual=%s\n", resf1.Evaluation.Variant)
	log.Printf("resf1 attributes: expected=map[], actual=%+v\n", resf1.Evaluation.Attributes)

	resf2, err := client.EvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, "feature2", evalCtxA)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resf2 enabled: expected=true, actual=%t\n", resf2.Evaluation.Enabled)
	log.Printf("resf2 variant: expected=, actual=%s\n", resf2.Evaluation.Variant)
	log.Printf("resf2 attributes: expected=map[max_items:200], actual=%+v\n", resf2.Evaluation.Attributes)

	resf3a, err := client.EvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, "feature3", evalCtxA)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resf3a enabled: expected=false, actual=%t\n", resf3a.Evaluation.Enabled)
	log.Printf("resf3a variant: expected=default, actual=%s\n", resf3a.Evaluation.Variant)
	log.Printf("resf3a attributes: expected=map[], actual=%+v\n", resf3a.Evaluation.Attributes)

	resf3b, err := client.EvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, "feature3", evalCtxB)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resf3b enabled: expected=true, actual=%t\n", resf3b.Evaluation.Enabled)
	log.Printf("resf3b variant: expected=users, actual=%s\n", resf3b.Evaluation.Variant)
	log.Printf("resf3b attributes: expected=map[], actual=%+v\n", resf3b.Evaluation.Attributes)

	resfalla, err := client.BulkEvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, []string{"feature1", "feature2", "feature3"}, evalCtxA)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resfalla feature1: %+v", resfalla.Evaluations["feature1"])
	log.Printf("resfalla feature2: %+v", resfalla.Evaluations["feature2"])
	log.Printf("resfalla feature3: %+v", resfalla.Evaluations["feature3"])

	resfallb, err := client.BulkEvaluateFeatureFlag(ctx, cfg.AppConfigCfgFeatureFlag, []string{"feature1", "feature2", "feature3"}, evalCtxB)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("resfallb feature1: %+v", resfallb.Evaluations["feature1"])
	log.Printf("resfallb feature2: %+v", resfallb.Evaluations["feature2"])
	log.Printf("resfallb feature3: %+v", resfallb.Evaluations["feature3"])
}
