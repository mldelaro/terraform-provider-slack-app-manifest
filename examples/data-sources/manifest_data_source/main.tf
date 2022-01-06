terraform {
  required_providers {
    slack = {
      version = "0.3.1"
      source = "hashicorp.com/edu/slack-app-manifest"
    }
  }
}

provider "slack" {
  token = "TOKEN_HERE"
}

data "slack_manifest" "example" {
  app_id = "A02TDSWCDDE"
}

output "some_manifest" {
  value = data.slack_manifest.example
}