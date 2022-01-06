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
	d.SetId(appId)

	return diags
}
