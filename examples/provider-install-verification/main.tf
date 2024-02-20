terraform {
  required_providers {
    dynaext = {
      source = "registry.terraform.io/jstermask/dynatrace-extension"
    }
  }
}

provider "dynaext" {}

data "dynaext_list_custom" "example" {}
