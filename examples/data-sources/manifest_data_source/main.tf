terraform {
  required_providers {
    slack = {
      version = "0.3.1"
      source  = "asu.edu/mldelaro/slack-app-manifest"
    }
  }
}

provider "slack" {
  token = var.slack_app_config_token
}

data "slack_manifest" "example" {
  app_id = var.slack_app_id
}

output "some_manifest" {
  value = data.slack_manifest.example
}
