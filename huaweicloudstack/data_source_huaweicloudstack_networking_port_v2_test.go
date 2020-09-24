package huaweicloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.huaweicloudstack_networking_port_v2.port_1", "id",
						"huaweicloudstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.huaweicloudstack_networking_port_v2.port_2", "id",
						"huaweicloudstack_networking_port_v2.port_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.huaweicloudstack_networking_port_v2.port_3", "id",
						"huaweicloudstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_port_v2.port_3", "all_fixed_ips.#", "1"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortDataSource_basic = `
resource "huaweicloudstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${huaweicloudstack_networking_network_v2.network_1.id}"
  cidr       = "10.0.0.0/24"
  ip_version = 4
}

data "huaweicloudstack_networking_secgroup_v2" "default" {
  name = "default"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${huaweicloudstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.huaweicloudstack_networking_secgroup_v2.default.id}",
  ]

  fixed_ip {
    subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_1.id}"
  }
}

resource "huaweicloudstack_networking_port_v2" "port_2" {
  name               = "port_2"
  network_id         = "${huaweicloudstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.huaweicloudstack_networking_secgroup_v2.default.id}",
  ]
}

data "huaweicloudstack_networking_port_v2" "port_1" {
  name           = "${huaweicloudstack_networking_port_v2.port_1.name}"
  admin_state_up = "${huaweicloudstack_networking_port_v2.port_1.admin_state_up}"

  security_group_ids = [
    "${data.huaweicloudstack_networking_secgroup_v2.default.id}",
  ]
}

data "huaweicloudstack_networking_port_v2" "port_2" {
  name           = "${huaweicloudstack_networking_port_v2.port_2.name}"
  admin_state_up = "${huaweicloudstack_networking_port_v2.port_2.admin_state_up}"
}

data "huaweicloudstack_networking_port_v2" "port_3" {
  fixed_ip = "${huaweicloudstack_networking_port_v2.port_1.all_fixed_ips.0}"
}
`
