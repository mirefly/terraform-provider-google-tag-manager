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

resource "gtm_tag_group" "test_tag_group" {
  elements = {
    "tag 1 in group" : {
      name  = "tag 1 in group"
      type  = "gaawe"
      notes = "Generated by terraform. Do not edit it."
      parameter = [
        {
          key   = "eventName"
          type  = "template"
          value = "event1"
        },
        {
          key  = "eventParameters"
          type = "list"
          list = [{
            type = "map"
            map = [
              {
                type  = "template"
                key   = "name"
                value = "eventName"
              },
              {
                type  = "template"
                key   = "value"
                value = "eventValue"
            }]
          }]
        },
        {
          key   = "measurementId",
          type  = "template"
          value = "G-A2ABC2ABCD"
        }
      ]
    },
    "tag 2 in group" : {
      name = "tag 2 in group"
      type = "html",
      parameter = [
        {
          type  = "template",
          key   = "html",
          value = "\u003cp\u003e this is custom html tag \u003c/p"
        },
        {
          type  = "boolean",
          key   = "supportDocumentWrite",
          value = "false"
        }
      ]
    }
  }
}

resource "gtm_trigger_group" "test_trigger_group" {
  elements = {
    "trigger 1 in group" : {
      name = "trigger 1 in group"
      type = "customEvent"
      custom_event_filter = [
        {
          type = "equals",
          parameter = [
            {
              type  = "template",
              key   = "arg0",
              value = "{{_event}}"
            },
            {
              type  = "template",
              key   = "arg1",
              value = "event-name-1"
            }
          ]
        }
      ]

    },
    "trigger 2 in group" : {
      name = "trigger 2 in group"
      type = "customEvent"
      custom_event_filter = [
        {
          type = "equals",
          parameter = [
            {
              type  = "template",
              key   = "arg0",
              value = "{{_event}}"
            },
            {
              type  = "template",
              key   = "arg1",
              value = "event-name-2"
            }
          ]
        }
      ]

    }
  }
}
