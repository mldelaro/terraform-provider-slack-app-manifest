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

resource "slack_manifest" "example" {
  sample_attribute = "foo"
}