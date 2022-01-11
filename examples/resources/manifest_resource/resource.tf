terraform {
  required_providers {
    slack = {
      version = "0.3.1"
      source = "hashicorp.com/edu/slack-app-manifest"
    }
  }
}

provider "slack" {
  token = "APP_TOKEN_HERE"
}

resource "slack_manifest" "example" {
  manifest = "{\"_metadata\":{\"major_version\":1,\"minor_version\":1},\"display_information\":{\"name\":\"SomeSlackApp\"},\"features\":{\"bot_user\":{\"display_name\":\"test-manifest\",\"always_online\":false},\"slash_commands\":[{\"command\":\"/servicenow-links-dev\",\"url\":\"https://some-lambda.execute-api.us-west-2.amazonaws.com/event/slash-command\",\"description\":\"SomeDescription\",\"usage_hint\":\"[help]\",\"should_escape\":false}]},\"oauth_config\":{\"redirect_urls\":[\"https://some-lambda.execute-api.us-west-2.amazonaws.com/oauth\"],\"scopes\":{\"bot\":[\"chat:write\",\"commands\",\"im:read\",\"app_mentions:read\",\"channels:history\",\"groups:history\",\"im:history\",\"mpim:history\"]}},\"settings\":{\"event_subscriptions\":{\"request_url\":\"https://some-lambda.execute-api.us-west-2.amazonaws.com/event\",\"bot_events\":[\"message.channels\",\"message.groups\",\"message.im\",\"message.mpim\"]},\"org_deploy_enabled\":false,\"socket_mode_enabled\":false,\"token_rotation_enabled\":false}}"
}