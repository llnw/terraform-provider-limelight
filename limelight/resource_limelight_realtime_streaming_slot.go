package limelight

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/llnw/llnw-sdk-go/configuration"
)

func resourceLimelightRealtimeStreamingSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceLimelightRealtimeStreamingSlotCreate,
		Read:   resourceLimelightRealtimeStreamingSlotRead,
		Delete: resourceLimelightRealtimeStreamingSlotDelete,
		Schema: map[string]*schema.Schema{
			"shortname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"north-america",
					"europe",
					"asia-pacific",
				}, false),
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     profilesElem,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"ip_geo_match": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mediavault_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"mediavault_secret_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"wait_for_provisioning": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceLimelightRealtimeStreamingSlotCreate(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	shortname := d.Get("shortname").(string)
	name := d.Get("name").(string)
	region := d.Get("region").(string)
	profiles := expandProfiles(d.Get("profile").(*schema.Set))
	password := d.Get("password").(string)
	ipGeoMatch := d.Get("ip_geo_match").(string)
	mediaVaultSecretKey, mediaVaultEnabled := d.GetOk("mediavault_secret_key")

	realtimeStreamingSlot := &configuration.RealtimeStreamingSlot{
		Name:                name,
		Region:              region,
		Profiles:            profiles,
		Password:            password,
		IPGeoMatch:          ipGeoMatch,
		MediaVaultEnabled:   mediaVaultEnabled,
		MediaVaultSecretKey: mediaVaultSecretKey.(string),
	}

	log.Printf("[INFO] Creating Realtime Streaming Slot: %s", name)
	slotResponse, _, err := c.CreateRealtimeStreamingSlot(shortname, realtimeStreamingSlot)

	if err != nil {
		return fmt.Errorf("error creating Realtime Streaming Slot: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", shortname, slotResponse.Id))

	waitForProvisioning := d.Get("wait_for_provisioning").(bool)
	if waitForProvisioning {
		log.Printf("[INFO] Waiting for provisioning of Realtime Stream Slot %s", name)
		err = waitForLimelightRealtimeStreamingSlotProvision(d, m)
		if err != nil {
			log.Printf("[WARN] Failed waiting for Realtime Streaming Slot provisioning of %s due to: %v", name, err)
		}
	}

	return resourceLimelightRealtimeStreamingSlotRead(d, m)
}

func resourceLimelightRealtimeStreamingSlotRead(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	shortname, slotID, err := resourceLimelightRealtimeStreamingSlotSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Fetching Realtime Streaming Slot %s", slotID)
	realtimeStreamingSlot, resp, err := c.GetRealtimeStreamingSlot(slotID, shortname)

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Realtime Streaming Slot %s not found", slotID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading Realtime Streaming Slot: %s", err)
	}

	// NOTE: password does not come back on GET requests
	d.Set("name", realtimeStreamingSlot.Name)
	d.Set("shortname", shortname)
	d.Set("region", realtimeStreamingSlot.Region)
	d.Set("profile", flattenProfiles(realtimeStreamingSlot.Profiles))
	d.Set("ip_geo_match", realtimeStreamingSlot.IPGeoMatch)
	d.Set("mediavault_enabled", realtimeStreamingSlot.MediaVaultEnabled)
	d.Set("mediavault_secret_key", realtimeStreamingSlot.MediaVaultSecretKey)
	d.Set("state", realtimeStreamingSlot.State)

	return nil
}

func resourceLimelightRealtimeStreamingSlotDelete(d *schema.ResourceData, m interface{}) error {
	c := getConfigurationClient(m)

	shortname, slotID, err := resourceLimelightRealtimeStreamingSlotSplitID(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting Realtime Streaming Slot %s", slotID)
	_, err = c.DeleteRealtimeStreamingSlot(slotID, shortname)

	if err != nil {
		return fmt.Errorf("error deleting Realtime Streaming slot: %s", err)
	}

	return nil
}

func waitForLimelightRealtimeStreamingSlotProvision(d *schema.ResourceData, m interface{}) error {

	shortname, slotID, err := resourceLimelightRealtimeStreamingSlotSplitID(d.Id())

	if err != nil {
		return err
	}

	client := getConfigurationClient(m)

	pendingStates := []string{configuration.SlotStatePending}
	targetStates := []string{configuration.SlotStateReady, configuration.SlotStateFailed}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {

			realtimeStreamingSlot, resp, err := client.GetRealtimeStreamingSlot(slotID, shortname)
			if err != nil && resp.StatusCode == http.StatusNotFound {
				d.Set("state", "NOT_FOUND")
				return nil, "NOT_FOUND", err
			}
			d.Set("state", realtimeStreamingSlot.State)
			return realtimeStreamingSlot, realtimeStreamingSlot.State, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		MinTimeout: 1 * time.Second,
		Delay:      1 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("failed waiting for Realtime Streaming Slot %s: %v", slotID, err)
	}
	return nil
}

func resourceLimelightRealtimeStreamingSlotSplitID(id string) (string, string, error) {
	shortName, slotID, err := splitSeparatedPair(id, ":")

	if err != nil {
		return "", "", fmt.Errorf("Realtime Streaming Slot ID in unexpected format (expected '<shortname>:<slot ID>'): %s", id)
	}

	return shortName, slotID, nil
}

var profilesElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"video_bitrate": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"audio_bitrate": &schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
	},
}

func flattenProfiles(expandedProfiles []configuration.RealtimeStreamingProfile) *schema.Set {
	flattenedProfiles := make([]interface{}, len(expandedProfiles), len(expandedProfiles))

	for i, v := range expandedProfiles {
		m := make(map[string]interface{})
		m["video_bitrate"] = v.VideoBitrate
		m["audio_bitrate"] = v.AudioBitrate
		flattenedProfiles[i] = m
	}

	return schema.NewSet(schema.HashResource(profilesElem), flattenedProfiles)
}

func expandProfiles(flattenedProfiles *schema.Set) []configuration.RealtimeStreamingProfile {
	expandedProfiles := make([]configuration.RealtimeStreamingProfile, flattenedProfiles.Len(), flattenedProfiles.Len())

	for i, v := range flattenedProfiles.List() {
		rawProfile := v.(map[string]interface{})
		opt := configuration.RealtimeStreamingProfile{
			VideoBitrate: rawProfile["video_bitrate"].(int),
			AudioBitrate: rawProfile["audio_bitrate"].(int),
		}
		expandedProfiles[i] = opt
	}

	return expandedProfiles
}
