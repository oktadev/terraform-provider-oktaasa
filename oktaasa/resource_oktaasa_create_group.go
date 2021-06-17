package oktaasa

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceOKTAASACreateGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOKTAASACreateGroupCreate,
		Read:   resourceOKTAASACreateGroupRead,
		Update: resourceOKTAASACreateGroupUpdate,
		Delete: resourceOKTAASACreateGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceOKTAASACreateGroupCreate(d *schema.ResourceData, m interface{}) error {

	token := m.(Bearer)

	//get settings from terraform config.
	oktaasaGroupName := d.Get("name").(string)

	log.Printf("[DEBUG] Creating group %s", oktaasaGroupName)

	type RolesOptions struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	}

	options := &RolesOptions{Name: oktaasaGroupName, Roles: []string{}}

	GroupDescriptionB, _ := json.Marshal(options)

	//make API call to assign Okta group to a project
	resp, err := SendPost(token.BearerToken, "/teams/"+teamName+"/groups", GroupDescriptionB)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when creating group: %s. Error: %s", oktaasaGroupName, err)
	}

	statusCode := resp.StatusCode()

	if statusCode < 300 {
		log.Printf("[DEBUG] Success. Group %s was created", oktaasaGroupName)
	} else if statusCode == 409 {
		log.Printf("[INFO] Group already exists")
	} else {
		return fmt.Errorf("[ERROR] Error happened while creating group %s. Error: %s", oktaasaGroupName, resp)
	}

	d.SetId(oktaasaGroupName)

	return resourceOKTAASACreateGroupRead(d, m)
}

type SftGroup struct {
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	DeletedAt string   `json:"deleted_at"`
}

func resourceOKTAASACreateGroupRead(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)
	groupName := d.Id()

	log.Printf("[DEBUG Running Create Group Read function.")
	log.Printf("[INFO] Group name is: %s", groupName)

	resp, err := SendGet(token.BearerToken, "/teams/"+teamName+"/groups/"+groupName)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when reading group state: %s. Error: %s", groupName, err)
	}

	status := resp.StatusCode()

	// API can return 200, but also have deleted_at or removed_at value.
	deleted, err := checkSoftDelete(resp.Body())
	if err != nil {
		return fmt.Errorf("[ERROR] Error when attempting to check for soft delete, while reading group state: %s. Error: %s", groupName, err)
	}

	if status == 200 && deleted == true {
		log.Printf("[INFO] Group %s was removed Need to recreate.", groupName)
		d.SetId("")
		return nil
	} else if status == 200 && deleted == false {
		log.Printf("[INFO] Group %s exists.", groupName)
		return nil
	} else if status == 404 {
		log.Printf("[INFO] group %s does not exist.", groupName)
		d.SetId("")
		return nil
	} else {
		return fmt.Errorf("[DEBUG] failed to read group state. Group: %s Status code: %d", groupName, status)
	}

}

func resourceOKTAASACreateGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceOKTAASACreateGroupCreate(d, m)
}

func resourceOKTAASACreateGroupDelete(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)
	groupName := d.Id()

	resp, err := SendDelete(token.BearerToken, "/teams/"+teamName+"/groups/"+groupName, make([]byte, 0))

	if err != nil {
		return fmt.Errorf("[ERROR] Error when deleting group: %s. Error: %s", groupName, err)
	}

	status := resp.StatusCode()

	if status < 300 || status == 404 {
		log.Printf("[INFO] Group %s was successfully deleted", groupName)
	} else {
		return fmt.Errorf("[ERROR] Something went wrong while deleting group %s. Error: %s", d.Id(), resp)
	}

	return nil
}
