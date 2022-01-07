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
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	appId := d.Get("app_id").(string)
	api := slack.New("TOKEN_HERE")
	manifest, err := api.ExportAppManifest(appId) //("A02TDSWCDDE")
	if err != nil {
		return diag.Errorf("Failed to make request via client.")
	}

	d.Set("display_name", manifest.DisplayInformation.Name)

	metadata := flattenMetadata(manifest.Metadata)
	if err := d.Set("_metadata", metadata); err != nil {
		return diag.FromErr(err)
	}

	features := flattenFeatures(manifest.Features)
	if err := d.Set("features", features); err != nil {
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

func flattenFeatures(feature *slack.Features) []interface{} {
	fs := make([]interface{}, 1, 1)

	f := make(map[string]interface{})
	bu := flattenBotUser(feature.BotUser)
	f["bot_user"] = bu

	sc := flattenSlashCommands(&feature.SlashCommands)
	f["slash_commands"] = sc
	fs[0] = f

	return fs
}

func flattenBotUser(botUser *slack.BotUser) []interface{} {
	bs := make([]interface{}, 1, 1)

	b := make(map[string]interface{})
	b["display_name"] = botUser.DisplayName
	b["always_online"] = botUser.AlwaysOnline
	bs[0] = b

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
