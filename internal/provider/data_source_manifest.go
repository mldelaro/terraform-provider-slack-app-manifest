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
