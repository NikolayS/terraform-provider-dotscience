package provider

import (
	"github.com/dotmesh-io/terraform-provider-dotscience/pkg/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hub_public_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_URL", nil),
				Description: "The public url of your dotscience hub.",
			},
			"hub_admin_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_USERNAME", "admin"),
				Description: "The username for the admin user.",
			},
			"hub_admin_password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOTSCIENCE_PASSWORD", nil),
				Description: "The password for the admin user.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"dotscience_runners": resourceRunners(),
		},
	}
	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		client := &api.Client{
			URL:      d.Get("hub_public_url").(string),
			Username: d.Get("hub_admin_username").(string),
			Password: d.Get("hub_admin_password").(string),
		}
		return client, nil
	}
	return provider
}
