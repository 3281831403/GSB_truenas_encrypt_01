# Run a replication task now (Terraform 1.14+).
# id is the numeric task id — reference the managing resource rather than
# hardcoding it (or discover it with `midclt call replication.query`).
action "truenas_replication_run" "example" {
  config {
    id = truenas_replication_task.nightly.id
  }
}
