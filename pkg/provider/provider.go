package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hub_public_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_PUBLIC_URL", nil),
				Description: "The public url of your dotscience hub.",
			},
			"hub_admin_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_ADMIN_USERNAME", "admin"),
				Description: "The username for the admin user.",
			},
			"hub_admin_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_ADMIN_PASSWORD", nil),
				Description: "The password for the admin user.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"dotscience_runners": resourceRunners(),
		},
	}
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}
	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	client := &Client{
		URL:      d.Get("hub_public_url").(string),
		Username: d.Get("hub_admin_username").(string),
		Password: d.Get("hub_admin_password").(string),
	}
	return client, nil
}
