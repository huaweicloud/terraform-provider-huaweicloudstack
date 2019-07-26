package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/groups"
	"log"
)

func TestAccASV1Group_basic(t *testing.T) {
	var asGroup groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Group_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1GroupExists("huaweicloudstack_as_group_v1.as_group_1", &asGroup),
				),
			},
		},
	})
}

func testAccCheckASV1GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloudstack autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloudstack_as_group_v1" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1GroupDestroy success!")

	return nil
}

func testAccCheckASV1GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		asClient, err := config.autoscalingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloudstack autoscaling client: %s", err)
		}

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Group not found")
		}
		log.Printf("[DEBUG] test found is: %#v", found)
		group = &found

		return nil
	}
}

var testASV1Group_basic = fmt.Sprintf(`
resource "huaweicloudstack_networking_secgroup_v2" "secgroup" {
  name        = "terraform"
  description = "This is a terraform test security group"
}

resource "huaweicloudstack_compute_keypair_v2" "key_1" {
  name = "key_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloudstack_as_configuration_v1" "as_config_1"{
  scaling_configuration_name = "as_config_1"
  instance_config = {
    image = "%s"
    disk {
      size = 40
      volume_type = "SATA"
      disk_type = "SYS"
    }
    key_name = "${huaweicloudstack_compute_keypair_v2.key_1.id}"
  }
}

resource "huaweicloudstack_as_group_v1" "as_group_1"{
  scaling_group_name = "as_group_1"
  scaling_configuration_id = "${huaweicloudstack_as_configuration_v1.as_config_1.id}"
  networks {
    id = "%s"
  }
  security_groups {
    id = "${huaweicloudstack_networking_secgroup_v2.secgroup.id}"
  }
  vpc_id = "%s"
}
`, OS_IMAGE_ID, OS_NETWORK_ID, OS_VPC_ID)
