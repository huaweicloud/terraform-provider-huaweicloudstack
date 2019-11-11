package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic(t *testing.T) {
	var netName = fmt.Sprintf("acc_net_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network(netName),
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic(netName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", netName),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	var netName = fmt.Sprintf("acc_net_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network(netName),
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet(netName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", netName),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	var netName = fmt.Sprintf("acc_net_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network(netName),
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID(netName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", netName),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network data source ID not set")
		}

		return nil
	}
}

func testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network(netName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet" {
  name = "huaweicloudstack_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}
`, netName)
}

func testAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic(netName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet" {
  name = "huaweicloudstack_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}

data "huaweicloudstack_networking_network_v2" "net" {
	name = "${huaweicloudstack_networking_network_v2.net.name}"
}
`, netName)
}

func testAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet(netName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet" {
  name = "huaweicloudstack_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}

data "huaweicloudstack_networking_network_v2" "net" {
	matching_subnet_cidr = "${huaweicloudstack_networking_subnet_v2.subnet.cidr}"
}
`, netName)
}

func testAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID(netName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_network_v2" "net" {
  name = "%s"
  admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet" {
  name = "huaweicloudstack_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}

data "huaweicloudstack_networking_network_v2" "net" {
	network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}
`, netName)
}
