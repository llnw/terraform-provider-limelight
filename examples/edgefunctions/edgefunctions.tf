provider "limelight" {}

provider "archive" {}

variable "shortname" {
  type = string
}

variable "published_hostname" {
  type = string
}

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

resource "limelight_delivery" "edgefunction_cdn_config" {
  shortname          = var.shortname
  published_hostname = var.published_hostname
  published_path     = "/my-function/"
  source_hostname    = "apis.llnw.com"
  source_path        = "/ef-api/v1/${var.shortname}/functions/${limelight_edgefunction.hello_world.name}/epInvoke/"
  service_profile    = "${var.shortname}-EdgeFunctions"
  protocol_set {
    published_protocol = "https"
    source_protocol    = "https"
  }
}
