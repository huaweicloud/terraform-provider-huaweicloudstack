package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func TestAccLBV2Member_basic(t *testing.T) {
	var member_1 pools.Member
	var member_2 pools.Member

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckULB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:             TestAccLBV2MemberConfig_basic,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false.
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MemberExists("huaweicloudstack_lb_member_v2.member_1", &member_1),
					testAccCheckLBV2MemberExists("huaweicloudstack_lb_member_v2.member_2", &member_2),
				),
			},
			{
				Config:             TestAccLBV2MemberConfig_update,
				ExpectNonEmptyPlan: true, // Because admin_state_up remains false.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloudstack_lb_member_v2.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("huaweicloudstack_lb_member_v2.member_2", "weight", "15"),
				),
			},
		},
	})
}

func testAccCheckLBV2MemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloudStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_lb_member_v2" {
			continue
		}

		poolId := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(networkingClient, poolId, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
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

		poolId := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(networkingClient, poolId, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*member = *found

		return nil
	}
}

var TestAccLBV2MemberConfig_basic = fmt.Sprintf(`
resource "huaweicloudstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "huaweicloudstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${huaweicloudstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "huaweicloudstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${huaweicloudstack_lb_listener_v2.listener_1.id}"
}

resource "huaweicloudstack_lb_member_v2" "member_1" {
  address = "192.168.1.10"
  protocol_port = 8080
  pool_id = "${huaweicloudstack_lb_pool_v2.pool_1.id}"
  subnet_id = "%s"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloudstack_lb_member_v2" "member_2" {
  address = "192.168.1.11"
  protocol_port = 8080
  pool_id = "${huaweicloudstack_lb_pool_v2.pool_1.id}"
  subnet_id = "%s"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID, OS_SUBNET_ID, OS_SUBNET_ID)

var TestAccLBV2MemberConfig_update = fmt.Sprintf(`
resource "huaweicloudstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "huaweicloudstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = "${huaweicloudstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "huaweicloudstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${huaweicloudstack_lb_listener_v2.listener_1.id}"
}

resource "huaweicloudstack_lb_member_v2" "member_1" {
  address = "192.168.1.10"
  protocol_port = 8080
  weight = 10
  admin_state_up = "true"
  pool_id = "${huaweicloudstack_lb_pool_v2.pool_1.id}"
  subnet_id = "%s"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "huaweicloudstack_lb_member_v2" "member_2" {
  address = "192.168.1.11"
  protocol_port = 8080
  weight = 15
  admin_state_up = "true"
  pool_id = "${huaweicloudstack_lb_pool_v2.pool_1.id}"
  subnet_id = "%s"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID, OS_SUBNET_ID, OS_SUBNET_ID)
