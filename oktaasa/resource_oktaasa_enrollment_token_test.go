package oktaasa

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTkn(t *testing.T) {
	//tkn := &EnrollmentToken{}
	//tknName := "test-acc-token"
	projectName := "test-acc-project2"
	project := &Project{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy2(project),
		Steps: []resource.TestStep{
			{
				Config: testAccProjectTokenCreateConfig2,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists2("oktaasa_project.test", project),
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
						"oktaasa_enrollment_token.test-token", "project_name", projectName,
					),
					resource.TestCheckResourceAttr(
						"oktaasa_enrollment_token.test-token", "description", "Token for TestAcc",
					),
				),
			},
			//Note: OKTAASA does not allow token or token description changes once created (hence there is no Update step)
		},
	})
}

func testAccProjectCheckExists2(rn string, p *Project) resource.TestCheckFunc {
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

func testAccProjectCheckDestroy2(p *Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+p.Name)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		status := resp.StatusCode()
		deleted, err := checkSoftDelete(resp.Body())
		if err != nil {
			return fmt.Errorf("error while checking deleted status: %s", err)
		}

		if status == 200 && !deleted {
			return fmt.Errorf("project still exists")
		}

		return nil
	}
}

func testAccTknCheckExists(rn string, p *EnrollmentToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		projectName := "test-acc-project2"

		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		// resource ID is token name
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/projects/"+projectName+"/server_enrollment_tokens/"+rs.Primary.ID)
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

func testAccTknCheckDestroy(p *EnrollmentToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/server_enrollment_tokens/"+p.Id)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		status := resp.StatusCode()
		deleted, err := checkSoftDelete(resp.Body())
		if err != nil {
			return fmt.Errorf("error while checking deleted status: %s", err)
		}

		if status == 200 && !deleted {
			return fmt.Errorf("token still exists")
		}

		return nil
	}
}

const testAccProjectTokenCreateConfig2 = `
resource "oktaasa_project" "test" {
    project_name = "test-acc-project2"
  	next_unix_uid = 60120
  	next_unix_gid = 63020
}

resource "oktaasa_enrollment_token" "test-token" {
    project_name = oktaasa_project.test.project_name
  	description = "Token for TestAcc"
}`
