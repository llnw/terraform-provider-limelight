package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/llnw/terraform-provider-limelight/limelight"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: limelight.Provider})
}
