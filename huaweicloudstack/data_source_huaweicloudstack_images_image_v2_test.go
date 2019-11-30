package huaweicloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const image_name = "Cirros-raw"

func TestAccImagesV2ImageDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckImage(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHCSImagesV2ImageDataSource_basic(image_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloudstack_images_image_v2.image_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "name", image_name),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "container_format", "bare"),
					/*resource.TestCheckResourceAttr(
							"data.huaweicloudstack_images_image_v2.image_1", "disk_format", "qcow2"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "min_disk_gb", "0"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "min_ram_mb", "0"),*/
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "protected", "true"),
					resource.TestCheckResourceAttr(
						"data.huaweicloudstack_images_image_v2.image_1", "visibility", "public"),
				),
			},
		},
	})
}

func TestAccImagesV2ImageDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckImage(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHCSImagesV2ImageDataSource_querySizeMin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloudstack_images_image_v2.image_1"),
				),
			},
			{
				Config: testAccHCSImagesV2ImageDataSource_querySizeMax,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.huaweicloudstack_images_image_v2.image_1"),
				),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccHCSImagesV2ImageDataSource_basic(image_name string) string {
	return fmt.Sprintf(`
data "huaweicloudstack_images_image_v2" "image_1" {
	most_recent = true
	name = "%s"
}
`, image_name)
}

const testAccHCSImagesV2ImageDataSource_querySizeMin = `
data "huaweicloudstack_images_image_v2" "image_1" {
	most_recent = true
	visibility = "public"
	size_min = "13000000"
    sort_key = "size"
    sort_direction = "desc"
}`

const testAccHCSImagesV2ImageDataSource_querySizeMax = `
data "huaweicloudstack_images_image_v2" "image_1" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}`
