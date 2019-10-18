package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloudstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_secgroup_v2.secgroup_1", "name", "huaweicloudstack_acctest_secgroup"),
				),
			},
		},
	})
}

func TestAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.huaweicloudstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_networking_secgroup_v2.secgroup_1", "name", "huaweicloudstack_acctest_secgroup"),
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

const testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group = `
resource "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
        name        = "huaweicloudstack_acctest_secgroup"
	description = "My neutron security group for huaweicloudstack acctest"
}
`

var testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
	name = "${huaweicloudstack_networking_secgroup_v2.secgroup_1.name}"
}
`, testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group)

var testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "huaweicloudstack_networking_secgroup_v2" "secgroup_1" {
	secgroup_id = "${huaweicloudstack_networking_secgroup_v2.secgroup_1.id}"
}
`, testAccHuaweiCloudStackNetworkingSecGroupV2DataSource_group)
