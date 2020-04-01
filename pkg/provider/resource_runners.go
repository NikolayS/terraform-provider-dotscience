package provider

import (
	"fmt"
	"log"
	"time"

	"github.com/dotmesh-io/terraform-provider-dotscience/pkg/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func deleteAllTasks(client *api.Client) error {
	runners, err := client.ListRunners()
	if err != nil {
		return fmt.Errorf("Error loading runners: %s", err)
	}
	log.Printf("found %d runners, terminating tasks", len(runners))
	for _, runner := range runners {
		log.Printf("terminating tasks on runner: %s", runner.ID)
		err = client.StopRunnerTasks(runner)
		if err != nil {
			return err
		}
	}
	return waitOnFunction("waitOnTaskTermination", time.Minute*5, time.Second*5, func() bool {
		runners, err := client.ListRunners()
		if err != nil {
			return false
		}
		for _, runner := range runners {
			for _, task := range runner.Tasks {
				if task.Status != "terminated" {
					return false
				}
			}
		}
		return true
	})
}

func deleteAllRunners(client *api.Client) error {
	runners, err := client.ListRunners()
	if err != nil {
		return fmt.Errorf("Error loading runners: %s", err)
	}
	log.Printf("found %d runners, deleting", len(runners))
	for _, runner := range runners {
		log.Printf("deleting runner: %s", runner.ID)
		err = client.DeleteRunner(runner)
		if err != nil {
			return err
		}
	}
	return waitOnFunction("waitOnRunnersDeleted", time.Minute*5, time.Second*5, func() bool {
		runners, err := client.ListRunners()
		if err != nil {
			return false
		}
		return len(runners) == 0
	})
}

// we set the id here so terraform knows the resource has been
// "created" and will not attempt to destroy it because of an empty id
// we don't actually do anything here - just tell terraform we
// are in "created" state and everything is ok
func resourceRunnersCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId("dotscience-runners")
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
	client := m.(*api.Client)
	// ping the API for version to check our credentials
	if _, err := client.Version(); err != nil {
		return fmt.Errorf("Error connecting to the dotscience API: %s", err)
	}

	err := deleteAllTasks(client)

	if err != nil {
		return err
	}

	err = deleteAllRunners(client)

	if err != nil {
		return err
	}

	return nil
}

func resourceRunners() *schema.Resource {
	return &schema.Resource{
		Create: resourceRunnersCreate,
		Read:   resourceRunnersRead,
		Delete: resourceRunnersDelete,

		Schema: map[string]*schema.Schema{},
	}
}
