package limelight

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLimelightIPRanges_basic(t *testing.T) {
	testResourceName := "data.limelight_ip_ranges.ips"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLimelightIPRangesBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "version"),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_ranges.#"),
				),
			},
		},
	})
}

func testAccDataSourceLimelightIPRangesBasicTemplate() string {
	return `
data "limelight_ip_ranges" "ips" {}
`
}
