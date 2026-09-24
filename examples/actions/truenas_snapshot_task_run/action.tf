# Run a periodic snapshot task now (Terraform 1.14+).
# id is the numeric task id — reference the managing resource rather than
# hardcoding it (or discover it with `midclt call pool.snapshottask.query`).
action "truenas_snapshot_task_run" "example" {
  config {
    id = truenas_periodic_snapshot_task.hourly.id
  }
}
