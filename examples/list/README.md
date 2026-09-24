# Listing / discovery with `terraform query`

List resources let you enumerate objects that already exist on a TrueNAS
system, independent of Terraform state. Each result carries the object's
**identity**, which pastes straight into an
`import { identity = { ... } }` block (see `examples/import-identity`).

Requires **Terraform v1.14+** and a configured `truenas` provider in the same
directory (endpoint + credentials — see `examples/provider`).

## Files

- `main.tfquery.hcl` — a minimal example: list users and pools.
- `discover-all.tfquery.hcl` — **full-system discovery**: one `list` block for
  every list-capable resource type (83 of them). Pointed at an existing box it
  prints a complete inventory of everything the provider can manage.

## Run

```sh
terraform init      # only needed once, for the provider
terraform query     # human-readable table
terraform query -json > inventory.jsonl   # machine-readable
```

Each row is `list.<type>.<name>  id=<identity>  <display name>`, e.g.:

```
list.truenas_pool.all      id=1            space
list.truenas_dataset.all   id=tank/apps    tank/apps
list.truenas_ssh_config.all id=ssh_config  ssh_config
```

## Notes

- Singleton config resources (e.g. `truenas_ssh_config`) return exactly one
  row with a fixed string id.
- Resources that share a backend query are type-filtered, so
  `truenas_dataset` lists only filesystems and `truenas_zvol` only volumes.
- `filesystem_acl` and `filesystem_permissions` have identity (for imports) but
  no list resource — they are keyed by an arbitrary path and cannot be
  enumerated.
