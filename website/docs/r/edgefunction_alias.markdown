---
layout: "limelight"
page_title: "Limelight: limelight_edgefunction_alias"
sidebar_current: "docs-limelight-resource-edgefunction-alias"
description: A resource that can be used to manage EdgeFunction Aliases.
---

# limelight_edgefunction_alias

This resource provides a way to manage EdgeFunction Aliases in Limelight Networks.
In order to create a new EdgeFunction Alias for a version other than `$LATEST`, the version
must be created out-of-band using the EdgeFunction CLI or REST API.
For more details see the [API docs](https://support.limelight.com/public/openapi/edgefunctions/index.html#tag/Aliases)

## Example Usage

```hcl
resource "limelight_edgefunction_alias" "my_alias" {
  shortname        = var.shortname
  name             = "PROD"
  function_name    = limelight_edgefunction.my_edgefunc.name
  function_version = "$LATEST"
  description      = "Alias for PROD"
}
```

## Argument Reference

The following arguments are supported:

* `shortname` - (Required) The account name (shortname).
* `name` - (Required) A unique name for the EdgeFunction alias.
* `description` - (Optional) A description for the EdgeFunction alias.
* `function_name` - (Required) The EdgeFunction's name to create the alias for.
* `function_version` - (Required) The EdgeFunction's version to create the alias for.
  If a version other than `$LATEST` is used, this version must already exist.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision_id` - Revision number of the EdgeFunction alias.

## Importing

An existing EdgeFunction Alias can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import limelight_edgefunction_alias.my_alias ALIAS_UUID
```

The above command imports the EdgeFunction Alias named `my_alias` with the ID `ALIAS_UUID` where
`ALIAS_UUID` is of the form `<shortname>:<function_name>:<alias_name>`.
