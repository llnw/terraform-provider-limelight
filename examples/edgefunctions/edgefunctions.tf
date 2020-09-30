# Configure the Limelight Provider
terraform {
  required_providers {
    limelight = {
      source  = "llnw/limelight"
      version = "~> 1.0.0"
    }
  }
}

provider "limelight" {
  username = var.llnw_username
  api_key  = var.llnw_api_key
}

provider "archive" {}

variable "llnw_username" {
  type = string
}

variable "llnw_api_key" {
  type = string
}

variable "shortname" {
  type = string
}

variable "published_hostname" {
  type = string
}

# The archive file created from the directory containing your EdgeFunction code
data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = "function"
  output_path = "function.zip"
}

# An EdgeFunction created from the zip archive above
resource "limelight_edgefunction" "hello_world" {
  shortname        = var.shortname
  name             = "hello_world_terraform"
  description      = "A simple hello world function, provisioned with Terraform"
  function_archive = data.archive_file.function_archive.output_path
  function_sha256  = filesha256(data.archive_file.function_archive.output_path)
  handler          = "hello_world.handler"
  runtime          = "python3"
  memory           = 256
  timeout          = 4000
  can_debug        = true
  environment_variable {
    name  = "NAME"
    value = "World"
  }
}

# A CDN configuration to allow the EdgeFunction to be invoked over HTTPS
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
