package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.router_acc", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloudstack_networking_router_v2.router_acc", "name", "router_acc_update"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_update_external_gw(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_update_external_gw_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.router_acc", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_update_external_gw_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloudstack_networking_router_v2.router_acc", "external_gateway", OS_EXTGW_ID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterExists(n string, router *routers.Router) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
		}

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingV2Router_basic = `
resource "huaweicloudstack_networking_router_v2" "router_acc" {
	name = "router_acc"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update = `
resource "huaweicloudstack_networking_router_v2" "router_acc" {
	name = "router_acc_update"
	admin_state_up = "true"
	distributed = "false"
}
`

const testAccNetworkingV2Router_update_external_gw_1 = `
resource "huaweicloudstack_networking_router_v2" "router_acc" {
	name = "router_acc"
	admin_state_up = "true"
	distributed = "false"
}
`

var testAccNetworkingV2Router_update_external_gw_2 = fmt.Sprintf(`
resource "huaweicloudstack_networking_router_v2" "router_acc" {
	name = "router_acc"
	admin_state_up = "true"
	distributed = "false"
	external_gateway = "%s"
}
`, OS_EXTGW_ID)
