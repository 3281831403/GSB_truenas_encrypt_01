# Query files (*.tfquery.hcl) are run with `terraform query` (Terraform
# v1.14+) and use a resource's list resource to enumerate objects that
# already exist on the TrueNAS system, independent of any Terraform state.
# Each result includes the resource's identity, which can be pasted into
# an `import { identity = { ... } }` block (see examples/import-identity).
#
# Run with:
#   terraform query

list "truenas_user" "all" {}

list "truenas_pool" "all" {}
