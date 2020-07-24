package limelight

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccResourceLimelightRealtimeStreamingSlot_minimal(t *testing.T) {
	testResourceName := "limelight_realtime_streaming_slot.test_streaming"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightRealtimeStreamingSlotCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightRealtimeStreamingSlotMinimalTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightRealtimeStreamingSlotExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", "terraform-min"),
					resource.TestCheckResourceAttr(testResourceName, "region", "europe"),
					resource.TestCheckResourceAttr(testResourceName, "mediavault_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "profile.#", "1"),
				),
			},
		},
	})
}

func testAccResourceLimelightRealtimeStreamingSlot_update(t *testing.T) {
	testResourceName := "limelight_realtime_streaming_slot.test_streaming"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccLimelightRealtimeStreamingSlotCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLimelightRealtimeStreamingSlotBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightRealtimeStreamingSlotExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", "terraform-update"),
					resource.TestCheckResourceAttr(testResourceName, "region", "europe"),
					// TODO: enable media vault when it doesn't take 20+ minutes to provision via API
					//resource.TestCheckResourceAttr(testResourceName, "mediavault_secret_key", "mpassw0rd"),
					resource.TestCheckResourceAttr(testResourceName, "profile.#", "1"),
				),
			},
			{
				Config: testAccLimelightRealtimeStreamingSlotBasicTemplate_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccLimelightRealtimeStreamingSlotExists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttr(testResourceName, "shortname", getShortname()),
					resource.TestCheckResourceAttr(testResourceName, "name", "terraform-update"),
					resource.TestCheckResourceAttr(testResourceName, "region", "north-america"),
					//resource.TestCheckResourceAttr(testResourceName, "mediavault_secret_key", "mpassw0rd2"),
					resource.TestCheckResourceAttr(testResourceName, "profile.#", "2"),
					// TODO: enable ip_geo_match when it doesn't take 20+ minutes to provision via API
					//resource.TestCheckResourceAttr(testResourceName, "ip_geo_match", "true"),
				),
			},
		},
	})
}

func testAccLimelightRealtimeStreamingSlotCheckDestroy(state *terraform.State) error {
	client := getConfigurationClient(testAccProvider.Meta().(map[string]interface{}))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "limelight_realtime_streaming_slot" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		shortname, slotID, err := resourceLimelightRealtimeStreamingSlotSplitID(resourceID)

		if err != nil {
			return err
		}

		_, resp, err := client.GetRealtimeStreamingSlot(slotID, shortname)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("error retrieving Realtime Streaming Slot with ID %s. Error: %v", resourceID, err)
		}

	}
	return fmt.Errorf("Realtime Streaming Slot still exists")
}

func testAccLimelightRealtimeStreamingSlotExists(testResourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		client := getConfigurationClient(testAccProvider.Meta().(map[string]interface{}))

		rs, ok := state.RootModule().Resources[testResourceName]
		if !ok {
			return fmt.Errorf("Realtime Streaming Slot %s not found in resources", testResourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Realtime Streaming Slot ID not set in resources ")
		}

		shortname, slotID, err := resourceLimelightRealtimeStreamingSlotSplitID(resourceID)

		if err != nil {
			return err
		}

		slot, resp, err := client.GetRealtimeStreamingSlot(slotID, shortname)
		if err != nil {
			return fmt.Errorf("error retrieving Realtime Streaming Slot: %s", resourceID)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("rrror while checking if Realtime Streaming Slot %s exists. HTTP return code was %d", resourceID, resp.StatusCode)
		}

		if slot.Id == slotID {
			return nil
		}

		return fmt.Errorf("Realtime Streaming Slot instance with ID %s wasn't found", resourceID)
	}
}

func testAccLimelightRealtimeStreamingSlotMinimalTemplate() string {
	return fmt.Sprintf(`
resource "limelight_realtime_streaming_slot" "test_streaming" {
	shortname             = "%s"
	name                  = "terraform-min"
	region                = "europe"
	wait_for_provisioning = true
	profile {
		video_bitrate = 1800000
		audio_bitrate = 192000
	}
}`, getShortname())
}

func testAccLimelightRealtimeStreamingSlotBasicTemplate() string {
	return fmt.Sprintf(`
resource "limelight_realtime_streaming_slot" "test_streaming" {
	shortname             = "%s"
	name                  = "terraform-update"
	region                = "europe"
	password              = "passw0rd"
	wait_for_provisioning = true
	profile {
		video_bitrate = 1800000
		audio_bitrate = 192000
	}
}`, getShortname())
}

func testAccLimelightRealtimeStreamingSlotBasicTemplate_update() string {
	return fmt.Sprintf(`
resource "limelight_realtime_streaming_slot" "test_streaming" {
	shortname             = "%s"
	name                  = "terraform-update"
	region                = "north-america"
	password              = "passw0rd2"
	wait_for_provisioning = true
	profile {
		video_bitrate = 1800000
		audio_bitrate = 192000
	}
	profile {
		video_bitrate = 2400000
		audio_bitrate = 192000
	  }
}`, getShortname())
}
