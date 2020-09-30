---
layout: "limelight"
page_title: "Limelight: limelight_delivery"
sidebar_current: "docs-limelight-resource-delivery"
description: A resource that can be used to configure Content Delivery.
---

# limelight_delivery

This resource provides a way to configure Content Delivery in Limelight Networks.

## Example Usage

```hcl
resource "limelight_delivery" "example_website" {
  shortname          = var.shortname
  published_hostname = "www.example.com"
  published_path     = "/"
  source_hostname    = "origin.example.com"
  source_path        = "/"

  protocol_set {
    published_protocol = "https"
    source_protocol    = "https"
    option {
      name       = "refresh_absmin"
      parameters = ["3600"]
    }
    option {
      name       = "refresh_absmax"
      parameters = ["3600"]
    }
    option {
      name       = "reply_send_header"
      parameters = ["X-Delivered-By", "LLNW"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `shortname` - (Required) The account name (shortname).
* `service_profile` - (Optional) Service profile to use. Defaults to `LLNW-Generic`.
* `published_hostname` - (Required) Published hostname for the content.
* `published_path` - (Required) Published path for the content.
* `source_hostname` - (Required) Source (origin) hostname for the content.
* `source_path` - (Required) Source path on the origin for the content.
* `protocol_set` - (Required) Protocol configuration for the delivery as child blocks (max of 2):
  * `published_protocol` - (Required) Published protocol to use (e.g. `http`, `https`).
  * `source_protocol` - (Required) Source protocol to use (e.g. `http`, `https`).
  * `source_port` - (Optional) Source port to use. Defaults to `80` for http and `443` for https.
  * `option` - (Optional) Protocol options to use specified as child blocks:
      * `name` - (Required) Option name.
      * `parameters` - (Required) List of string parameters for the option.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The delivery ID.
* `version_number` - The delivery version.

## Importing

An existing Delivery configuration can be [imported](https://www.terraform.io/docs/import/index.html) into this
resource, via the following command:

```
terraform import limelight_delivery.example_website UUID
```

The above command imports the Delivery named `example_website` with the ID `UUID`.