### VARS ###

variable "slack_app_config_token" {
  type        = string
  description = "Slack App Configuration token for use with Slack API"
}

variable "slack_app_id" {
  type        = string
  description = "Slack App ID to lookup"
}
