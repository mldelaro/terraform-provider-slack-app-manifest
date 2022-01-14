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

resource "slack_manifest" "example" {
  manifest = file("some-bot-app-manifest.json")
}
