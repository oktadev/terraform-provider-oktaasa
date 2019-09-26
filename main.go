package main

import (
	"github.com/oktstage/terraform-provider-asa/asa"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: asa.Provider})
}
