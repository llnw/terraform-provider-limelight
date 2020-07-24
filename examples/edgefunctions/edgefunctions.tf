provider "limelight" {}

provider "archive" {}

data "archive_file" "function_archive" {
  type        = "zip"
  source_dir  = "function"
  output_path = "function.zip"
}

resource "limelight_edgefunction" "hello_world" {
  shortname        = "llnwfaas"
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
