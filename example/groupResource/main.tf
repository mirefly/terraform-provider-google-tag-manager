terraform {
  required_providers {
    gtm = {
      source  = "mirefly/google-tag-manager"
      version = "0.0.1"
    }
  }
}

provider "gtm" {
  credential_file            = "../credentials-77b14e38b4dd.json"
  account_id                 = "6105084028"
  container_id               = "119458552"
  workspace_name             = "my-workspace"
  max_api_queries_per_minute = 15
}

resource "gtm_variable_group" "test_variable_group" {
  elements = {
    "variable 1 in group" : {
      name  = "variable 1 in group"
      type  = "v"
      notes = "Generated by terraform. Do not edit it."
      parameter = [
        {
          key   = "name"
          type  = "template"
          value = "parameters.alph1a"
        }
      ]
    },
    "variable 2 in group" : {
      name  = "variable 2 in group"
      type  = "v"
      notes = "Generated by terraform. Do not edit it."
      parameter = [
        {
          key   = "name"
          type  = "template"
          value = "parameters.beta"
        }
      ]
    }
  }
}