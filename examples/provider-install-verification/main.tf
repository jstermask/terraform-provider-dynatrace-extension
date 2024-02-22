terraform {
  required_providers {
    dynatrace-extension = {
      source = "registry.terraform.io/jstermask/dynatrace-extension"
    }
  }
}


resource "dynatrace-extension_extension" "custom_jmx_testext_extension" {
  payload = file("./custom.jmx.testext.json")
}
