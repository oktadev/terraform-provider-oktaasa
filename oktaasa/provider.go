package oktaasa

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"oktaasa_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAASA_KEY", nil),
				Description: "OKTAASA API key.",
			},

			"oktaasa_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAASA_KEY_SECRET", nil),
				Description: "OKTAASA API secret.",
			},

			"oktaasa_team": {
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
