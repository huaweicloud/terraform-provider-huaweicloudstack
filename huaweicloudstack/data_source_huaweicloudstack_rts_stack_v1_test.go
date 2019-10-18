package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRtsStackV1DataSource_basic(t *testing.T) {
	var stackName = fmt.Sprintf("terra-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsStackV1DataSource_basic(stackName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRtsStackV1DataSourceID("data.huaweicloudstack_rts_stack_v1.stacks"),
					resource.TestCheckResourceAttr("data.huaweicloudstack_rts_stack_v1.stacks", "name", stackName),
					resource.TestCheckResourceAttr("data.huaweicloudstack_rts_stack_v1.stacks", "disable_rollback", "true"),
					resource.TestCheckResourceAttr("data.huaweicloudstack_rts_stack_v1.stacks", "parameters.%", "4"),
					resource.TestCheckResourceAttr("data.huaweicloudstack_rts_stack_v1.stacks", "status", "CREATE_COMPLETE"),
				),
			},
		},
	})
}

func testAccCheckRtsStackV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find rts data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("RTS data source ID not set ")
		}

		return nil
	}
}

func testAccRtsStackV1DataSource_basic(stackName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_rts_stack_v1" "stack_1" {
  name = "%s"
  disable_rollback= true
  timeout_mins=60
  template_body = <<JSON
          {
			"outputs": {
              "str1": {
                 "description": "The description of the nat server.",
                 "value": {
                   "get_resource": "random"
                 }
	          }
            },
            "heat_template_version": "2013-05-23",
            "description": "A HOT template that create a single server and boot from volume.",
            "parameters": {
              "key_name": {
                "type": "string",
                "description": "Name of existing key pair for the instance to be created.",
                "default": "KeyPair-click2cloud"
	          }
	        },
            "resources": {
               "random": {
                  "type": "OS::Heat::RandomString",
                  "properties": {
                  "length": "6"
                  }
	          }
	       }
}
JSON
}

data "huaweicloudstack_rts_stack_v1" "stacks" {
        name = "${huaweicloudstack_rts_stack_v1.stack_1.name}"
}
`, stackName)
}
