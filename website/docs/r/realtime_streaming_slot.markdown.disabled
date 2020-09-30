---
layout: "limelight"
page_title: "Limelight: limelight_realtime_streaming_slot"
sidebar_current: "docs-limelight-resource-realtime-streaming-slot"
description: A resource that can be used to configure Realtime Streaming Slots.
---

# limelight_realtime_streaming_slot

This resource provides a way to configure Realtime Streaming Slots in Limelight Networks.

## Example Usage

```hcl
resource "limelight_realtime_streaming_slot" "test_streaming" {
  shortname             = var.shortname
  name                  = "terraform-rts"
  region                = "north-america"
  password              = "secretpassw0rd"
  wait_for_provisioning = true
  profile {
	  video_bitrate = 1800000
	  audio_bitrate = 192000
  }
  profile {
	  video_bitrate = 2400000
	  audio_bitrate = 192000
  }
}
```

## Argument Reference

The following arguments are supported:

* `shortname` - (Required) The account name (shortname).
* `name` - (Required) Name of the Realtime Streaming Slot.
* `region` - (Required) Region for the Realtime Streaming Slot. Must be one of `north-america`, `europe` or `asia-pacific`.
* `password` - (Optional) Password to use for the Realtime Streaming Slot.
* `ip_geo_match` - (Optional) IP/Geo matching for the Realtime Streaming Slot. Note this can cause the Slot
  provisioning to take in excess of 20 minutes.
* `mediavault_secret_key` - (Optional) Enables and sets the secret key for Media Vault. Note this can cause
  the Slot to take in excess of 20 minutes to complete provisioning.
* `wait_for_provisioning` - (Optional) Boolean flag to enable waiting for provisioning of the Realtime Streaming
  Slot. Default is `true`.
* `profile` - (Required) One or more profiles for the Realtime Streaming Slot as child blocks:
  * `video_bitrate` - (Required) Video bitrate for the Realtime Streaming Slot.
  * `audio_bitrate` - (Required) Audio bitrate for the Realtime Streaming Slot.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - Current provisioning state of the Realtime Streaming Slot.
  
## Importing

Not supported.
