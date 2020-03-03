package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-oktaasa/oktaasa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: oktaasa.Provider})
}
