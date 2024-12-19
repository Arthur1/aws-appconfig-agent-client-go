resource "aws_appconfig_application" "this" {
  name = local.application_name
}

resource "aws_appconfig_environment" "this" {
  name           = local.environment_name
  application_id = aws_appconfig_application.this.id
}

resource "aws_appconfig_deployment_strategy" "immediately" {
  name                           = "demo-immediately"
  deployment_duration_in_minutes = 0
  final_bake_time_in_minutes     = 0
  growth_factor                  = 100
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
}

// conf1
resource "aws_appconfig_configuration_profile" "conf1" {
  application_id = aws_appconfig_application.this.id
  location_uri   = "hosted"
  name           = local.configuration1_name
  type           = "AWS.Freeform"
}

resource "aws_appconfig_hosted_configuration_version" "conf1" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf1.configuration_profile_id
  content_type             = "application/json"
  content = jsonencode({
    hello = "world"
  })
}

resource "aws_appconfig_deployment" "conf1" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf1.configuration_profile_id
  configuration_version    = aws_appconfig_hosted_configuration_version.conf1.version_number
  deployment_strategy_id   = aws_appconfig_deployment_strategy.immediately.id
  environment_id           = aws_appconfig_environment.this.environment_id
}

// conf2
resource "aws_appconfig_configuration_profile" "conf2" {
  application_id = aws_appconfig_application.this.id
  location_uri   = "hosted"
  name           = local.configuration2_name
  type           = "AWS.Freeform"
}

resource "aws_appconfig_hosted_configuration_version" "conf2" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf2.configuration_profile_id
  content_type             = "text/plain"
  content                  = "Hello, world!"
}

resource "aws_appconfig_deployment" "conf2" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf2.configuration_profile_id
  configuration_version    = aws_appconfig_hosted_configuration_version.conf2.version_number
  deployment_strategy_id   = aws_appconfig_deployment_strategy.immediately.id
  environment_id           = aws_appconfig_environment.this.environment_id
}

// conf_featureflag
resource "aws_appconfig_configuration_profile" "conf_featureflag" {
  application_id = aws_appconfig_application.this.id
  location_uri   = "hosted"
  name           = local.configuration_featureflag_name
  type           = "AWS.AppConfig.FeatureFlags"
}

resource "aws_appconfig_hosted_configuration_version" "conf_featureflag" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf_featureflag.configuration_profile_id
  content_type             = "application/json"
  content = jsonencode({
    version = "1"
    flags = {
      feature1 = {
        name       = "feature 1"
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
      }
      feature2 = {
        name       = "feature 2"
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
        attributes = {
          max_items = {
            constraints = {
              type     = "number"
              required = true
            }
          }
        }
      }
      feature3 = {
        name       = "feature 3"
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
      }
    }
    values = {
      feature1 = {
        enabled    = false
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
      }
      feature2 = {
        enabled    = true
        max_items  = 200
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
      }
      feature3 = {
        _variants = [
          {
            name    = "users"
            rule    = <<-EOT
            (or
              (eq $userId "userB")
              (eq $userId "userC")
            )
            EOT
            enabled = true
          },

          {
            name    = "default"
            enabled = false
          },
        ]
        _createdAt = "2024-08-28T10:00:00.000Z"
        _updatedAt = "2024-08-28T10:00:00.000Z"
      }
    }
  })
}

resource "aws_appconfig_deployment" "conf_featureflag" {
  application_id           = aws_appconfig_application.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.conf_featureflag.configuration_profile_id
  configuration_version    = aws_appconfig_hosted_configuration_version.conf_featureflag.version_number
  deployment_strategy_id   = aws_appconfig_deployment_strategy.immediately.id
  environment_id           = aws_appconfig_environment.this.environment_id
}
