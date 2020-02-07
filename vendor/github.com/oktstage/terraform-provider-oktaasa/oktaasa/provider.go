package oktaasa

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"oktaasa_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAASA_KEY", nil),
				Description: "OKTAASA API key.",
			},

			"oktaasa_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAASA_KEY_SECRET", nil),
				Description: "OKTAASA API secret.",
			},

			"oktaasa_team": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAASA_TEAM", nil),
				Description: "OKTAASA Team.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oktaasa_project":          resourceOKTAASAProject(),
			"oktaasa_enrollment_token": resourceOKTAASAToken(),
			"oktaasa_assign_group":     resourceOKTAASAAssignGroup(),
			"oktaasa_create_group":     resourceOKTAASACreateGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		key:    d.Get("oktaasa_key").(string),
		secret: d.Get("oktaasa_secret").(string),
		team:   d.Get("oktaasa_team").(string),
	}

	return config.Authorization()

}
