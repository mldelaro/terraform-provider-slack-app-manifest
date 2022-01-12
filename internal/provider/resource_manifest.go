package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mldelaro/slack"
)

func resourceManifest() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Resource that manages App Manifests.",

		CreateContext: resourceManifestCreate,
		ReadContext:   resourceManifestRead,
		UpdateContext: resourceManifestUpdate,
		DeleteContext: resourceManifestDelete,

		Schema: map[string]*schema.Schema{
			"manifest": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"app_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"oauth_authorize_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_secret": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"verification_token": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"signing_secret": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceManifestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	slackApiClient := meta.(*slack.Client)
	appManifest := d.Get("manifest").(string)
	newManifestResponse, err := slackApiClient.CreateAppManifest(appManifest)
	if err != nil {
		return diag.Errorf("Failed to make request to create new manifest via client.")
	}

	d.SetId(newManifestResponse.AppId)
	if err := d.Set("manifest", appManifest); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_id", newManifestResponse.AppId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("oauth_authorize_url", newManifestResponse.OAuthAuthorizeUrl); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("credentials", flattenCredentials(newManifestResponse.Credentials)); err != nil {
		return diag.FromErr(err)
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	resourceManifestRead(ctx, d, meta)

	return diags
}

func flattenCredentials(credentials *slack.Credentials) []interface{} {
	cs := make([]interface{}, 1, 1)

	c := make(map[string]interface{})
	c["client_id"] = credentials.ClientId
	c["client_secret"] = credentials.ClientSecret
	c["verification_token"] = credentials.VerificationToken
	c["signing_secret"] = credentials.SigningSecret
	cs[0] = c

	return cs
}

func resourceManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	slackApiClient := meta.(*slack.Client)
	appId := d.Id()
	manifest, err := slackApiClient.ExportAppManifest(appId)
	if err != nil {
		return diag.Errorf("Failed to make request to read manifest via client.")
	}

	strManifest, err := json.Marshal(manifest)
	if err != nil {
		return diag.Errorf("Failed to marshal app manifest.")
	}

	d.Set("manifest", strManifest)

	return diags
}

func resourceManifestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	slackApiClient := meta.(*slack.Client)
	appManifest := d.Get("manifest").(string)
	appId := d.Get("app_id").(string)
	updateManifestResponse, err := slackApiClient.UpdateAppManifest(appId, appManifest)
	if err != nil {
		return diag.Errorf("Failed to make request to update manifest via client.")
	}

	d.SetId(updateManifestResponse.AppId)
	if err := d.Set("app_id", updateManifestResponse.AppId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("manifest", appManifest); err != nil {
		return diag.FromErr(err)
	}
	if !updateManifestResponse.Ok {
		return diag.Errorf("Received non-ok response from client.")
	}
	resourceManifestRead(ctx, d, meta)
	return diags
}

func resourceManifestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	slackApiClient := meta.(*slack.Client)
	appId := d.Id()
	slackResponse, err := slackApiClient.DeleteAppManifest(appId)
	if err != nil {
		return diag.Errorf("Failed to make request to delete manifest via client.")
	}

	if !slackResponse.Ok {
		return diag.Errorf("Received non-ok response from client.")
	}

	// d.SetId("") is automatically called assuming delete returns no errors
	return diags
}
