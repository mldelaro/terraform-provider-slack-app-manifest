package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mldelaro/slack"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("SLACK_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"slack_manifest": dataSourceManifest(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"slack_manifest": resourceManifest(),
			},
		}

		p.ConfigureContextFunc = providerConfigure

		return p
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	slackApiToken := d.Get("token").(string)
	var diags diag.Diagnostics

	if slackApiToken != "" {
		c :=  slack.New(slackApiToken)
		return c, diags
	}

	c :=  slack.New("")
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Unable to create Slack API Client",
		Detail:   "Unable to auth user for authenticated Slack API client",
	})
	return c, diags
}
