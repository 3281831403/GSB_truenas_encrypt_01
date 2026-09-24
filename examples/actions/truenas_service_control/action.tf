# Restart a service (Terraform 1.14+).
action "truenas_service_control" "example" {
  config {
    service = "cifs"
    verb    = "RESTART"
  }
}
