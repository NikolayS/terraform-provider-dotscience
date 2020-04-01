package main

import (
	"github.com/dotmesh-io/terraform-provider-dotscience/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
