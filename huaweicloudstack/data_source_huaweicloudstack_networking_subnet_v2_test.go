package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2SubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork(
						"data.huaweicloudstack_networking_subnet_v2.subnet_1",
						"huaweicloudstack_networking_network_v2.network_test1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_subnet_v2.subnet_1", "name", "subnet_test1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_dhcpEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_ipVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
				),
			},
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_gatewayIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_networkIdAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHCSNetworkingSubnetV2DataSource_networkIdAttribute,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.huaweicloudstack_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork(
						"data.huaweicloudstack_networking_subnet_v2.subnet_1",
						"huaweicloudstack_networking_network_v2.network_test1"),
					testAccCheckNetworkingPortV2ID("huaweicloudstack_networking_port_v2.port_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingPortV2ID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find port resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Port resource ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetV2DataSourceGoodNetwork(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmt.Errorf("Can't find network resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmt.Errorf("Network resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["network_id"] {
			return fmt.Errorf("Network id and subnet network_id don't match")
		}

		return nil
	}
}

const testAccHCSNetworkingSubnetV2DataSource_subnet = `
resource "huaweicloudstack_networking_network_v2" "network_test1" {
  name = "network_test1"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet_test1" {
  name = "subnet_test1"
  cidr = "172.16.199.0/24"
  network_id = "${huaweicloudstack_networking_network_v2.network_test1.id}"
}
`

var testAccHCSNetworkingSubnetV2DataSource_basic = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  name = "${huaweicloudstack_networking_subnet_v2.subnet_test1.name}"
}
`, testAccHCSNetworkingSubnetV2DataSource_subnet)

var testAccHCSNetworkingSubnetV2DataSource_cidr = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  network_id = "${huaweicloudstack_networking_network_v2.network_test1.id}"
  cidr = "172.16.199.0/24"
}
`, testAccHCSNetworkingSubnetV2DataSource_subnet)

var testAccHCSNetworkingSubnetV2DataSource_dhcpEnabled = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  network_id = "${huaweicloudstack_networking_network_v2.network_test1.id}"
  dhcp_enabled = true
}
`, testAccHCSNetworkingSubnetV2DataSource_subnet)

var testAccHCSNetworkingSubnetV2DataSource_ipVersion = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  network_id = "${huaweicloudstack_networking_network_v2.network_test1.id}"
  ip_version = 4
}
`, testAccHCSNetworkingSubnetV2DataSource_subnet)

var testAccHCSNetworkingSubnetV2DataSource_gatewayIP = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  gateway_ip = "${huaweicloudstack_networking_subnet_v2.subnet_test1.gateway_ip}"
}
`, testAccHCSNetworkingSubnetV2DataSource_subnet)

var testAccHCSNetworkingSubnetV2DataSource_networkIdAttribute = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_subnet_v2" "subnet_1" {
  subnet_id = "${huaweicloudstack_networking_subnet_v2.subnet_test1.id}"
}

resource "huaweicloudstack_networking_port_v2" "port_1" {
  name               = "test_port"
  network_id         = "${data.huaweicloudstack_networking_subnet_v2.subnet_1.network_id}"
  admin_state_up  = "true"
}

`, testAccHCSNetworkingSubnetV2DataSource_subnet)
