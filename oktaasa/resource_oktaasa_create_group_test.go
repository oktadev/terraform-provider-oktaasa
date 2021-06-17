package oktaasa

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup(t *testing.T) {
	group := &Group{}
	groupName := "test-acc-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGroupCheckDestroy(group),
		Steps: []resource.TestStep{
			{
				Config: testAccGroupCreateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccGroupCheckExists("oktaasa_create_group.test-group", group),
					resource.TestCheckResourceAttr(
						"oktaasa_create_group.test-group", "name", groupName,
					),
				),
			},
			//Note: OKTAASA does not allow a group name change once created (hence there is no Update step)
		},
	})
}

func testAccGroupCheckExists(rn string, p *Group) resource.TestCheckFunc {
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

func testAccGroupCheckDestroy(p *Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(Bearer)

		resp, err := SendGet(config.BearerToken, "/teams/"+config.TeamName+"/groups/"+p.Name)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}

		status := resp.StatusCode()
		deleted, err := checkSoftDelete(resp.Body())
		if err != nil {
			return fmt.Errorf("error while checking deleted status: %s", err)
		}

		if status == 200 && !deleted {
			return fmt.Errorf("group still exists")
		}

		return nil
	}
}

const testAccGroupCreateConfig = `
resource "oktaasa_create_group" "test-group" {
    name = "test-acc-group"
}`
