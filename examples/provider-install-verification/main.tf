terraform {
  required_providers {
    dynaext = {
      source = "registry.terraform.io/jstermask/dynatrace-extension"
    }
  }
}

provider "dynaext" {}

resource "dynaext_extension" "example" {
  payload = file("./custom.jmx.testext.json")
}
