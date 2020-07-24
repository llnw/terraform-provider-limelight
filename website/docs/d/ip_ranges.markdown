---
layout: "limelight"
page_title: "Limelight: limelight_ip_ranges"
sidebar_current: "docs-limelight-datasource-ip-ranges"
description: An IP Ranges data source.
---

# limelight_ip_ranges

This data source provides information about the default IP Ranges for the Edge Cache server configured
in Limelight Networks.

## Example Usage

```hcl
data "limelight_ip_ranges" "ip_ranges" {}
```

## Argument Reference

None

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `ip_ranges` - A `string` list where each element is an IP address.
* `version` - The version for the current list of `ip_ranges`.
