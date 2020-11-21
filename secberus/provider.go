package secberus

import (
	"context"
	"errors"

	"github.com/RexBelli/go-secberus/secberus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SECBERUS_API_KEY", nil),
				Description: "API key with which to authenticate",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("SECBERUS_USERNAME", nil),
				Description: "username with which to authenticate",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SECBERUS_PASSWORD", nil),
				Description: "password with which to authenticate",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"secberus_organization":              resourceOrganization(),
			"secberus_organization_user_binding": resourceOrganizationUserBinding(),
			"secberus_rule":                      resourceRule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"secberus_compliances":   dataSourceCompliances(),
			"secberus_organizations": dataSourceOrganizations(),
			"secberus_resources":     dataSourceResources(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	apikey, hasapikey := d.Get("api_key").(string)
	username, hasUsername := d.Get("username").(string)
	password, hasPassword := d.Get("password").(string)

	if !hasapikey && (!hasUsername || !hasPassword) {
		return nil, diag.FromErr(errors.New("no api_key or username/password provided"))
	}
	if apikey == "" && (username == "" || password == "") {
		return nil, diag.FromErr(errors.New("no api_key or username/password provided"))
	}

	if apikey != "" {
		c, err := secberus.NewClientWithAPIKey(apikey)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return c, diags
	}
	c, err := secberus.NewClientWithCredentials(username, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
