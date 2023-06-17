resource "gtm_trigger" "test_trigger_1" {
  name = "test trigger 1"
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
          value = "event-name"
        }
      ]
    }
  ]
}
