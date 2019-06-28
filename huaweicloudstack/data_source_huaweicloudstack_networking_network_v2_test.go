package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network,
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", "huaweicloudstack_test_network"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network,
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", "huaweicloudstack_test_network"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network,
			},
			{
				Config: testAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.huaweicloudstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_network_v2.net", "name", "huaweicloudstack_test_network"),
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

const testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network = `
resource "huaweicloudstack_networking_network_v2" "net" {
        name = "huaweicloudstack_test_network"
        admin_state_up = "true"
}

resource "huaweicloudstack_networking_subnet_v2" "subnet" {
  name = "huaweicloudstack_test_subnet"
  cidr = "192.168.198.0/24"
  no_gateway = true
  network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}
`

var testAccHuaweiCloudStackNetworkingNetworkV2DataSource_basic = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_network_v2" "net" {
	name = "${huaweicloudstack_networking_network_v2.net.name}"
}
`, testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network)

var testAccHuaweiCloudStackNetworkingNetworkV2DataSource_subnet = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_network_v2" "net" {
	matching_subnet_cidr = "${huaweicloudstack_networking_subnet_v2.subnet.cidr}"
}
`, testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network)

var testAccHuaweiCloudStackNetworkingNetworkV2DataSource_networkID = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_network_v2" "net" {
	network_id = "${huaweicloudstack_networking_network_v2.net.id}"
}
`, testAccHuaweiCloudStackNetworkingNetworkV2DataSource_network)
