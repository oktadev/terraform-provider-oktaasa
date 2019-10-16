package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/oktstage/terraform-provider-asa/asa"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: asa.Provider})
}
