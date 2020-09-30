---
layout: "limelight"
page_title: "Limelight: limelight_edgefunction"
sidebar_current: "docs-limelight-resource-edgefunction"
description: A resource that can be used to manage EdgeFunctions.
---

# limelight_edgefunction

This resource provides a way to manage EdgeFunctions in Limelight Networks.
For more details see the [API docs](https://support.limelight.com/public/openapi/edgefunctions/index.html#tag/Function-Management)

## Example Usage

```hcl
provider "archive" {}

data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = "function"
  output_path = "function.zip"
}

resource "limelight_edgefunction" "hello_world" {
  shortname        = var.shortname
  name             = "hello_world_terraform"
  description      = "A simple hello world function, provisioned with Terraform"
  function_archive = data.archive_file.function_archive.output_path
  function_sha256  = filesha256(data.archive_file.function_archive.output_path)
  handler          = "hello_world.handler"
  runtime          = "python3"
  memory           = 256
  timeout          = 2000
  can_debug        = false
  environment_variable {
    name  = "NAME"
    value = "World"
  }
}
```

## Argument Reference

The following arguments are supported:

* `shortname` - (Required) The account name (shortname).
* `name` - (Required) A unique name for the EdgeFunction.
* `description` - (Optional) A description for the EdgeFunction.
* `function_archive` - (Required) Path to the function archive (zip file).
* `handler` - (Required) Handler that's run when the EdgeFunction is invoked.
* `runtime` - (Required) The runtime for the EdgeFunction.
* `memory` - (Optional) The memory allocated to the EdgeFunction. Defaults to `256`. CPU is allocated
  proportional to memory.
* `timeout` - (Optional) Timeout for the EdgeFunction execution in milliseconds. Defaults to `5000`.
* `can_debug` - (Optional) Boolean flag to enable debug IO. Defaults to `false`.
* `function_sha256` - (Required) The SHA256 value of the `function_archive`.
* `reserved_concurrency` - (Optional) Sets the reserved concurrency for the EdgeFunction. Defaults to `0`.
* `environment_variable` - (Optional) Zero or more environment variables for the EdgeFunction as child blocks:
  * `name` - (Required) The environment variable name.
  * `value` - (Required) The environment variable value.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision_id` - Revision number of the EdgeFunction.

## Importing

An existing EdgeFunction can be [imported](https://www.terraform.io/docs/import/index.html) into this resource, via the
following command:

```
terraform import limelight_edgefunction.my_func FUNCTION_ID
```

The above command imports the EdgeFunction named `my_func` with the ID `FUNCTION_ID` where `FUNCTION_ID`
is of the form `<shortname>:<function_name>`.
