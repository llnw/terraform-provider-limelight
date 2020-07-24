package limelight

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceLimelightDelivery_minimal(t *testing.T) {
	testResourceName := "limelight_delivery.test_delivery"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightDeliveryCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightDeliveryMinimalTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightDeliveryExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_number"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "published_hostname",
						fmt.Sprintf("terraform-test.%s.s.llnwi.net", getShortname())),
					resource.TestCheckResourceAttr(testResourceName, "published_path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "service_profile", "LLNW-Generic"),
					resource.TestCheckResourceAttr(testResourceName, "source_hostname", "dummy-origin.llnw.net"),
					resource.TestCheckResourceAttr(testResourceName, "source_path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.published_protocol", "https"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.source_protocol", "https"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.option.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceLimelightDelivery_update(t *testing.T) {
	testResourceName := "limelight_delivery.test_delivery"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightDeliveryCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightDeliveryBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightDeliveryExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_number"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "published_hostname",
						fmt.Sprintf("terraform-test-basic.%s.s.llnwi.net", getShortname())),
					resource.TestCheckResourceAttr(testResourceName, "published_path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "service_profile", "LLNW-Generic"),
					resource.TestCheckResourceAttr(testResourceName, "source_hostname", "dummy-origin-basic.llnw.net"),
					resource.TestCheckResourceAttr(testResourceName, "source_path", "/"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.published_protocol", "https"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.source_protocol", "https"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.option.#", "1"),
				),
			},
			{
				Config: testAccLimelightDeliveryBasicTemplate_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightDeliveryExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "version_number"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "published_hostname",
						fmt.Sprintf("terraform-test-basic.%s.s.llnwi.net", getShortname())),
					resource.TestCheckResourceAttr(testResourceName, "published_path", "/publish"),
					resource.TestCheckResourceAttr(testResourceName, "service_profile", "LLNW-Generic"),
					resource.TestCheckResourceAttr(testResourceName, "source_hostname", "dummy-origin-basic.llnw.net"),
					resource.TestCheckResourceAttr(testResourceName, "source_path", "/source"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.published_protocol", "http"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.source_protocol", "http"),
					resource.TestCheckResourceAttr(testResourceName, "protocol_set.0.option.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceLimelightDelivery_importBasic(t *testing.T) {
	testResourceName := "limelight_delivery.test_delivery"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightDeliveryCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightDeliveryImportTemplate(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLimelightDeliveryCheckDestroy(state *terraform.State) error {
	client := getConfigurationClient(testAccProvider.Meta().(map[string]interface{}))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "limelight_delivery" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, resp, err := client.GetDeliveryServiceInstance(resourceID)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("error retrieving delivery configuration with ID %s. Error: %v", resourceID, err)
		}

	}
	return fmt.Errorf("delivery configuration still exists")
}

func testAccLimelightDeliveryExists(testResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		client := getConfigurationClient(testAccProvider.Meta().(map[string]interface{}))

		rs, ok := state.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("delivery configuration %s not found in resources", testResourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("delivery configuration ID not set in resources ")
		}

		serviceInst, resp, err := client.GetDeliveryServiceInstance(resourceID)
		if err != nil {
			return fmt.Errorf("error retrieving delivery configuration for %s: %s", resourceID, err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error checking if delivery configuration %s exists. HTTP return code was %d", resourceID, resp.StatusCode)
		}

		if serviceInst.UUID == resourceID {
			return nil
		}

		return fmt.Errorf("delivery configuration instance with ID %s wasn't found", resourceID)
	}
}

func testAccLimelightDeliveryMinimalTemplate() string {
	return fmt.Sprintf(`
resource "limelight_delivery" "test_delivery" {
	shortname          = "%s"
	published_hostname = "terraform-test.%s.s.llnwi.net"
	published_path     = "/"
	source_hostname    = "dummy-origin.llnw.net"
	source_path        = "/"
	
	protocol_set {
		published_protocol = "https"
		source_protocol    = "https"
	}
}`, getShortname(), getShortname())
}

func testAccLimelightDeliveryBasicTemplate() string {
	return fmt.Sprintf(`
resource "limelight_delivery" "test_delivery" {
	shortname          = "%s"
	published_hostname = "terraform-test-basic.%s.s.llnwi.net"
	published_path     = "/"
	source_hostname    = "dummy-origin-basic.llnw.net"
	source_path        = "/"
	
	protocol_set {
		published_protocol = "https"
		source_protocol    = "https"
		option {
			name       = "reply_send_header"
			parameters = ["X-LLNW-Test", "123"]
		}
	}
}`, getShortname(), getShortname())
}

func testAccLimelightDeliveryBasicTemplate_update() string {
	return fmt.Sprintf(`
resource "limelight_delivery" "test_delivery" {
	shortname          = "%s"
	published_hostname = "terraform-test-basic.%s.s.llnwi.net"
	published_path     = "/publish"
	source_hostname    = "dummy-origin-basic.llnw.net"
	source_path        = "/source"
	
	protocol_set {
		published_protocol = "http"
		source_protocol    = "http"
		option {
			name       = "reply_send_header"
			parameters = ["X-LLNW-Test", "123"]
		}
		option {
			name       = "genreply"
			parameters = ["200"]
		}
	}
}`, getShortname(), getShortname())
}

func testAccLimelightDeliveryImportTemplate() string {
	return fmt.Sprintf(`
resource "limelight_delivery" "test_delivery" {
	shortname          = "%s"
	published_hostname = "terraform-test-import.%s.s.llnwi.net"
	published_path     = "/"
	source_hostname    = "dummy-origin-import.llnw.net"
	source_path        = "/"
	
	protocol_set {
		published_protocol = "https"
		source_protocol    = "https"
	}
}`, getShortname(), getShortname())
}
