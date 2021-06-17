package oktaasa

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceOKTAASAToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceOKTAASATokenCreate,
		Read:   resourceOKTAASATokenRead,
		Update: resourceOKTAASATokenUpdate,
		Delete: resourceOKTAASATokenDelete,

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Computed
			"token_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type EnrollmentToken struct {
	Id    string `json:"id"`
	Value string `json:"token"`
}

func resourceOKTAASATokenCreate(d *schema.ResourceData, m interface{}) error {
	// Bearer session token
	token := m.(Bearer)

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)
	description := d.Get("description").(string)

	// create project in OKTAASA
	tokenDescription := map[string]interface{}{"description": description}
	tokenDescriptionB, _ := json.Marshal(tokenDescription)

	//make API call to create project
	resp, err := SendPost(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName+"/server_enrollment_tokens", tokenDescriptionB)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when creating enrollment token: %s. Error: %s", description, err)
	}

	status := resp.StatusCode()

	if status > 400 {
		return fmt.Errorf("[ERROR] Unexpected error when creating token %d", status)
	}

	enrollmentToken := EnrollmentToken{}

	jsonErr := json.Unmarshal(resp.Body(), &enrollmentToken)
	if jsonErr != nil {
		log.Printf("[DEBUG] Error storing enrollment token: %s", jsonErr)
	}

	// update resource ID with token ID.
	d.SetId(enrollmentToken.Id)

	return resourceOKTAASATokenRead(d, m)
}

func resourceOKTAASATokenRead(d *schema.ResourceData, m interface{}) error {
	sessionToken := m.(Bearer)
	tokenId := d.Id()

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)

	resp, err := SendGet(sessionToken.BearerToken, "/teams/"+teamName+"/projects/"+projectName+"/server_enrollment_tokens/"+tokenId)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when reading token state. Token: %s. Error: %s", tokenId, err)
	}

	status := resp.StatusCode()

	if status == 200 {
		log.Printf("[DEBUG] Token %s exists", tokenId)

		var tokenInfo EnrollmentToken
		err := json.Unmarshal([]byte(resp.Body()), &tokenInfo)

		if err != nil {
			return fmt.Errorf("[ERROR] Error when reading token state. Token: %s. Error: %s", tokenId, err)
		}

		d.Set("token_value", tokenInfo.Value)

	} else if status == 404 {
		log.Printf("[DEBUG] No token %s in this project", tokenId)
		d.SetId("")
		return nil
	} else {
		return fmt.Errorf("[ERROR] Something went wrong while retrieving a list of enrollment tokens. Error: %s", resp)
	}
	return nil
}

func resourceOKTAASATokenUpdate(d *schema.ResourceData, m interface{}) error {
	// not possible to update token.
	return nil
}

func resourceOKTAASATokenDelete(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)
	tokenId := d.Id()

	resp, err := SendDelete(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName+"/server_enrollment_tokens/"+tokenId, make([]byte, 0))

	if err != nil {
		return fmt.Errorf("[ERROR] Error when deleting token: %s. Error: %s", tokenId, err)
	}

	status := resp.StatusCode()

	if status < 300 || status == 404 {
		log.Printf("[INFO] Enrollment token %s of a project %s was successfully deleted", d.Id(), projectName)
	} else {
		return fmt.Errorf("[ERROR] Error while deleting token: %s", resp)
	}

	return nil
}
