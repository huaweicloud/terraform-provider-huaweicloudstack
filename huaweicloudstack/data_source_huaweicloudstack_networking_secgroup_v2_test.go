package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	var secName = fmt.Sprintf("acc_sec_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group(secName),
			},
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic(secName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloudstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_secgroup_v2.secgroup_1", "name", secName),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	var secName = fmt.Sprintf("acc_sec_%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group(secName),
			},
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID(secName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloudstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_secgroup_v2.secgroup_1", "name", secName),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

func testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group(secName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group for huaweicloudstack acctest"
}
`, secName)
}

func testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic(secName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group for huaweicloudstack acctest"
}

data "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
	name = "${huaweicloudstack_networking_secgroup_v2.secgroup_1.name}"
}
`, secName)
}

func testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID(secName string) string {
	return fmt.Sprintf(`
resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group for huaweicloudstack acctest"
}

data "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${huaweicloudstack_networking_secgroup_v2.secgroup_1.id}"
}
`, secName)
}
