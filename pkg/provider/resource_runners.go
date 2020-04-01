package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// we set the id here so terraform knows the resource has been
// "created" and will not attempt to destroy it because of an empty id
// we don't actually do anything here - just tell terraform we
// are in "created" state and everything is ok
func resourceRunnersCreate(d *schema.ResourceData, m interface{}) error {
	//d.SetId("dotscience-runners")
	return resourceRunnersRead(d, m)
}

// these 2 are no-ops, we don't actually want terraform to know
// about or update our list of runners
func resourceRunnersRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

// this is the key handler - it's job is to delete any runners
// via the gateway api and only return when they are deleted
func resourceRunnersDelete(d *schema.ResourceData, m interface{}) error {
	// client := (m).Client
	// // ping the API for version to check our credentials
	// if _, err := client.Version(); err != nil {
	// 	return fmt.Errorf("Error connecting to the dotscience API: %s", err)
	// }
	// // load a list of the runners
	// runners, err := client.ListRunners()
	// if err != nil {
	// 	return fmt.Errorf("Error loading runners: %s", err)
	// }

	// if logging.IsDebugOrHigher() {
	// 	log.Printf("[DEBUG] found runner list count: %d", len(runners))
	// }

	return nil
}

func resourceRunners() *schema.Resource {
	return &schema.Resource{
		Create: resourceRunnersCreate,
		Read:   resourceRunnersRead,
		Delete: resourceRunnersDelete,

		Schema: map[string]*schema.Schema{
			// "address": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	ForceNew: true,
			// },
		},
	}
}
