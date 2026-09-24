# Stop an app (Terraform 1.14+).
# app_name is the app's INSTANCE name, unique per system — installing the same
# catalog app twice yields two distinct names. Reference the managing resource
# rather than hardcoding it (or discover it with `midclt call app.query`).
action "truenas_app_stop" "example" {
  config {
    app_name = truenas_app.nextcloud.name
  }
}
