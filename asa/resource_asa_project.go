package asa

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceASAProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceASAProjectCreate,
		Read:   resourceASAProjectRead,
		Update: resourceASAProjectUpdate,
		Delete: resourceASAProjectDelete,

		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"next_unix_uid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60101,
			},
			"next_unix_gid": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  63001,
			},
		},
	}
}

func resourceASAProjectCreate(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)

	//get project_name from terraform config.
	project_name := d.Get("project_name").(string)
	next_unix_uid := d.Get("next_unix_uid")
	next_unix_gid := d.Get("next_unix_gid")

	// create project in ASA
	project := map[string]interface{}{"name": project_name,
		"create_server_users": true,
		"next_unix_uid":       next_unix_uid,
		"next_unix_gid":       next_unix_gid}
	projectB, _ := json.Marshal(project)

	d.SetId(project_name)
	log.Printf("[DEBUG] Project POST body: %s", projectB)

	//make API call to create project
	resp, err := SendPost(token.BearerToken, "/teams/"+teamName+"/projects", projectB)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when creating project: %s. Error: %s", project_name, err)
	}

	status := resp.StatusCode()

	if status == 201 {
		log.Printf("[INFO] Project %s was successfully created", project_name)
	} else {
		log.Printf("[ERROR] Something went wrong while creating project. Error: %s", resp)
	}

	return resourceASAProjectRead(d, m)
}

type ProjectList struct {
	Projects []Project `json:"list"`
}

type Project struct {
	Name      string `json:"name"`
	DeletedAt string `json:"deleted_at"`
}

func resourceASAProjectRead(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)
	projectName := d.Id()

	resp, err := SendGet(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName)

	if err != nil {
		return fmt.Errorf("[ERROR] Error when reading project state: %s. Error: %s", projectName, err)
	}

	status := resp.StatusCode()

	// API can return 200, but also have deleted_at or removed_at value.
	deleted, err := checkSoftDelete(resp.Body())

	if err != nil {
		return fmt.Errorf("[ERROR] Error when attempting to check for soft delete, while reading project state: %s. Error: %s", projectName, err)
	}

	if status == 200 && deleted == true {
		log.Printf("[INFO] Project %s was removed.", projectName)
		d.SetId("")
		return nil
	} else if status == 200 && deleted == false {
		log.Printf("[INFO] Project %s exists.", projectName)
		return nil
	} else if status == 404 {
		log.Printf("[INFO] Project %s does not exist", projectName)
		d.SetId("")
		return nil
	} else {
		return fmt.Errorf("[DEBUG] failed to read project state. Project: %s Status code: %d", projectName, status)
	}
}

func resourceASAProjectUpdate(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)
	nextUnixUid := d.Get("next_unix_uid")
	nextUnixGid := d.Get("next_unix_gid")

	// create project in ASA
	project := map[string]interface{}{"name": projectName,
		"create_server_users": true,
		"next_unix_uid":       nextUnixUid,
		"next_unix_gid":       nextUnixGid}
	projectB, _ := json.Marshal(project)

	d.SetId(projectName)
	log.Printf("[DEBUG] Project POST body: %s", projectB)

	//make API call to create project
	resp, err := SendPut(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName, projectB)

	if err != nil {
		return fmt.Errorf("[ERROR] Error updating project settings. Project: %s. Error: %s", projectName, err)
	}

	status := resp.StatusCode()

	if status == 204 {
		log.Printf("[INFO] Project %s was successfully updated", projectName)
	} else {
		return fmt.Errorf("[ERROR] Something went wrong while updating the project %s. Error: %s", projectName, resp)

	}

	return resourceASAProjectRead(d, m)
}

func resourceASAProjectDelete(d *schema.ResourceData, m interface{}) error {
	token := m.(Bearer)

	//get project_name from terraform config.
	projectName := d.Get("project_name").(string)

	resp, err := SendDelete(token.BearerToken, "/teams/"+teamName+"/projects/"+projectName, make([]byte, 0))

	if err != nil {
		return fmt.Errorf("[ERROR] Error when deleting project: %s. Error: %s", projectName, err)
	}

	status := resp.StatusCode()

	if status < 300 || status == 400 {
		log.Printf("[INFO] Project %s was successfully deleted", projectName)
	} else {
		log.Printf("[ERROR] Something went wrong while deleting project %s. Error: %s", projectName, resp)
	}

	return nil
}
