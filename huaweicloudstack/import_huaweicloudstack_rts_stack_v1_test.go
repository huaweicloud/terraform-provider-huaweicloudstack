package huaweicloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRtsStackV1_importBasic(t *testing.T) {
	resourceName := "huaweicloudstack_rts_stack_v1.stack_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRtsStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRtsStackV1_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
