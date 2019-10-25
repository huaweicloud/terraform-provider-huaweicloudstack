package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2RouterInterface_basic_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("huaweicloudstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_basic_port(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterface_basic_port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2RouterInterfaceExists("huaweicloudstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_networking_router_interface_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router interface still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterInterfaceExists(n string) resource.TestCheckFunc {
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

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router interface not found")
		}

		return nil
	}
}

const testAccNetworkingV2RouterInterface_basic_subnet = `
resource "huaweicloudstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_router_interface_v2" "int_1" {
  subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  router_id = "${huaweicloudstack_networking_router_v2.router_1.id}"
}

resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2RouterInterface_basic_port = `
resource "huaweicloudstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_router_interface_v2" "int_1" {
  router_id = "${huaweicloudstack_networking_router_v2.router_1.id}"
  port_id = "${huaweicloudstack_networking_port_v2.port_1.id}"
}

resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}
`
