package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mldelaro/slack"
)

func dataSourceManifest() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Reads App Manifest via Slack API's app.manifest.export path.",

		ReadContext: dataSourceManifestRead,

		Schema: map[string]*schema.Schema{
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"_metadata": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"major_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minor_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"display_information": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"features": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bot_user": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"display_name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"always_online": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"slash_commands": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"command": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"usage_hint": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"should_escape": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"oauth_config": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_urls": &schema.Schema{
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"scopes": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bot": &schema.Schema{
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"settings": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_subscriptions": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"request_url": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"bot_events": &schema.Schema{
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"org_deploy_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"socket_mode_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"token_rotation_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},

					},
				},
			},
		},
	}
}

func dataSourceManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	slackApiClient := meta.(*slack.Client)
	appId := d.Get("app_id").(string)
	manifest, err := slackApiClient.ExportAppManifest(appId)
	if err != nil {
		return diag.Errorf("Failed to make request via client.")
	}

	d.Set("display_name", manifest.DisplayInformation.Name)

	metadata := flattenMetadata(manifest.Metadata)
	if err := d.Set("_metadata", metadata); err != nil {
		return diag.FromErr(err)
	}

	displayInformation := flattenDisplayInformation(manifest.DisplayInformation)
	if err := d.Set("display_information", displayInformation); err != nil {
		return diag.FromErr(err)
	}

	features := flattenFeatures(manifest.Features)
	if err := d.Set("features", features); err != nil {
		return diag.FromErr(err)
	}

	oauthConfig := flattenOAuthConfig(manifest.OAuthConfig)
	if err := d.Set("oauth_config", oauthConfig); err != nil {
		return diag.FromErr(err)
	}

	settings := flattenSettings(manifest.Settings)
	if err := d.Set("settings", settings); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(appId)

	return diags
}

func flattenMetadata(metadata *slack.Metadata) []interface{} {
	mds := make([]interface{}, 1, 1)

	md := make(map[string]interface{})
	md["major_version"] = metadata.Majorversion
	md["minor_version"] = metadata.Minorversion
	mds[0] = md

	return mds
}

func flattenDisplayInformation(displayInformation *slack.DisplayInformation) []interface{} {
	ds := make([]interface{}, 1, 1)

	d := make(map[string]interface{})
	d["name"] = displayInformation.Name
	ds[0] = d

	return ds
}

func flattenFeatures(feature *slack.Features) []interface{} {
	fs := make([]interface{}, 1, 1)

	if feature != nil {
		f := make(map[string]interface{})
		bu := flattenBotUser(feature.BotUser)
		f["bot_user"] = bu

		sc := flattenSlashCommands(&feature.SlashCommands)
		f["slash_commands"] = sc
		fs[0] = f

		return fs
	}

	return make([]interface{}, 0)
}

func flattenBotUser(botUser *slack.BotUser) []interface{} {
	bs := make([]interface{}, 1, 1)

	if botUser != nil {
		b := make(map[string]interface{})
		b["display_name"] = botUser.DisplayName
		b["always_online"] = botUser.AlwaysOnline
		bs[0] = b
	}

	return bs
}

func flattenSlashCommands(slashCommands *[]slack.SlashCommandManifest) []interface{} {
	if slashCommands != nil {
		scs := make([]interface{}, len(*slashCommands), len(*slashCommands))

		for i, slashCommand := range *slashCommands {
			sc := make(map[string]interface{})

			sc["command"] = slashCommand.Command
			sc["url"] = slashCommand.Url
			sc["description"] = slashCommand.Description
			sc["usage_hint"] = slashCommand.UsageHint
			sc["should_escape"] = slashCommand.ShouldEscape

			scs[i] = sc
		}
		return scs
	}
	return make([]interface{}, 0)
}

func flattenOAuthConfig(oauthConfig *slack.OAuthConfig) []interface{} {
	oacs := make([]interface{}, 1, 1)

	if oauthConfig != nil {
		oac := make(map[string]interface{})
		oac["redirect_urls"] = oauthConfig.RedirectUrls
		oac["scopes"] = flattenScopes(&oauthConfig.Scopes)
		oacs[0] = oac

		return oacs
	}

	return make([]interface{}, 0)
}


func flattenScopes(scopes *slack.Scopes) []interface{} {
	ss := make([]interface{}, 1, 1)

	if scopes != nil {
		s := make(map[string]interface{})
		s["bot"] = scopes.Bot
		ss[0] = s
	}

	return ss
}

func flattenRedirectUrls(slashCommands *[]slack.SlashCommandManifest) []interface{} {
	if slashCommands != nil {
		scs := make([]interface{}, len(*slashCommands), len(*slashCommands))

		for i, slashCommand := range *slashCommands {
			sc := make(map[string]interface{})

			sc["command"] = slashCommand.Command
			sc["url"] = slashCommand.Url
			sc["description"] = slashCommand.Description
			sc["usage_hint"] = slashCommand.UsageHint
			sc["should_escape"] = slashCommand.ShouldEscape

			scs[i] = sc
		}
		return scs
	}
	return make([]interface{}, 0)
}

func flattenSettings(settings *slack.Settings) []interface{} {
	ss := make([]interface{}, 1, 1)

	if settings != nil {
		s := make(map[string]interface{})
		s["event_subscriptions"] = flattenEventSubscriptions(settings.EventSubscriptions)
		s["org_deploy_enabled"] = settings.OrgDeployEnabled
		s["socket_mode_enabled"] = settings.SocketModeEnabled
		s["token_rotation_enabled"] = settings.TokenRotationEnabled
		ss[0] = s

		return ss
	}

	return make([]interface{}, 0)
}

func flattenEventSubscriptions(eventSubscriptions *slack.EventSubscriptions) []interface{} {
	ess := make([]interface{}, 1, 1)

	if eventSubscriptions != nil {
		es := make(map[string]interface{})
		es["request_url"] = eventSubscriptions.RequestUrl
		es["bot_events"] = eventSubscriptions.BotEvents
		ess[0] = es
	}

	return ess
}
