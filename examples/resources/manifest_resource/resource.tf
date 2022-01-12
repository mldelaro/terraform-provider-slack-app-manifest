terraform {
  required_providers {
    slack = {
      version = "0.3.1"
      source  = "hashicorp.com/edu/slack-app-manifest"
    }
  }
}

provider "slack" {
  token = var.slack_app_config_token
}

/*
resource "slack_manifest" "example" {
  manifest = file("some-bot-app-manifest.json")
}
*/
