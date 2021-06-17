package oktaasa

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var projectName string

func TestAccGroupAssign(t *testing.T) {
	groupassign := &Group{}

	project := &Project{}
	projectName := "test-acc-project_g"

	//group := &Group{}
	groupName := "test-acc-group_g"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGroupAssignCheckDestroy(groupassign),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectGroupCreateConfig3,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists3("oktaasa_project.test", project),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_uid", "60120",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_project.test", "next_unix_gid", "63020",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_create_group.test-group", "name", groupName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "group_name", groupName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "server_access", "true",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "server_admin", "true",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "create_server_group", "true",
					),
				),
			},
			{
				Config: testAccGroupAssignUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccGroupAssignCheckExists("oktaasa_assign_group.test-acc-group-assignment", groupassign),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "group_name", groupName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "server_access", "true",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "server_admin", "false",
					),
					resource.TestCheckResourceAttr(
						"oktaasa_assign_group.test-acc-group-assignment", "create_server_group", "true",
					),
				),
			},
		},
	})
}

func testAccProjectCheckExists3(rn string, p *Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		// resource ID is project name
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		err = json.Unmarshal(resp.Body(), p)
		if err != nil {
			return fmt.Errorf("error unmarshaling data source response: %s", err)
		}

		return nil
	}
}

func testAccGroupCheckExists2(rn string, p *Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		// resource ID is group name
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/groups/"+rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		err = json.Unmarshal(resp.Body(), p)
		if err != nil {
			return fmt.Errorf("error unmarshaling data source response: %s", err)
		}

		return nil
	}
}

func testAccGroupAssignCheckExists(rn string, p *Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		// resource ID
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+projectName+"/groups/"+rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		err = json.Unmarshal(resp.Body(), p)
		if err != nil {
			return fmt.Errorf("error unmarshaling data source response: %s", err)
		}

		return nil
	}
}

func testAccGroupAssignCheckDestroy(p *Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+projectName+"/groups/"+p.Name)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		status := resp.StatusCode()
		deleted, err := checkSoftDelete(resp.Body())
		if err != nil {
			return fmt.Errorf("error while checking deleted status: %s", err)
		}

		if status == 200 && !deleted {
			//COME BACK
			return fmt.Errorf("project still exists")
		}

		return nil
	}
}

const testAccProjectGroupCreateConfig3 = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project_g"
  	next_unix_uid = 60120
  	next_unix_gid = 63020
}

resource "oktaasa_create_group" "test-group" {
    name = "test-acc-group_g"
}

resource "oktaasa_assign_group" "test-acc-group-assignment" {
    project_name = oktaasa_project.test.project_name
  	group_name = oktaasa_create_group.test-group.name
	server_access = true
	server_admin = true
	create_server_group = true
}`

const testAccGroupAssignUpdateConfig = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project_g"
  	next_unix_uid = 60120
  	next_unix_gid = 63020
}

resource "oktaasa_create_group" "test-group" {
    name = "test-acc-group_g"
}

resource "oktaasa_assign_group" "test-acc-group-assignment" {
    project_name = oktaasa_project.test.project_name
  	group_name = oktaasa_create_group.test-group.name
	server_access = true
	server_admin = false
	create_server_group = true
}`
