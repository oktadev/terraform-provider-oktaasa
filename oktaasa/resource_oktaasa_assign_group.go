package oktaasa

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOKTAASAAssignGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceOKTAASAAssignGroupCreate,
		Read:   resourceOKTAASAAssignGroupRead,
		Update: resourceOKTAASAAssignGroupUpdate,
		Delete: resourceOKTAASAAssignGroupDelete,

		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"server_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"create_server_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceOKTAASAAssignGroupCreate(d *schema.ResourceData, m interface{}) error {

	token := m.(Bearer)

	//get settings from terraform config.
	project_name := d.Get("project_name").(string)
	oktaGroupName := d.Get("group_name").(string)
	serverAccess := d.Get("server_access")
	serverAdmin := d.Get("server_admin")
	groupSync := d.Get("create_server_group")

	log.Printf("[DEBUG] Assigning group %s to the project: %s", oktaGroupName, project_name)

	// set POST parameters for group assignment
	oktaGroupSettings := map[string]interface{}{
		"group":               oktaGroupName,
		"server_access":       serverAccess,
		"server_admin":        serverAdmin,
		"create_server_group": groupSync}

	oktaGroupDescriptionB, _ := json.Marshal(oktaGroupSettings)

	//make API call to assign Okta group to a project
	resp, err := SendPost(token.BearerToken, "/teams/"+teamName+"/projects/"+project_name+"/groups", oktaGroupDescriptionB)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when assigning group: %s. Error: %s", oktaGroupName, err)
	}

	statusCode := resp.StatusCode()
	if statusCode < 300 {
		log.Printf("[DEBUG] Success. Group %s was assigned to %s", oktaGroupName, project_name)
	} else {
		return fmt.Errorf("[ERROR] Error happened while assigning group %s to a project: %s", oktaGroupName, resp)
	}

	d.SetId(oktaGroupName)

	return resourceOKTAASAAssignGroupRead(d, m)
}

type Group struct {
	Name         string `json:"group"`
	ServerAccess bool   `json:"server_access"`
	ServerAdmin  bool   `json:"server_admin"`
	GroupSync    bool   `json:"create_server_group"`
}

func resourceOKTAASAAssignGroupRead(d *schema.ResourceData, m interface{}) error {
	sessionToken := m.(Bearer)
	groupName := d.Id()

	log.Printf("[INFO] Group ID is: %s", groupName)
	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)

	resp, err := SendGet(sessionToken.BearerToken, "/teams/"+teamName+"/projects/"+projectName+"/groups/"+groupName)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when reading group state: %s. Error: %s", groupName, err)
	}

	status := resp.StatusCode()

	// API can return 200, but also have deleted_at or removed_at value.
	deleted, err := checkSoftDelete(resp.Body())
	if err != nil {
		return fmt.Errorf("[ERROR] Error when attempting to check for soft delete, while reading group state: %s. Error: %s", groupName, err)
	}

	if status == 200 && deleted {
		log.Printf("[INFO] Group %s was removed from project %s", groupName, projectName)
		d.SetId("")
		return nil
	} else if status == 200 && !deleted {
		log.Printf("[INFO] Group %s is assigned to project %s ", groupName, projectName)

		var group Group

		err := json.Unmarshal(resp.Body(), &group)
		if err != nil {
			return fmt.Errorf("unable to unmarshal group settings: %w", err)
		}

		d.SetId(group.Name)
		d.Set("server_access", group.ServerAccess)
		d.Set("server_admin", group.ServerAdmin)
		d.Set("create_server_group", group.GroupSync)

		return nil
	} else if status == 404 {
		log.Printf("[INFO] group %s is no assigned to the project. %s", groupName, projectName)
		d.SetId("")
		return nil
	} else {
		return fmt.Errorf("[DEBUG] failed to read group assignment state. Project: %s Group: %s Status code: %d", projectName, groupName, status)
	}
}

func resourceOKTAASAAssignGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceOKTAASAAssignGroupCreate(d, m)
}

func resourceOKTAASAAssignGroupDelete(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)
	groupName := d.Id()

	resp, err := SendDelete(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName+"/groups/"+groupName, make([]byte, 0))

	if err != nil {
		return fmt.Errorf("[ERROR] Error when deleting group: %s. Error: %s", groupName, err)
	}

	status := resp.StatusCode()

	if status < 300 || status == 404 {
		log.Printf("[INFO] Group %s was successfully deleted", projectName)
	} else {
		return fmt.Errorf("[ERROR] Something went wrong while deleting group %s. Error: %s", d.Id(), resp)
	}

	return nil
}
