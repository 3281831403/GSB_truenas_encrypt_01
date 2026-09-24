# Restart the TrueNAS web UI after a config change (Terraform 1.14+).
action "truenas_ui_restart" "example" {
  config {
    delay = 0
  }
}
