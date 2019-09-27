package asa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"asa_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ASA_KEY", nil),
				Description: "ASA API key.",
			},

			"asa_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ASA_KEY_SECRET", nil),
				Description: "ASA API secret.",
			},

			"asa_team": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ASA_TEAM", nil),
				Description: "ASA Team.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"asa_project":          resourceASAProject(),
			"asa_enrollment_token": resourceASAToken(),
			"asa_assign_group":     resourceASAAssignGroup(),
			"asa_create_group":     resourceASACreateGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		key:    d.Get("asa_key").(string),
		secret: d.Get("asa_secret").(string),
		team:   d.Get("asa_team").(string),
	}

	return config.Authorization()

}
