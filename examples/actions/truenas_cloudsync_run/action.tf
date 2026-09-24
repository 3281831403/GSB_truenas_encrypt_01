# Run a cloud sync task now (Terraform 1.14+).
# id is the numeric task id — reference the managing resource rather than
# hardcoding it (or discover it with `midclt call cloudsync.query`).
action "truenas_cloudsync_run" "example" {
  config {
    id = truenas_cloudsync_task.offsite.id
  }
}
