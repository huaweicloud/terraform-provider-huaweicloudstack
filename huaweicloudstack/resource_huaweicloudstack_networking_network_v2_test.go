package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Network_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Network_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
				),
			},
			{
				ResourceName:      "huaweicloudstack_networking_network_v2.network_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkingV2Network_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloudstack_networking_network_v2.network_1", "name", "network_2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router
	var routerName = fmt.Sprintf("acc_router_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Network_netstack(routerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.router_acc", &router),
					testAccCheckNetworkingV2RouterInterfaceExists(
						"huaweicloudstack_networking_router_interface_v2.ri_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_fullstack(t *testing.T) {
	var instance servers.Server
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var secgroup groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Network_fullstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("huaweicloudstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("huaweicloudstack_networking_router_v2.vpc_1", &router),
					testAccCheckNetworkingV2SecGroupExists("huaweicloudstack_networking_secgroup_v2.secgroup_1", &secgroup),
					testAccCheckNetworkingV2PortExists("huaweicloudstack_networking_port_v2.port_1", &port),
					testAccCheckComputeV2InstanceExists("huaweicloudstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_multipleSegmentMappings(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Network_multipleSegmentMappings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("huaweicloudstack_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2NetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_networking_network_v2" {
			continue
		}

		_, id := ExtractValFromNid(rs.Primary.ID)
		_, err := networks.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2NetworkExists(n string, network *networks.Network) resource.TestCheckFunc {
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

		_, id := ExtractValFromNid(rs.Primary.ID)
		found, err := networks.Get(networkingClient, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

const testAccNetworkingV2Network_basic = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Network_update = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_2"
}
`

func testAccNetworkingV2Network_netstack(routerName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_router_v2" "router_acc" {
  name = "%s"
}

resource "huaweicloudstack_networking_router_interface_v2" "ri_1" {
  router_id = "${huaweicloudstack_networking_router_v2.router_acc.id}"
  subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
}
`, routerName)
}

const testAccNetworkingV2Network_fullstack = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
}

resource "huaweicloudstack_networking_router_v2" "vpc_1" {
  name = "vpc_1"
}

resource "huaweicloudstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${huaweicloudstack_networking_router_v2.vpc_1.id}"
  subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloudstack_networking_router_route_v2" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.100"

  router_id = "${huaweicloudstack_networking_router_v2.vpc_1.id}"
  depends_on = ["huaweicloudstack_networking_router_interface_v2.router_interface_1"]
}

resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  security_group_ids = ["${huaweicloudstack_networking_secgroup_v2.secgroup_1.id}"]
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
    ip_address =  "192.168.199.23"
  }
  depends_on = ["huaweicloudstack_networking_router_interface_v2.router_interface_1"]
}

resource "huaweicloudstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["${huaweicloudstack_networking_secgroup_v2.secgroup_1.name}"]

  network {
    port = "${huaweicloudstack_networking_port_v2.port_1.id}"
  }
}
`

const testAccNetworkingV2Network_multipleSegmentMappings = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name = "network_1"
  segments =[
    {
      segmentation_id = 2,
      network_type = "vxlan"
    }
  ],
  admin_state_up = "true"
}
`
